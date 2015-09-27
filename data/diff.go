package data

import (
	"crypto/sha256"
	"github.com/golang/glog"
	"io"
	"os"
	"reflect"
)

func NewDiff(filename string, sign Fingerprint) []Block {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		glog.Fatal(err)
	}

	finfo, _ := file.Stat()

	var delta []Block
	processBlock(file, 0, finfo.Size(), sign, &delta)
	return delta
}

func processBlock(r io.Reader, rptr int64, filesz int64, s Fingerprint, delta *[]Block) {

	blksz := int64(s.Blocksz)
	brem := int64(blksz)
	if (rptr + blksz) > filesz {
		brem = filesz - rptr
	}
	glog.V(2).Infof("Process Block :rptr %d filesz %d blocksz %d brem %d \n", rptr, filesz, s.Blocksz, brem)
	glog.V(4).Infof("Delta %v \n", *delta)
	if brem == 0 {
		glog.V(2).Infof("All read\n ")
		return
	}

	buf := make([]byte, brem)
	n, err := io.ReadFull(r, buf)
	if err != nil || int64(n) != brem {
		glog.Fatalf("Error %v read %d bytes", err, n)
	}

	checksum, state := Checksum(buf)
	matchblock, matched := matchBlock(checksum, sha256.Sum256(buf), s)
	if matched {
		glog.V(2).Infof("Matched block %v \n", matchblock)
		*delta = append(*delta, matchblock)
		rptr += int64(brem)
		processBlock(r, rptr, filesz, s, delta)
	} else {
		glog.V(2).Infof("Block not matched\n")
		*delta = append(*delta, Block{isdatablock: true, Start: rptr})
		processRolling(r, state, rptr, filesz, s, delta)
	}

}

func processRolling(r io.Reader, st *State, rptr int64, filesz int64, s Fingerprint, delta *[]Block) {

	diff := *delta
	db := &diff[len(diff)-1]
	glog.V(4).Infof("db.data %v \n", db)
	brem := filesz - (rptr + int64(len(st.window)))
	glog.V(4).Infof("Rolling State: State %v \n", *st)
	glog.V(3).Infof("Rolling Info: rptr %d filesz %d blocksz %d brem %d \n", rptr, filesz, s.Blocksz, brem)
	glog.V(4).Infof("Delta %v \n", *delta)

	if brem == 0 {
		db.data = append(db.data, st.window...)
		*delta = diff
		glog.V(4).Infof("db.data %v \n", db.data)
		return
	}
	fb := st.window[0]
	db.data = append(db.data, fb)
	b := make([]byte, 1)
	_, e := io.ReadFull(r, b)
	if e != nil {
		glog.Fatal(e)
	}
	rptr += 1
	checksum := st.UpdateWindow(b[0])
	matchblock, matched := matchBlock(checksum, sha256.Sum256(st.window), s)
	if matched {
		*delta = append(diff, matchblock)
		rptr += int64(len(st.window))
		processBlock(r, rptr, filesz, s, delta)
	} else {
		processRolling(r, st, rptr, filesz, s, delta)
	}
}

func matchBlock(checksum uint32, sha256 [sha256.Size]byte, s Fingerprint) (mblock Block, matched bool) {
	glog.V(3).Infof("comparing input checksum %d ", checksum)
	for _, block := range s.BlockMap {

		if reflect.DeepEqual(block.Checksum32, checksum) && reflect.DeepEqual(sha256, block.Sha256hash) {
			glog.V(2).Infof("found match ")
			return block, true
		}
	}
	return Block{}, false

}
