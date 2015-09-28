package main

import (
	"encoding/gob"
	"flag"
	"github.com/golang/glog"
	"github.com/monmohan/xferspdy/data"
	"os"
	"path/filepath"
)

var (
	fngprt = flag.String("fingerprint", "", "Path to the finger print file of older version")
	fPath  = flag.String("file", "", "Path to new version of the file to diff with")
)

func main() {
	flag.Parse()
	if *fPath == "" {
		glog.Fatal("File path to generate diff is required")
	}
	if *fngprt == "" {
		glog.Fatal("Fingerprint file path, of the older version, is required")
	}

	glog.V(2).Infof("File path %s , Fingerprint file %s \n", *fPath, *fngprt)
	fpfile, err := os.Open(*fngprt)
	defer fpfile.Close()
	if err != nil {
		glog.Fatalf("Error in reading finger print file %v \n, Error :", fngprt, err)
	}
	var fp data.Fingerprint
	dec := gob.NewDecoder(fpfile)
	err = dec.Decode(&fp)

	if err != nil {
		glog.Fatalf("Error in decoding finger print file %v \n, Error :", fngprt, err)
	}

	diff := data.NewDiff(*fPath, fp)

	dir, fname := filepath.Split(*fPath)

	nfile, err := os.OpenFile(filepath.Join(dir, fname+".diff"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)

	defer nfile.Close()

	if err != nil {
		glog.Fatalf("Error in creating diff file %v \n, Error :", filepath.Join(dir, fname+".diff"), err)
	}

	enc := gob.NewEncoder(nfile)
	err = enc.Encode(diff)
	if err != nil {
		glog.Fatalf("Error in encoding diff file %v \n, Error :", filepath.Join(dir, fname+".diff"), err)
	}
	glog.V(2).Infof("Diff created %v \n ", nfile.Name())

	glog.Flush()

}
