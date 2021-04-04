package xferspdy

import (
	"crypto/sha256"
	"io"
	"os"

	"github.com/golang/glog"
)

type PatchOp interface{}

// NewDiff computes a diff between a given file and Fingerprint created from some other file
// The diff is represented as a slice of Blocks. Matching Blocks are represented just by their hashes, start and end byte position
// Non-matching blocks are raw binary arrays.
func NewDiffOptimal(filename string, sign Fingerprint) []PatchOp {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		glog.Fatal(err)
	}

	finfo, _ := file.Stat()

	var delta []PatchOp
	processDiffV2(file, finfo.Size(), sign, &delta)
	glog.V(3).Infof("Delta created %v\n", delta)
	return delta
}

type CopyOp struct {
	Start, End int64
}
type InsertOp struct {
	Bytes []byte
}

func processBlockV2(r io.Reader, rptr int64, filesz int64, s Fingerprint, delta *[]PatchOp) processingResult {

	blksz := int64(s.Blocksz)
	brem := int64(blksz)
	if (rptr + blksz) > filesz {
		brem = filesz - rptr
	}
	glog.V(2).Infof("Process Block :rptr %d filesz %d blocksz %d brem %d \n", rptr, filesz, s.Blocksz, brem)
	glog.V(4).Infof("Delta %v \n", *delta)
	if brem == 0 {
		glog.V(2).Infof("All read\n ")
		return processingResult{false, Block{}, nil, rptr, true}
	}

	buf := make([]byte, brem)
	n, err := io.ReadFull(r, buf)
	if err != nil || int64(n) != brem {
		glog.Fatalf("Error %v read %d bytes", err, n)
	}

	checksum, state := Checksum(buf)
	matchblock, matched := matchBlock(checksum, sha256.Sum256(buf), s)
	return processingResult{matched, matchblock, state, rptr, false}

}

func processRollingV2(r io.Reader, st *State, rptr int64, filesz int64, s Fingerprint, delta *[]PatchOp) processingResult {
	diff := *delta
	db := (diff[len(diff)-1]).(Block)
	glog.V(4).Infof("db.RawBytes %v \n", db)
	brem := filesz - (rptr + int64(len(st.window)))
	glog.V(4).Infof("Rolling State: State %v \n", *st)
	glog.V(3).Infof("Rolling Info: rptr %d filesz %d blocksz %d brem %d \n", rptr, filesz, s.Blocksz, brem)
	glog.V(4).Infof("Delta %v \n", *delta)

	if brem == 0 {
		db.RawBytes = append(db.RawBytes, st.window...)
		diff[len(diff)-1] = InsertOp{Bytes: db.RawBytes}
		glog.V(4).Infof("db.RawBytes %v \n", db.RawBytes)
		return processingResult{false, Block{}, nil, rptr, true}
	}
	fb := st.window[0]
	db.RawBytes = append(db.RawBytes, fb)
	diff[len(diff)-1] = db
	b := make([]byte, 1)
	_, e := io.ReadFull(r, b)
	if e != nil {
		glog.Fatal(e)
	}
	rptr++
	checksum := st.UpdateWindow(b[0])
	matchblock, matched := matchBlock(checksum, sha256.Sum256(st.window), s)
	return processingResult{matched, matchblock, st, rptr, false}
}

func processDiffV2(r io.Reader, filesz int64, s Fingerprint, delta *[]PatchOp) {

	var (
		state     *State
		rptr      int64
		result    processingResult
		blockMode bool
	)
	blockMode = true
	for {
		if blockMode {
			result = processBlockV2(r, rptr, filesz, s, delta)
			rptr = result.readPtr
			state = result.windowState
			if result.eof {
				return
			}
			if result.blockMatch {
				*delta = append(*delta, CopyOp{Start: result.matchedBlock.Start, End: result.matchedBlock.End})
				rptr += int64(len(state.window))
				continue
			}
			glog.V(3).Infof("Block not matched\n")
			*delta = append(*delta, Block{HasData: true, Start: rptr})
			blockMode = false
		}
		result = processRollingV2(r, state, rptr, filesz, s, delta)
		rptr = result.readPtr
		state = result.windowState

		if result.eof {
			return
		}
		if result.blockMatch {
			//Last Block is InsertOp now
			diff := *delta
			lastRollingBlock := (diff[len(diff)-1]).(Block)
			diff[len(diff)-1] = InsertOp{Bytes: lastRollingBlock.RawBytes}
			*delta = append(*delta, CopyOp{Start: result.matchedBlock.Start, End: result.matchedBlock.End})
			rptr += int64(len(state.window))
			blockMode = true
			continue
		}

	}

}

type DiffFunc func(string, Fingerprint) interface{}

func DiffFnOld(fileName string, sign Fingerprint) interface{} {
	return NewDiff(fileName, sign)
}

func DiffFnOpitmal(fileName string, sign Fingerprint) interface{} {
	return NewDiffOptimal(fileName, sign)
}
