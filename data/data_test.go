package data

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
	sign := NewFingerprint("../testdata/Adler32testresource", 2048)
	fmt.Printf(" %v\n", sign.Blocksz)

}

func TestRollingChecksum(t *testing.T) {
	fmt.Println("testing checksum")
	file, e := os.Open("../testdata/samplefile")
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

	num_iter := 100
	st := 3076
	for num_iter > 0 {
		x := data[st:mid]
		libsum := adler32.Checksum(x)
		libroll, state := Checksum(x)
		fmt.Printf("Libsum %d libroll %d \n", libsum, libroll)
		if !reflect.DeepEqual(libsum, libroll) {
			fmt.Printf("Libsum %d libroll %d \n", libsum, libroll)
			t.FailNow()
		}
		st += 1
		x = data[st : mid+1]
		libsum = adler32.Checksum(x)
		libroll = state.UpdateWindow(data[mid])

		if !reflect.DeepEqual(libsum, libroll) {
			fmt.Printf("Libsum %d libroll %d \n", libsum, libroll)
			t.FailNow()
		}
		num_iter -= 1
		mid += 1
	}

}
