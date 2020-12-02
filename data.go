// Copyright 2015 Monmohan Singh. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xferspdy provides the basic interfaces around binary diff and patching process
package xferspdy

import (
	"crypto/sha256"
	"fmt"
	"hash/adler32"
	"io"
	"os"
	"sync"

	"github.com/golang/glog"
)

var (
	DEFAULT_GENERATOR = &FingerprintGenerator{ConcurrentMode: true, NumWorkers: 8}
)

type FingerprintGenerator struct {
	Source         io.Reader
	BlockSize      uint32
	ConcurrentMode bool
	NumWorkers     int
}

// Block represent a byte slice from the file. For each block, following are computed.
//
// * Adler-32 and SHA256 checksum,
//
// * Start and End byte pos of the block,
//
// * Whether or not its a data block -If this is a data block, RawBytes will capture the byte data represented by this block
type Block struct {
	Start, End int64
	Checksum32 uint32
	Sha256hash [sha256.Size]byte
	HasData    bool
	RawBytes   []byte
}

// Fingerprint of a given File, encapsulates the following mapping -
//   Adler-32 hash of Block --> SHA256 hash of Block -->Block
// Also stores the block size and the source
type Fingerprint struct {
	Blocksz  uint32
	BlockMap map[uint32]map[[sha256.Size]byte]Block
	Source   string
}

func (b Block) String() string {
	return fmt.Sprintf("Start %d End %d adler %d HasData %v \n", b.Start, b.End, b.Checksum32, b.HasData)
}

func (f Fingerprint) String() string {
	buf := fmt.Sprintf("Block size=%d, Source=%s\n", f.Blocksz, f.Source)
	for k, v := range f.BlockMap {
		buf += fmt.Sprintf("Checksum32=%d\n", k)
		for sha, blk := range v {
			buf += fmt.Sprintf("\tSHA Hash=%d,Block=%v\n", sha, blk)
		}
	}
	return buf
}

// Generate creates a finger print using the FingerprintGenerator.
// Processing i.e. concurrent or sequential depends on the generator field ConcurrentMode
func (g *FingerprintGenerator) Generate() *Fingerprint {
	if g.ConcurrentMode {
		return g.genConcurrent()
	} else {
		return g.genSequential()
	}
}

// NewFingerprintFromReader creates a Fingerprint for a given reader and blocksize.
// By default it does concurrent processing of blocks to generate fingerprint.
// However if the number of blocks is small <50 , then caller should use sequential generation,
// since the concurrent processing would not add much value.
// Or use the function NewFingerrprint(file, blocksize) when dealing with files, which switches
// mode based on the number of blocks.
// Number of blocks can be calculated as file size/block size
func NewFingerprintFromReader(r io.Reader, blocksz uint32) *Fingerprint {
	DEFAULT_GENERATOR.Source = r
	DEFAULT_GENERATOR.BlockSize = blocksz
	return DEFAULT_GENERATOR.Generate()

}
func (g *FingerprintGenerator) genSequential() *Fingerprint {
	bufz := make([]byte, g.BlockSize)

	n, start := 0, int64(0)

	var (
		err   error
		block Block
	)

	fngprt := Fingerprint{
		Blocksz: g.BlockSize, BlockMap: make(map[uint32]map[[sha256.Size]byte]Block)}

	for err == nil {
		n, err = g.Source.Read(bufz)
		if err == nil {
			block = Block{Start: start, End: start + int64(n),
				Checksum32: adler32.Checksum(bufz[0:n]),
				Sha256hash: sha256.Sum256(bufz[0:n])}
			addBlock(&fngprt, &block)
			start = block.End
		} else {
			if err == io.EOF {
				glog.V(2).Infoln("Fingerprint generation: Reader read complete")
			} else {
				glog.Fatal(err)
			}

		}

	}

	return &fngprt

}

// NewFingerprint creates a Fingerprint for a given reader and blocksize.
func (g *FingerprintGenerator) genConcurrent() *Fingerprint {
	fngprt := Fingerprint{
		Blocksz: g.BlockSize, BlockMap: make(map[uint32]map[[sha256.Size]byte]Block)}

	blkin := readBlocks(g.Source, g.BlockSize, g.NumWorkers)
	blkout := fillBlocks(blkin, g.NumWorkers)
	for b := range blkout {
		addBlock(&fngprt, b)
	}

	return &fngprt

}

// NewFingerprint creates a Fingerprint for a given file and blocksize.
// By default it does concurrent processing of blocks to generate fingerprint.
// The generation is switched to sequential mode if the number of blocks is less than 50.
func NewFingerprint(filename string, blocksize uint32) *Fingerprint {
	file, e := os.Open(filename)
	defer file.Close()
	if e != nil {
		glog.Fatalf("Unable to open file %s %s", filename, e)
	}
	fileInfo, _ := file.Stat()
	numblocks := (fileInfo.Size() / int64(blocksize))
	var f *Fingerprint
	if numblocks < 50 {
		//switch to sequential mode
		g := &FingerprintGenerator{Source: file, ConcurrentMode: false, BlockSize: blocksize}
		f = g.Generate()

	} else {
		//use default generator
		f = NewFingerprintFromReader(file, blocksize)
	}

	f.Source = filename
	return f

}

// addBlock adds the hashed block to the Fingerprint struct
func addBlock(f *Fingerprint, b *Block) {

	glog.V(3).Infof("Adding Block %v ", *b)
	if sha2blk := f.BlockMap[b.Checksum32]; sha2blk == nil {
		f.BlockMap[b.Checksum32] = make(map[[sha256.Size]byte]Block)
	}
	f.BlockMap[b.Checksum32][b.Sha256hash] = *b

}

// readBlocks reads blocksize bytes from the reader into memory
// numhashers determines the buffer size of the output channel where the method places the blocks which
// been read in
func readBlocks(r io.Reader, blocksize uint32, numhashers int) chan *Block {

	blkin := make(chan *Block, numhashers)
	n, start := 0, int64(0)
	go func() {
		var err error
		defer close(blkin)
		for err == nil {
			bufz := make([]byte, blocksize)

			n, err = r.Read(bufz)

			if err == nil {
				block := Block{Start: start, End: start + int64(n), RawBytes: bufz, HasData: true}
				blkin <- &block
				start += int64(n)
			} else {
				if err == io.EOF {
					glog.V(2).Infoln("Fingerprint generation: Reader read complete")

				} else {
					glog.Fatal(err)
				}

			}
		}

	}()
	return blkin

}

// fillBlocks takes an input channel with the bytes read from disk and creates the Checksum and SHAHashes
// numworkers is the number of go routines used for processing
func fillBlocks(in chan *Block, numhashers int) chan *Block {
	out := make(chan *Block)
	var wg sync.WaitGroup

	wg.Add(numhashers)
	for i := 0; i < numhashers; i++ {
		go func() {
			for blkptr := range in {
				buf := blkptr.RawBytes[0:(blkptr.End - blkptr.Start)]
				blkptr.Checksum32 = adler32.Checksum(buf)
				blkptr.Sha256hash = sha256.Sum256(buf)
				blkptr.RawBytes = nil
				blkptr.HasData = false
				out <- blkptr
			}
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
