package data

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"io"
	"os"
	"reflect"
	"testing"
)

var logLevel = flag.Int("lv", 3, "log level")

func TestDeltaSameFile(t *testing.T) {
	fmt.Println("===TestDeltaSameFile===..")
	fname := "../testdata/26bytefile"
	sign := NewFingerprint(fname, 16)
	fmt.Printf(" %v\n", *sign)
	delta := NewDiff(fname, *sign)
	fmt.Printf("Delta: %v\n", delta)
	fmt.Println("===TestDeltaSameFile END===\n..")

}

func TestDelta2ByteExtraInEnd(t *testing.T) {
	fmt.Println("==TestDelta2ByteExtraInEnd==")
	fname := "../testdata/26bytefile"
	sign := NewFingerprint(fname, 24)
	fmt.Printf(" %v\n", *sign)
	fname = "../testdata/28bytefile"
	delta := NewDiff(fname, *sign)
	fmt.Printf("Delta: %v\n", delta)
	fmt.Println("==TestDelta2ByteExtraInEnd END==\n")

}
func TestDelta2ByteExtraInMid(t *testing.T) {
	fmt.Println("==TestDelta2ByteExtraInMid, block size 5 ==")
	ofname := "../testdata/10bytefile"
	sign := NewFingerprint(ofname, 5)
	fmt.Printf(" %v\n", *sign)
	nfname := "../testdata/12bytemidchgfile"
	delta := NewDiff(nfname, *sign)
	fmt.Printf("Delta: %v\n", delta)
	fmt.Println("==TestDelta2ByteExtraInMid block size 8 ==\n")
	sign = NewFingerprint(ofname, 8)

	delta = NewDiff(nfname, *sign)
	fmt.Printf("Delta: %v\n", delta)

}

