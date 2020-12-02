// Copyright 2015 Monmohan Singh. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package xferspdy

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/golang/glog"
)

func TestFilePatchSimpleText(t *testing.T) {
	otext := []byte(`Go is building a garbage collector (GC) not only for 2015 but for 2025 and beyond: 
		A GC that supports today’s software development and scales along with new software and hardware throughout the next decade. 
		Such a future has no place for stop-the-world GC pauses, which have been an 
		impediment to broader uses of safe and secure languages such as Go.`)
	mtext := []byte(`Go is building a garbage collector (GC) not only for 2015 but for 2025 and beyond: 
		A GC that supports today’s software development and scales along with new software and hardware throughout the next decade. 
		Such a future has no place for stop-the-world GC pauses, which have been an 
		impediment to broader uses of safe and secure languages such as Go.Go 1.5, the first glimpse of this future, 
		achieves GC latencies well below the 10 millisecond goal we set a year ago.`)
	fmt.Printf("log level %v\n", *logLevel)
	flag.Lookup("v").Value.Set(fmt.Sprint(*logLevel))
	blksz := 32
	ofname := "/tmp/TextFilePatchSimple_o"
	ofile, _ := os.OpenFile(ofname, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	ofile.Write(otext)
	ofile.Close()
	sign := NewFingerprint(ofname, uint32(blksz))
	glog.V(4).Infof("Fingerprint for file %v\n %v\n", ofname, *sign)
	nfname := "/tmp/TextFilePatchSimple_1"
	nfile, _ := os.OpenFile(nfname, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	nfile.Write(mtext)
	defer nfile.Close()
	delta := NewDiff(nfname, *sign)
	glog.V(4).Infof("Delta = %v ", delta)

	expfname := "/tmp/TextFilePatchSimple_2"
	expfile, _ := os.OpenFile(expfname, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	defer expfile.Close()
	Patch(delta, *sign, expfile)
	//read from new file and delta and compare file bytes
	nr, _ := os.Open(nfname)
	er, _ := os.Open(expfname)
	nbytes, _ := ioutil.ReadAll(nr)
	expbytes, _ := ioutil.ReadAll(er)
	if !reflect.DeepEqual(expbytes, nbytes) {
		t.Fatalf("bytes don't match after patch nbytes=%v\n exp=%v\n", nbytes, expbytes)
	} else {
		glog.V(4).Infof("bytes match after patch nbytes=%v\n exp=%v\n", nbytes, expbytes)
	}
	glog.Flush()
}

type TestFiles struct {
	baseFile     string
	modifiedFile string
	patchedFile  string
}

func TestPatchManyFiles(t *testing.T) {
	testdata := []TestFiles{
		{"testdata/doc_v1.docx", "testdata/doc_v2.docx", "/tmp/doc_patched.docx"},
		{"testdata/samplepdf.pdf", "testdata/samplepdf_v2.pdf", "/tmp/samplepdf_patched.pdf"},
		{"testdata/sampleimg.jpg", "testdata/sampleimg_v2.jpg", "/tmp/sampleimg_patched.jpg"},
	}
	fmt.Printf("log level %v\n", *logLevel)
	flag.Lookup("v").Value.Set(fmt.Sprint(*logLevel))
	blksz := 2048

	for _, v := range testdata {
		fmt.Printf("Test to patch %s\n", v.baseFile)
		sign := NewFingerprint(v.baseFile, uint32(blksz))
		glog.V(4).Infof("Fingerprint for file %v\n %v\n", v.baseFile, *sign)

		delta := NewDiff(v.modifiedFile, *sign)
		glog.V(4).Infof("Delta = %v ", delta)

		patchedFile, _ := os.OpenFile(v.patchedFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
		defer patchedFile.Close()
		Patch(delta, *sign, patchedFile)
		//read from new file and delta and compare file bytes
		br, _ := os.Open(v.baseFile)
		nr, _ := os.Open(v.modifiedFile)
		er, _ := os.Open(v.patchedFile)
		originalBytes, _ := ioutil.ReadAll(br)
		v2Bytes, _ := ioutil.ReadAll(nr)
		patchedBytes, _ := ioutil.ReadAll(er)
		mustMatch := reflect.DeepEqual(patchedBytes, v2Bytes)
		mustNotMatch := reflect.DeepEqual(patchedBytes, originalBytes)
		if !mustMatch {
			t.Fatalf("Patched Bytes from File %s , don't match the v2 file %s", v.patchedFile, v.modifiedFile)
		}
		if mustNotMatch {
			t.Fatalf("Patched Bytes from File %s , match the v1 file %s", v.patchedFile, v.baseFile)
		}
		fmt.Printf("Matching succeeded base =%s, modified=%s, patched=%s\n", v.baseFile, v.modifiedFile, v.patchedFile)

		glog.Flush()
	}
}
