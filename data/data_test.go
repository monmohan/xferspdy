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

func TestSignatureCreate(t *testing.T) {
	t.Skip("not now..")
	sign := NewSignature("/msingh/projects/genknow/gitcheatsheet", 16)
	fmt.Printf(" %v\n", *sign)

}

func TestRollingChecksum(t *testing.T) {
	fmt.Println("testing checksum")
	file, e := os.Open("/msingh/projects/gocode/testdata/Adler32testresource")
	defer file.Close()

	if e != nil {
		log.Fatal(e)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	mid := len(data) >> 1
	mid = 1100

	num_iter := 1
	st := 0
	for num_iter > 0 {
		x := data[st:mid]
		libsum := adler32.Checksum(x)
		libroll, state := Checksum(x)

		if !reflect.DeepEqual(libsum, libroll) {
			fmt.Printf("Libsum %d libroll %d \n", libsum, libroll)
			t.Fail()
		}
		st += 1
		x = data[st : mid+1]
		libsum = adler32.Checksum(x)
		libroll = state.UpdateWindow(data[mid])

		if !reflect.DeepEqual(libsum, libroll) {
			fmt.Printf("Libsum %d libroll %d \n", libsum, libroll)
			t.Fail()
		}
		num_iter -= 1
		mid += 1
	}

}
