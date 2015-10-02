package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/monmohan/xferspdy/data"
	"os"
	"path/filepath"
)

var (
	fPatch = flag.String("patch", "", "Path to the patch file ")
	fPath  = flag.String("base", "", "Path to older version of the file to apply patch on")
)

func main() {
	flag.Parse()
	if *fPath == "" || *fPatch == "" {
		glog.Fatal("Argument missing")
	}

	glog.V(2).Infof("File path %s , Fingerprint file %s \n", *fPath, *fPatch)

	pf, err := os.Open(*fPatch)

	defer pf.Close()
	if err != nil {
		glog.Fatalf("Error in reading patch file %v \n, Error :", *fPatch, err)
	}

	var pd []data.Block
	dec := gob.NewDecoder(pf)
	err = dec.Decode(&pd)

	glog.V(4).Infof("Patch file read %v \n", pd)

	if err != nil {
		glog.Fatalf("Error in decoding patch file %v \n, Error :", *fPatch, err)
	}

	dir, fname := filepath.Split(*fPath)

	target, err := os.OpenFile(filepath.Join(dir, "Patched_"+fname), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)

	defer target.Close()

	if err != nil {
		glog.Fatalf("Error in applying patch  %v \n, Error :", filepath.Join(dir, fname+".patched"), err)
	}

	data.PatchFile(pd, *fPath, target)

	fmt.Printf("Patch applied, Target file generated - %v \n ", target.Name())

	glog.Flush()

}
