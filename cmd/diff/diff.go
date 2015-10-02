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
	if *fPath == "" || *fngprt == "" {
		glog.Fatal("Argument missing")
	}

	glog.V(2).Infof("File path %s , Fingerprint file %s \n", *fPath, *fngprt)

	fpfile, err := os.Open(*fngprt)

	defer fpfile.Close()
	if err != nil {
		glog.Fatalf("Error in reading finger print file %v \n, Error :", *fngprt, err)
	}
	var fp data.Fingerprint
	dec := gob.NewDecoder(fpfile)
	err = dec.Decode(&fp)

	glog.V(4).Infof("Read fingerprint %v \n", fp)

	if err != nil {
		glog.Fatalf("Error in decoding finger print file %v \n, Error :", fp, err)
	}

	diff := data.NewDiff(*fPath, fp)

	dir, fname := filepath.Split(*fPath)

	nfile, err := os.OpenFile(filepath.Join(dir, fname+".patch"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)

	defer nfile.Close()

	if err != nil {
		glog.Fatalf("Error in creating patch file %v \n, Error :", filepath.Join(dir, fname+".patch"), err)
	}

	enc := gob.NewEncoder(nfile)
	err = enc.Encode(diff)
	if err != nil {
		glog.Fatalf("Error in encoding Patch file %v \n, Error :", filepath.Join(dir, fname+".patch"), err)
	}
	glog.V(2).Infof("Patch created %v \n ", nfile.Name())

	glog.Flush()

}
