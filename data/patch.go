// Copyright 2015 Monmohan Singh. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//Package data provides the basic interfaces around binary diff and patching process
package data

import (
	"github.com/golang/glog"
	"io"
	"os"
)

//Patch is a wrapper on PatchFile (current version supports patching only files)
func Patch(delta []Block, sign Fingerprint, t io.Writer) {
	PatchFile(delta, sign.Source, t)
}

//PatchFile creates an updated (patched file), given a source file and Diff
//The source to the patch is the base version of the file and its fingerprint
//The diff is the delta computed between the Fingerprint and the new version
func PatchFile(delta []Block, source string, t io.Writer) {
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
			io.CopyN(t, s, block.End-block.Start)
			wptr += ds
		}
	}
}
