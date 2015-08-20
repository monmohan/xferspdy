package data

import (
	"fmt"
	"io"
	"os"
	"reflect"
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
	//t.SkipNow()
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
	if len(sign.BlockMap) != (basesz/blksz)+1 {
		t.Errorf("bad signature, length error %v", len(sign.BlockMap))
		t.FailNow()
	}

	delta := NewDiff(ofname, *sign)

	for i, blk := range delta {

		if blk.Start != sign.BlockMap[i].Start && blk.End != sign.BlockMap[i].End {
			t.Error("failed diff %v \n at blk %v ", delta, blk)
			t.FailNow()
		} else {
			t.Log("Diff and signature block match,\n", blk)
		}
	}
	fmt.Printf("Signature : %v\n", sign)

	fmt.Printf("Delta: %v\n", delta)

}

func TestFewBlocksWithMorebytes(t *testing.T) {
	fmt.Println("==TestFewBlocksWithMorebytes1===")
	blksz := 32
	basesz := 2002
	basefile := "../testdata/samplefile"
	bfile, _ := os.Open(basefile)

	defer bfile.Close()
	ofname := "../testdata/TestFewBlocksWithMorebytes_o"
	ofile, _ := os.OpenFile(ofname, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	io.CopyN(ofile, bfile, int64(basesz))
	bfile.Seek(0, 0)
	ofile.Close()

	sign := NewSignature(ofname, uint32(blksz))

	nfname := "../testdata/TestFewBlocksWithMorebytes_1"
	extraBytes := []byte("xxxx")
	nfile, _ := os.OpenFile(nfname, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	nfile.Write(extraBytes) //append in the begining
	io.CopyN(nfile, bfile, int64(basesz))
	nfile.Write(extraBytes) //append in the end
	nfile.Close()

	delta := NewDiff(nfname, *sign)

	if len(delta) != (len(sign.BlockMap) + 1) {
		t.Fatalf("Error , wrong number of blocks in delta %v\n, signature %v\n", delta, *sign)

	}

	//check first block
	blk := delta[0]
	if !blk.isdatablock || blk.Start != 0 {
		t.Fatalf("First block is not a data block %v \n", blk)
	}
	if !reflect.DeepEqual(extraBytes, blk.data) {
		t.Fatalf("First block extra data mismatch %v \n", blk)
	}

	//check last block
	blk = delta[len(delta)-1]
	lblkSt := len(extraBytes) + basesz - (basesz % blksz)
	fmt.Printf("expected last block start %v\n", lblkSt)
	if !blk.isdatablock || blk.Start != int64(lblkSt) {
		t.Fatalf("Last block is not a data block %v \n", blk)
	}

	delta = delta[1 : len(delta)-1]

	for i, blk := range delta {
		fmt.Printf("Comparing Block number %d , blk %v \n", i, blk)
		if blk.Start != sign.BlockMap[i].Start && blk.End != sign.BlockMap[i].End {
			t.Fatalf("failed diff %v \n at blk %v ", delta, blk)
		}
	}

}
