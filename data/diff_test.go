package data

import (
	"fmt"
	"testing"
)

func TestDeltaSameFile(t *testing.T) {
	fmt.Println("===TestDeltaSameFile===..")
	fname := "../testdata/26bytefile"
	sign := NewSignature(fname, 16)
	fmt.Printf(" %v\n", *sign)
	delta := NewDiff(fname, *sign)
	fmt.Printf("Delta: %v\n", delta)
	fmt.Println("===TestDeltaSameFile END===\n..")

}

func TestDelta2ByteExtraInEnd(t *testing.T) {
	fmt.Println("==TestDelta2ByteExtraInEnd==")
	fname := "../testdata/26bytefile"
	sign := NewSignature(fname, 24)
	fmt.Printf(" %v\n", *sign)
	fname = "../testdata/28bytefile"
	delta := NewDiff(fname, *sign)
	fmt.Printf("Delta: %v\n", delta)
	fmt.Println("==TestDelta2ByteExtraInEnd END==\n")

}
func TestDelta2ByteExtraInMid(t *testing.T) {
	fmt.Println("==TestDelta2ByteExtraInMid, block size 5 ==")
	ofname := "../testdata/10bytefile"
	sign := NewSignature(ofname, 5)
	fmt.Printf(" %v\n", *sign)
	nfname := "../testdata/12bytemidchgfile"
	delta := NewDiff(nfname, *sign)
	fmt.Printf("Delta: %v\n", delta)
	fmt.Println("==TestDelta2ByteExtraInMid block size 8 ==\n")
	sign = NewSignature(ofname, 8)

	delta = NewDiff(nfname, *sign)
	fmt.Printf("Delta: %v\n", delta)

}
