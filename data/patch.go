package data

import (
	"github.com/golang/glog"
	"io"
	"os"
)

func Patch(delta []Block, sign Fingerprint, t io.Writer) {
	PatchFile(delta, sign.Source, t)
}

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
