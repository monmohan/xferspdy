package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	//"github.com/monmohan/xferspdy/data"
)

func main() {
	testFlag()
	testBufIO()
}

func testFlag() {
	var name string
	flag.StringVar(&name, "name", "unknown", "specify a name")
	flag.Parse()
	fmt.Printf("Hello %s\n", name)

}

func testBufIO() {
	testfile := "/msingh/projects/genknow/gitcheatsheet"
	file, _ := os.Open(testfile)
	bufreader := bufio.NewReader(file)

	data, _ := bufreader.Peek(16)
	fmt.Println(string(data))
	few := make([]byte, 5)
	n, _ := bufreader.Read(few)
	fmt.Printf("number of bytes read %d ", n)
	fmt.Println(string(few))
}
