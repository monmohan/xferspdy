package data

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"reflect"
)

func NewDiff(filename string, sign Signature) []Block {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	finfo, _ := file.Stat()
	r := bufio.NewReader(file)
	delta := make([]Block, 0)
	processBlock(r, 0, finfo.Size(), sign, &delta)
	return delta
}

func processBlock(r *bufio.Reader, rptr int64, filesz int64, s Signature, delta *[]Block) {

	blksz := int64(s.Blocksz)
	brem := int64(blksz)
	if (rptr + blksz) > filesz {
		brem = filesz - rptr
	}
	fmt.Printf("Process Block :rptr %d filesz %d blocksz %d brem %d delta %v\n", rptr, filesz, s.Blocksz, brem, *delta)
	if brem == 0 {
		fmt.Println("All read\n ")
		return
	}

	buf := make([]byte, brem)
	n, err := r.Read(buf)
	if err != nil {
		fmt.Printf("Error %v read %d bytes", err, n)
	}
	//fmt.Printf("Buffer read %v \n", buf)
	checksum, state := Checksum(buf)
	matchblock, matched := matchBlock(checksum, sha256.Sum256(buf), s)
	diff := *delta
	if matched {
		fmt.Printf("Matched block %v \n", matchblock)
		*delta = append(diff, matchblock)
		rptr += int64(brem)
		processBlock(r, rptr, filesz, s, delta)
	} else {
		fmt.Printf("Block not matched\n")
		processRolling(r, state, rptr, filesz, s, delta)
	}

}

func processRolling(r *bufio.Reader, st *State, rptr int64, filesz int64, s Signature, delta *[]Block) {

	db := Block{isdatablock: true, Start: rptr}
	diff := *delta
	brem := filesz - (rptr + int64(len(st.window)))
	fmt.Printf(" Rolling : st %v rptr %d filesz %d blocksz %d brem %d delta %v\n", *st, rptr, filesz, s.Blocksz, brem, *delta)

	if brem == 0 {
		db.data = append(db.data, st.window...)
		*delta = append(diff, db)
		return
	}
	fb := st.window[0]
	db.data = append(db.data, fb)
	b, e := r.ReadByte()
	if e != nil {
		log.Fatal(e)
	}
	rptr += 1
	checksum := st.UpdateWindow(b)
	matchblock, matched := matchBlock(checksum, sha256.Sum256(st.window), s)
	if matched {
		*delta = append(diff, matchblock)
		rptr += int64(len(st.window))
		processBlock(r, rptr, filesz, s, delta)
	} else {
		processRolling(r, st, rptr, filesz, s, delta)
	}
}

func matchBlock(checksum uint32, sha256 [sha256.Size]byte, s Signature) (mblock Block, matched bool) {
	fmt.Printf("comparing input checksum %d ", checksum)
	for _, block := range s.BlockMap {
		//fmt.Printf("Comparing with block %v", block)
		if reflect.DeepEqual(block.Checksum32, checksum) && reflect.DeepEqual(sha256, block.Sha256hash) {
			fmt.Println("found match ")
			return block, true
		}
	}
	return Block{}, false

}
