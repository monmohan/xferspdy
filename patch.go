// Copyright 2015 Monmohan Singh. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xferspdy

import (
	"io"
	"os"

	"github.com/golang/glog"
)

//Patch is a wrapper on PatchFile (current version only supports patching of local files)
func Patch(delta []Block, sign Fingerprint, t io.Writer) {
	PatchFile(delta, sign.Source, t)
}

//Patch is a wrapper on PatchFile (current version only supports patching of local files)
func ApplyPatchOps(delta []PatchOp, sign Fingerprint, t io.Writer) {
	ApplyPatchOpsFile(delta, sign.Source, t)
}

// PatchFile takes a source file and Diff as input, and writes out to the Writer.
// The source file would normally be the base version of the file  and
// the Diff is the delta computed by using the Fingerprint generated for the base file and the new version of the file
func PatchFile(delta []Block, source string, t io.Writer) error {
	s, e := os.Open(source)
	defer s.Close()
	wptr := int64(0)
	for _, block := range delta {
		if block.HasData {
			glog.V(3).Infof("Writing RawBytes block , wptr=%v , num bytes = %v \n", wptr, len(block.RawBytes))
			_, e = t.Write(block.RawBytes)
			glog.V(4).Infof("Writing bytes = %v \n", block.RawBytes)
			if e != nil {
				glog.Fatal(e)
			}
			wptr += int64(len(block.RawBytes))
		} else {
			s.Seek(block.Start, 0)
			ds := block.End - block.Start
			glog.V(3).Infof("Writing RawBytes block, Block=%v\n , wptr=%v , num bytes = %v \n", block, wptr, ds)
			if _, e := io.CopyN(t, s, block.End-block.Start); e != nil {
				return e
			}
			wptr += ds
		}
	}
	return nil
}

// PatchFile takes a source file and Diff as input, and writes out to the Writer.
// The source file would normally be the base version of the file  and
// the Diff is the delta computed by using the Fingerprint generated for the base file and the new version of the file
func ApplyPatchOpsFile(delta []PatchOp, source string, t io.Writer) error {
	s, e := os.Open(source)
	defer s.Close()
	wptr := int64(0)
	for _, op := range delta {
		switch op := op.(type) {

		case InsertOp:
			glog.V(3).Infof("Writing RawBytes block , wptr=%v , num bytes = %v \n", wptr, len(op.Bytes))
			_, e = t.Write(op.Bytes)
			glog.V(4).Infof("Writing bytes = %v \n", op.Bytes)
			if e != nil {
				glog.Fatal(e)
			}
			wptr += int64(len(op.Bytes))

		case CopyOp:
			s.Seek(op.Start, 0)
			ds := op.End - op.Start
			glog.V(3).Infof("Writing RawBytes block, Block=%v\n , wptr=%v , num bytes = %v \n", op, wptr, ds)
			if _, e := io.CopyN(t, s, op.End-op.Start); e != nil {
				return e
			}
			wptr += ds
		}
	}
	return nil
}

type PatchFunc func(interface{}, Fingerprint, io.Writer)
type PatchFileFunc func(interface{}, string, io.Writer)

func PatchFnOptimal(patchOps interface{}, sign Fingerprint, w io.Writer) {
	ApplyPatchOps(patchOps.([]PatchOp), sign, w)
}
func PatchFnOld(blocks interface{}, sign Fingerprint, w io.Writer) {
	Patch(blocks.([]Block), sign, w)
}

func PatchFileOptimal(patchOps interface{}, source string, w io.Writer) {
	ApplyPatchOpsFile(patchOps.([]PatchOp), source, w)
}
func PatchFileOld(blocks interface{}, source string, w io.Writer) {
	PatchFile(blocks.([]Block), source, w)
}
