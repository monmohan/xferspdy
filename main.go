package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"os"
)

func main() {
	//testFlag()
	var cline = flag.NewFlagSet(os.Args, flag.ExitOnError)
	fmt.Printf("command line value %v\n", cline)
	//testBufIO()
	//flag.Parse()
	//testLog()
}

func testFlag() {
	var name string
	flag.StringVar(&name, "name", "unknown", "specify a name")
	flag.Parse()
	fmt.Printf("Hello %s\n", name)

}

func testLog() {
	r := []byte("xxxx")
	glog.V(2).Infof("Something logged %v", r)

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
