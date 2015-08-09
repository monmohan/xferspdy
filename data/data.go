package data

import (
	"crypto/sha256"
	"fmt"
	"hash/adler32"
	"io"
	"log"
	"os"
)

type Signature struct {
	Blocksz  uint32
	BlockMap []Block
}

type Block struct {
	Start, End  int64
	Checksum32  uint32
	Sha256hash  [sha256.Size]byte
	isdatablock bool
	data        []byte
}

func (b Block) String() string {
	return fmt.Sprintf("Start %d End %d adler %d data %v\n", b.Start, b.End, b.Checksum32, b.data)
}

func NewSignature(filename string, blocksize uint32) *Signature {
	bufz := make([]byte, blocksize)
	file, e := os.Open(filename)
	defer file.Close()

	if e != nil {
		log.Fatal(e)
	}

	n, start := 0, int64(0)
	var err error = nil
	var block Block
	signature := Signature{Blocksz: blocksize}
	for err == nil {
		n, err = file.Read(bufz)
		//tfmt.Printf("Read file %d bytes read , error= %v \n", n, err)
		if err == nil {
			block = Block{Start: start, End: start + int64(n),
				Checksum32: adler32.Checksum(bufz[0:n]),
				Sha256hash: sha256.Sum256(bufz[0:n])}
			signature.BlockMap = append(signature.BlockMap, block)
			start = block.End
		} else {
			if err == io.EOF {
				fmt.Println("File read complete")
			} else {
				log.Fatal(err)

			}

		}

	}
	return &signature

}
