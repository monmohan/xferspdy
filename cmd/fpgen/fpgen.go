package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/monmohan/xferspdy/data"
	"log"
	"os"
	"path/filepath"
)

var (
	logLevel = flag.Int("loglevel", 3, "log level")
	fPath    = flag.String("file", "", "File path to create the fingerprint, mandatory")
	blockSz  = flag.Uint64("blocksz", 2*1024, "Block Size, default block size is 2KB")
)

func main() {
	flag.Lookup("v").Value.Set(fmt.Sprint(*logLevel))
	flag.Parse()
	if *fPath == "" {
		glog.Fatal("File path is required")
	}
	glog.Infof("File path %s , Block Size %d \n", *fPath, *blockSz)

	fgprt := data.NewSignature(*fPath, uint32(*blockSz))
	glog.V(4).Infof("Signature  %v \n", *fgprt)
	glog.Flush()

	dir, fname := filepath.Split(*fPath)

	fpfile, err := os.OpenFile(filepath.Join(dir, fname+".fingerprint"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	defer fpfile.Close()

	if err != nil {
		log.Fatalf("Error in creating finger print file %v \n, Error :", filepath.Join(dir, fname+".fingerprint"), err)
	}

	enc := gob.NewEncoder(fpfile)
	enc.Encode(*fgprt)
	glog.Info("Signature created %v \n ", fpfile)

}
