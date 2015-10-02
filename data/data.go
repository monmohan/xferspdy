package data

import (
	"crypto/sha256"
	"fmt"
	"github.com/golang/glog"
	"hash/adler32"
	"io"
	"os"
)

type Fingerprint struct {
	Blocksz  uint32
	BlockMap []Block
	Source   string
}

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

func NewFingerprint(filename string, blocksize uint32) *Fingerprint {
	bufz := make([]byte, blocksize)
	file, e := os.Open(filename)
	defer file.Close()

	if e != nil {
		glog.Fatal(e)
	}

	n, start := 0, int64(0)
	var err error = nil
	var block Block
	fngprt := Fingerprint{Blocksz: blocksize, Source: filename}

	for err == nil {
		n, err = file.Read(bufz)
		if err == nil {
			block = Block{Start: start, End: start + int64(n),
				Checksum32: adler32.Checksum(bufz[0:n]),
				Sha256hash: sha256.Sum256(bufz[0:n])}
			fngprt.BlockMap = append(fngprt.BlockMap, block)
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
