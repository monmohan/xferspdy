// Copyright 2015 Monmohan Singh. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//Package xferspdy provides the basic interfaces around binary diff and patching process
package xferspdy

import (
	"crypto/sha256"
	"fmt"
	"github.com/golang/glog"
	"hash/adler32"
	"io"
	"os"
)

//Fingerprint of a given File.
//Read file as a set of Blocks
//Calculate an Adler-32 and SHA256 hashes for the Block.
//Creata a lookup map from the hashes to the Block
type Fingerprint struct {
	Blocksz  uint32
	BlockMap map[uint32]map[[sha256.Size]byte]Block
	Source   string
}

//Block represent a byte array chunk
//For each block, following are computed
//Adler-32 and SHA256 checksum
//Start and End byte pos of the block
//Whether or not its a data block
//If this is a data block, RawBytes will capture the byte data
//represented by this block
type Block struct {
	Start, End int64
	Checksum32 uint32
	Sha256hash [sha256.Size]byte
	HasData    bool
	RawBytes   []byte
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

//NewFingerprint computes the fingerprint for a file and given blocksize
func NewFingerprint(filename string, blocksize uint32) *Fingerprint {
	bufz := make([]byte, blocksize)
	file, e := os.Open(filename)
	defer file.Close()

	if e != nil {
		glog.Fatal(e)
	}

	n, start := 0, int64(0)

	var (
		err   error
		block Block
	)

	fngprt := Fingerprint{
		Blocksz: blocksize, BlockMap: make(map[uint32]map[[sha256.Size]byte]Block), Source: filename}

	for err == nil {
		n, err = file.Read(bufz)
		if err == nil {
			block = Block{Start: start, End: start + int64(n),
				Checksum32: adler32.Checksum(bufz[0:n]),
				Sha256hash: sha256.Sum256(bufz[0:n])}
			addBlock(&fngprt, block)
			start = block.End
		} else {
			if err == io.EOF {
				glog.V(2).Infoln("File read complete")
			} else {
				glog.Fatal(err)
			}

		}

	}
	return &fngprt

}

func addBlock(f *Fingerprint, b Block) {
	glog.V(2).Infof("Adding Block %v ", b)
	if sha2blk := f.BlockMap[b.Checksum32]; sha2blk == nil {
		f.BlockMap[b.Checksum32] = make(map[[sha256.Size]byte]Block)
	}
	f.BlockMap[b.Checksum32][b.Sha256hash] = b

}
