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
	Start, End int
	Checksum32 uint32
	Sha256hash [sha256.Size]byte
}

func NewSignature(filename string, blocksize uint32) *Signature {
	bufz := make([]byte, blocksize)
	file, e := os.Open(filename)
	defer file.Close()

	if e != nil {
		log.Fatal(e)
	}

	n, start := 0, 0
	var err error = nil
	var block Block
	signature := Signature{Blocksz: blocksize}
	for err == nil {
		n, err = file.Read(bufz)
		fmt.Printf("Read file %d bytes read , error= %v \n", n, err)
		if err == nil {
			block = Block{Start: start, End: start + n - 1, Checksum32: adler32.Checksum(bufz[0:n]), Sha256hash: sha256.Sum256(bufz[0:n])}
			fmt.Printf("Block %v\n", block)
			signature.BlockMap = append(signature.BlockMap, block)
			start = block.End + 1
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
