package data

import (
	"fmt"
	"testing"
)

func TestDeltaSameFile(t *testing.T) {
	t.Skip("dd")
	fmt.Println("Testing diff between files with same binary..")
	fname := "../testdata/26bytefile"
	sign := NewSignature(fname, 16)
	fmt.Printf(" %v\n", *sign)
	delta := NewDiff(fname, *sign)
	fmt.Printf("Delta: %v\n", delta)

}

func TestDelta2ByteExtraInEnd(t *testing.T) {
	fmt.Println("Testing diff between files with one having 2 bytes extra in the end..")
	fname := "../testdata/26bytefile"
	sign := NewSignature(fname, 24)
	fmt.Printf(" %v\n", *sign)
	fname = "../testdata/28bytefile"
	delta := NewDiff(fname, *sign)
	fmt.Printf("Delta: %v\n", delta)

}
