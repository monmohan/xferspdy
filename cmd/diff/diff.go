// Copyright 2015 Monmohan Singh. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/monmohan/xferspdy"
	"os"
	"path/filepath"
)

var (
	fngprt = flag.String("fingerprint", "", "Path to the finger print file of older version, REQUIRED")
	fPath  = flag.String("file", "", "Path to new version of the file to diff with, REQUIRED")
)

func main() {
	flag.Parse()
	if *fPath == "" || *fngprt == "" {
		glog.Fatal("Argument missing")
	}

	glog.V(2).Infof("File path %s , Fingerprint file %s \n", *fPath, *fngprt)

	fpfile, err := os.Open(*fngprt)

	defer fpfile.Close()
	if err != nil {
		glog.Fatalf("Error in reading finger print file %v \n, Error : %s", *fngprt, err)
	}
	var fp xferspdy.Fingerprint
	dec := gob.NewDecoder(fpfile)
	err = dec.Decode(&fp)

	glog.V(4).Infof("Read fingerprint %v \n", fp)

	if err != nil {
		glog.Fatalf("Error in decoding finger print file %v \n, Error : %s", fp, err)
	}

	diff := xferspdy.NewDiff(*fPath, fp)

	dir, fname := filepath.Split(*fPath)

	nfile, err := os.OpenFile(filepath.Join(dir, fname+".patch"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)

	defer nfile.Close()

	if err != nil {
		glog.Fatalf("Error in creating patch file %v \n, Error : %s", filepath.Join(dir, fname+".patch"), err)
	}

	enc := gob.NewEncoder(nfile)
	err = enc.Encode(diff)
	if err != nil {
		glog.Fatalf("Error in encoding Patch file %v \n, Error :%s", filepath.Join(dir, fname+".patch"), err)
	}
	fmt.Printf("Patch file created - %v \n ", nfile.Name())

	glog.Flush()

}
