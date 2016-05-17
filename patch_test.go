// Copyright 2015 Monmohan Singh. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package xferspdy

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

//var logLevel = flag.Int("lv", 3, "log level")

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

func TestFilePatchWordDocument(t *testing.T) {
	fmt.Println("Test to patch a word document")
	fmt.Printf("log level %v\n", *logLevel)
	flag.Lookup("v").Value.Set(fmt.Sprint(*logLevel))
	blksz := 2048

	ofname := "testdata/doc_v1.docx"
	sign := NewFingerprint(ofname, uint32(blksz))
	glog.V(4).Infof("Fingerprint for file %v\n %v\n", ofname, *sign)

	nfname := "testdata/doc_v2.docx"
	delta := NewDiff(nfname, *sign)
	glog.V(4).Infof("Delta = %v ", delta)

	expfname := "/tmp/doc_patched.docx"
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
