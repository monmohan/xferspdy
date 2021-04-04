// +build darwin dragonfly freebsd linux netbsd openbsd

package xferspdy

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/golang/glog"
)

//var fileToUpdate=flag.String("update-file","","File to update some bytes")

func TestRandomFilePatchOld(t *testing.T) {

	runTest(DiffFnOld, PatchFnOld)
}

func TestRandomFilePatchOptimal(t *testing.T) {

	runTest(DiffFnOpitmal, PatchFnOptimal)
}

func runTest(diffFn DiffFunc, patchFn PatchFunc) {
	gob.Register(InsertOp{})
	gob.Register(CopyOp{})
	encodePatch := true
	rand.Seed(time.Now().UnixNano())
	fgprtBlockSize := 1000
	blockSize := 1000
	numBlock := 1000
	fileV1, fileV2, patchFile := "TestRandomFilePatch_v1", "TestRandomFilePatch_v2", "TestRandomFilePatch.patch"
	numBytesToChange := 100
	bytes := make([]byte, numBytesToChange)
	genRandomFile(fileV1, blockSize, numBlock)

	sign := NewFingerprint(fileV1, uint32(fgprtBlockSize))
	glog.V(4).Infof("Fingerprint for file %v\n %v\n", fileV1, *sign)

	//create a V2 file copied from V1 and change few bytes
	v1, _ := os.Open(fileV1)
	defer v1.Close()
	v2, _ := os.Create(fileV2)
	_, err := io.Copy(v2, v1)
	for i := 0; i < numBytesToChange; i++ {
		bytes[i] = 64
	}

	//seekLocations := []int{0, numBytesToChange, rand.Intn(blockSize * numBlock), rand.Intn(blockSize * numBlock), -1}
	seekLocations := []int{rand.Intn(blockSize * numBlock)}

	for i, seekLoc := range seekLocations {
		(func() {
			v2, err = os.OpenFile(fileV2, os.O_RDWR, os.ModePerm)
			defer v2.Close()
			if err != nil {
				log.Fatal(err)
			}
			seek := io.SeekStart
			if seekLoc == -1 {
				seek = io.SeekEnd
				seekLoc = 0
			}
			_, err = v2.Seek(int64(seekLoc), seek)
			if err != nil {
				log.Fatal(err)
			}

			v2.Write(bytes)
		})()
		func() {
			delta := diffFn(fileV2, *sign)
			if encodePatch {
				pfJSON, err := os.Create(fmt.Sprintf("%s_%d.%s", patchFile, i, "json"))
				if err != nil {
					log.Fatal(err.Error())
				}
				enc := json.NewEncoder(pfJSON)
				err = enc.Encode(delta)
				if err != nil {
					log.Fatal(err.Error())
				}
				pfGob, err := os.Create(fmt.Sprintf("%s_%d.%s", patchFile, i, "gob"))
				if err != nil {
					log.Fatal(err.Error())
				}
				enc2 := gob.NewEncoder(pfGob)
				err = enc2.Encode(delta)
				if err != nil {
					log.Fatal(err.Error())
				}

			}
			patchedFile, _ := os.Create(fmt.Sprintf("%s_%d.%s", fileV2, i, "patched"))
			defer patchedFile.Close()
			patchFn(delta, *sign, patchedFile)
			compareBytes(v2.Name(), patchedFile.Name())

		}()

	}

}

func compareBytes(fileV2 string, patchedFile string) {
	v2, err := os.OpenFile(fileV2, os.O_RDONLY, os.ModePerm)
	pf, err := os.OpenFile(patchedFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	if v2_stat, err := v2.Stat(); err != nil {
		pf_stat, _ := pf.Stat()
		if pf_stat.Size() != v2_stat.Size() {
			log.Fatal("failed size match between patched and v2 file")
		}
	}
	buf1 := make([]byte, 1024)
	buf2 := make([]byte, 1024)
	fail := false
	for num, err := v2.Read(buf1); num != 0 && err != io.EOF; {
		glog.V(4).Infof("Read %d bytes from %s, Err=%v\n", num, v2.Name(), err)
		num2, err2 := pf.Read(buf2)
		glog.V(4).Infof("Read %d bytes from %s, Err= %v\n", num2, pf.Name(), err2)
		bytesEqual := bytes.Equal(buf1[0:num], buf2[0:num2])
		if err2 != nil || num2 != num || !bytesEqual {
			fail = true
			break
		}
		num, err = v2.Read(buf1)
	}
	if fail {
		log.Fatal("File comparison failed")
	}
}

func genRandomFile(fileName string, blockSize int, numBlocks int) {
	cmd := exec.Command("dd", "if=/dev/urandom", fmt.Sprintf("of=%s", fileName), fmt.Sprintf("bs=%d", blockSize), fmt.Sprintf("count=%d", numBlocks))
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func TestGobVsJsonOnDisk(t *testing.T) {
	type Dimension interface{}
	type Point struct {
		X, Y, Z int64
	}
	dims := []Dimension{}
	for i := 0; i < 100000; i++ {
		dims = append(dims, Point{int64(i), int64(i * 10), int64(i * 100)})
	}

	pfJSON, err := os.Create(fmt.Sprintf("%s.%s", "TestGobVsJsonOnDisk", "json"))
	if err != nil {
		log.Fatal(err.Error())
	}
	enc := json.NewEncoder(pfJSON)
	err = enc.Encode(dims)
	if err != nil {
		log.Fatal(err.Error())
	}
	gob.Register(Point{})
	pfGob, err := os.Create(fmt.Sprintf("%s.%s", "TestGobVsJsonOnDiskInterface", "gob"))
	if err != nil {
		log.Fatal(err.Error())
	}
	enc2 := gob.NewEncoder(pfGob)
	err = enc2.Encode(dims)
	if err != nil {
		log.Fatal(err.Error())
	}
	points := []Point{}
	for i := 0; i < 100000; i++ {
		points = append(points, Point{int64(i), int64(i * 10), int64(i * 100)})
	}
	pfGob2, err := os.Create(fmt.Sprintf("%s.%s", "TestGobVsJsonOnDiskStruct", "gob"))
	if err != nil {
		log.Fatal(err.Error())
	}
	enc3 := gob.NewEncoder(pfGob2)
	err = enc3.Encode(points)
	if err != nil {
		log.Fatal(err.Error())
	}
}
