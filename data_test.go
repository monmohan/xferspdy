// Copyright 2015 Monmohan Singh. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package xferspdy

import (
	"fmt"
	"hash/adler32"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
)

func TestFingerprintCreate(t *testing.T) {
	//t.Skip("not now..")
	sign := NewFingerprint("testdata/Adler32testresource", 2048)
	fmt.Printf(" %v\n", sign.Blocksz)

}

func TestRollingChecksum(t *testing.T) {
	fmt.Println("testing checksum")
	file, e := os.Open("testdata/samplefile")
	defer file.Close()

	if e != nil {
		log.Fatal(e)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	mid := 5000
	//mid = 1100

	numIter := 100
	st := 3076
	for numIter > 0 {
		x := data[st:mid]
		libsum := adler32.Checksum(x)
		libroll, state := Checksum(x)
		fmt.Printf("Libsum %d libroll %d \n", libsum, libroll)
		if !reflect.DeepEqual(libsum, libroll) {
			fmt.Printf("Libsum %d libroll %d \n", libsum, libroll)
			t.FailNow()
		}
		st++
		x = data[st : mid+1]
		libsum = adler32.Checksum(x)
		libroll = state.UpdateWindow(data[mid])

		if !reflect.DeepEqual(libsum, libroll) {
			fmt.Printf("Libsum %d libroll %d \n", libsum, libroll)
			t.FailNow()
		}
		numIter--
		mid++
	}

}

func Example() {
	//Create fingerprint of a file
	fingerprint := NewFingerprint("/path/foo_v1.binary", 1024)

	//Say the file was updated
	//Lets generate the diff
	diff := NewDiff("/path/foo_v2.binary", *fingerprint)

	//diff is sufficient to recover/recreate the modified file, given the base/source and the diff.
	modifiedFile, _ := os.OpenFile("/path/foo_v2_from_v1.binary", os.O_CREATE|os.O_WRONLY, 0777)

	//This writes the output to modifiedFile (Writer). The result will be the same binary as /path/foo_v2.binary
	PatchFile(diff, "/path/foo_v1.binary", modifiedFile)

}
