package data

import (
	"github.com/golang/glog"
	"io"
	"os"
)

func Patch(delta []Block, sign Fingerprint, t io.Writer) {
	s, e := os.Open(sign.Source)
	defer s.Close()
	wptr := int64(0)
	for _, block := range delta {
		if block.isdatablock {
			glog.V(3).Infof("Writing data block , wptr=%v , num bytes = %v \n", wptr, len(block.data))
			_, e = t.Write(block.data)
			glog.V(4).Infof("Writing bytes = %v \n", block.data)
			if e != nil {
				glog.Fatal(e)
			}
			wptr += int64(len(block.data))
		} else {
			s.Seek(block.Start, 0)
			ds := block.End - block.Start
			glog.V(3).Infof("Writing data block, Block=%v\n , wptr=%v , num bytes = %v \n", block, wptr, ds)
			io.CopyN(t, s, block.End-block.Start)
			wptr += ds
		}
	}
}
