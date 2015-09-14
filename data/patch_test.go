package data

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"os"
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
	ofname := "../testdata/TextFilePatchSimple_o"
	ofile, _ := os.OpenFile(ofname, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	ofile.Write(otext)
	ofile.Close()
	sign := NewSignature(ofname, uint32(blksz))
	glog.V(4).Infof("Signature for file %v\n %v\n", ofname, *sign)
	nfname := "../testdata/TextFilePatchSimple_1"
	nfile, _ := os.OpenFile(nfname, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	nfile.Write(mtext)
	nfile.Close()
	delta := NewDiff(nfname, *sign)

	expfname := "../testdata/TextFilePatchSimple_2"
	expfile, _ := os.OpenFile(expfname, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	defer expfile.Close()
	Patch(delta, *sign, expfile)

}
