package data

import (
	"fmt"
	"io"
	"os"
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

var alphabets = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func TestSameBlocks(t *testing.T) {
	fmt.Println("==TestSameBlocks==")
	blksz := 32
	basesz := 1000
	basefile := "../testdata/samplefile"
	bfile, _ := os.Open(basefile)

	defer bfile.Close()
	ofname := "../testdata/TestSameBlocks"
	ofile, _ := os.OpenFile(ofname, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 777)
	io.CopyN(ofile, bfile, int64(basesz))
	ofile.Close()

	sign := NewSignature(ofname, uint32(blksz))
	fmt.Printf("Signature : %v\n", sign)
	if len(sign.BlockMap) != (basesz/blksz)+1 {
		t.Errorf("bad signature, length error %v", len(sign.BlockMap))
		t.FailNow()
	}

	delta := NewDiff(ofname, *sign)

	for i, blk := range delta {

		if blk.Start != sign.BlockMap[i].Start && blk.End != sign.BlockMap[i].End {
			t.Error("failed diff %v \n at blk %v ", delta, blk)
		} else {
			t.Log("Diff and signature block match,\n", blk)
		}
	}
	fmt.Printf("Delta: %v\n", delta)

}