func TestSameBlocks(t *testing.T) {
	fmt.Println("==TestSameBlocks==")

	blksz := 1024
	basesz := 10000
	basefile := "../testdata/samplefile"
	bfile, _ := os.Open(basefile)
    defer bfile.Close()

	ofname := "/tmp/TestSameBlocks"
	ofile, _ := os.OpenFile(ofname, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	io.CopyN(ofile, bfile, int64(basesz))
	ofile.Close()
	
	
	sign := NewFingerprint(ofname, uint32(blksz))
	
	delta := NewDiff(ofname, *sign)

	for _, blk := range delta {

		if b, ok := matchBlock(blk.Checksum32, blk.Sha256hash, *sign); !(ok && (b.Start == blk.Start && b.End == blk.End)) {
			t.Errorf("failed diff %v \n at blk %v ", delta, blk)
			t.FailNow()
		}
		t.Log("Diff and Fingerprint block match,\n", blk)

	}
	
	fmt.Printf("Sign:=>	 %v\n",sign)
	
	fmt.Printf("Delta: %v\n", delta)

}

func TestFewBlocksWithMorebytes(t *testing.T) {
	fmt.Println("==TestFewBlocksWithMorebytes, add bytes in begin and end of file===\n")
	fmt.Printf("log level %v\n", *logLevel)
	flag.Lookup("v").Value.Set(fmt.Sprint(*logLevel))

	fmt.Println("log v value ", flag.Lookup("v").Value)
	blksz := 64 * 1024
	basesz := 200000
	basefile := "../testdata/samplefile"
	bfile, _ := os.Open(basefile)

	defer bfile.Close()
	ofname := "/tmp/TestFewBlocksWithMorebytes_o"
	ofile, _ := os.OpenFile(ofname, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	io.CopyN(ofile, bfile, int64(basesz))
	bfile.Seek(0, 0)
	ofile.Close()

	sign := NewFingerprint(ofname, uint32(blksz))
	glog.V(4).Infof("Fingerprint for file %v\n %v\n", ofname, *sign)

	nfname := "/tmp/TestFewBlocksWithMorebytes_1"
	extraBytes := []byte("xxxx")
	nfile, _ := os.OpenFile(nfname, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	nfile.Write(extraBytes) //append in the begining
	io.CopyN(nfile, bfile, int64(basesz))
	nfile.Write(extraBytes) //append in the end
	nfile.Close()

	delta := NewDiff(nfname, *sign)
	glog.V(2).Infof("Resulting Delta %v\n", delta)
	additionalblks := 1
	if basesz%blksz == 0 {
		additionalblks = 2
	}

	if len(delta) != (len(sign.BlockMap) + additionalblks) {
		t.Fatalf("Error , wrong number of blocks in delta %v\n, Fingerprint %v\n", delta, *sign)

	}

	//check first block
	blk := delta[0]
	if !blk.HasData || blk.Start != 0 {
		t.Fatalf("First block is not a RawBytes block %v \n", blk)
	}
	if !reflect.DeepEqual(extraBytes, blk.RawBytes) {
		t.Fatalf("First block extra RawBytes mismatch %v \n", blk)
	}

	//check last block
	blk = delta[len(delta)-1]
	lblkSt := len(extraBytes) + basesz - (basesz % blksz)
	glog.V(0).Infof("expected last block start %v\n", lblkSt)
	if !blk.HasData || blk.Start != int64(lblkSt) {
		t.Fatalf("Last block is not a RawBytes block %v \n", blk)
	}

	delta = delta[1 : len(delta)-1]

	for i, blk := range delta {
		glog.V(0).Infof("Comparing Block number %d , blk %v \n", i, blk)
		_, matched := matchBlock(blk.Checksum32, blk.Sha256hash, *sign)
		if !matched {
			t.Fatalf("Failed, delta block doesn't match %v \n", blk)
		}
	}
	glog.Flush()
}

func TestFirstLastBlockDataDeleted(t *testing.T) {
	fmt.Println("==TestFirstLastBlockDataDeleted===\n")
	fmt.Printf("log level %v\n", *logLevel)
	flag.Lookup("v").Value.Set(fmt.Sprint(*logLevel))

	blksz := 1024
	basesz := 200000
	delBytes := make([]byte, 1000)

	basefile := "../testdata/samplefile"
	bfile, _ := os.Open(basefile)

	defer bfile.Close()
	ofname := "/tmp/TestFirstLastBlockDataDeleted_o"
	ofile, _ := os.OpenFile(ofname, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	io.CopyN(ofile, bfile, int64(basesz))
	bfile.Seek(0, 0)
	ofile.Close()

	sign := NewFingerprint(ofname, uint32(blksz))
	glog.V(4).Infof("Fingerprint for file %v\n %v\n", ofname, *sign)

	nfname := "/tmp/TestFirstLastBlockDataDeleted_1"

	nfile, _ := os.OpenFile(nfname, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	//move read pointer
	io.ReadFull(bfile, delBytes)
	//read fewer bytes
	io.CopyN(nfile, bfile, int64(basesz-(2*len(delBytes))))
	nfile.Close()

	delta := NewDiff(nfname, *sign)
	glog.V(2).Infof("Resulting Delta %v\n", delta)
	glog.Flush()
	additionalblks := -1
	if rem := basesz % blksz; rem > len(delBytes) {
		additionalblks = 0
	}

	if len(delta) != (len(sign.BlockMap) + additionalblks) {
		t.Fatalf("Error , wrong number of blocks in delta %v\n, Fingerprint %v\n", delta, *sign)

	}

	//check first block
	blk := delta[0]
	if !blk.HasData || blk.Start != 0 {
		t.Fatalf("First block is not a RawBytes block %v \n", blk)
	}

	//check last block
	blk = delta[len(delta)-1]
	lastBlockIsDatablk := ((basesz - len(delBytes)) % blksz) != 0
	if lastBlockIsDatablk != blk.HasData {
		t.Fatalf("Last block is has a wrong block type , expected RawBytes block %v\n, Block %v \n", lastBlockIsDatablk, blk)
	}

	delta = delta[1 : len(delta)-1]

	for i, blk := range delta {
		glog.V(0).Infof("Comparing Block number %d , blk %v \n", i, blk)
		_, matched := matchBlock(blk.Checksum32, blk.Sha256hash, *sign)
		if !matched {
			t.Fatalf("Failed, delta block doesn't match %v \n", blk)
		}
	}

}
func TestRandomChanges(t *testing.T) {
	fmt.Println("==TestRandomChanges===\n")
	fmt.Printf("log level %v\n", *logLevel)
	flag.Lookup("v").Value.Set(fmt.Sprint(*logLevel))

	blksz := 1024
	basesz := 200000

	basefile := "../testdata/samplefile"
	bfile, _ := os.Open(basefile)

	defer bfile.Close()
	ofname := "/tmp/TestRandomChanges_o"
	ofile, _ := os.OpenFile(ofname, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	io.CopyN(ofile, bfile, int64(basesz))
	bfile.Seek(0, 0)
	ofile.Close()

	sign := NewFingerprint(ofname, uint32(blksz))
	glog.V(4).Infof("Fingerprint for file %v\n %v\n", ofname, *sign)

	nfname := "/tmp/TestRandomChanges_1"
	totalBlks := len(sign.BlockMap)
	if totalBlks < 4 {
		t.Fatal("number of blocks should be atleast 4")
	}
	nfile, _ := os.OpenFile(nfname, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	buf := make([]byte, blksz)
	io.ReadFull(bfile, buf)
	//drop first few bytes
	dropBytes := 4
	_, err := nfile.Write(buf[dropBytes:])
	if err != nil {
		t.Fatalf("write failed %v", err)
	}
	//copy couple of blocks
	io.CopyN(nfile, bfile, int64(2*blksz))
	//read a block in mem and change last couple of bytes
	io.ReadFull(bfile, buf)
	buf[len(buf)-1] = buf[len(buf)-1] + 1
	buf[len(buf)-2] = buf[len(buf)-2] + 1
	buf = append(buf, 0) //append random byte
	_, err = nfile.Write(buf)
	if err != nil {
		t.Fatalf("write failed %v", err)
	}
	//copy some more
	//copy couple of blocks
	io.CopyN(nfile, bfile, int64(blksz*(totalBlks-4)))

	//add one more block
	io.CopyN(nfile, bfile, int64(blksz))
	nfile.Close()

	delta := NewDiff(nfname, *sign)
	glog.V(2).Infof("Resulting Delta %v\n", delta)
	glog.Flush()

	if !delta[0].HasData || !delta[len(delta)-1].HasData || !delta[3].HasData {
		t.Fatalf(" First/last/3rd block is not a RawBytes block %v \n", delta)
	}

	/*if len(delta) != (len(sign.BlockMap) + 1) {
		t.Fatalf("wrong delta size , delta=%v\n, Fingerprint=%v\n", delta, sign)
	}*/

}
