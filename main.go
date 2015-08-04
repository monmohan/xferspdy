package main

import (
	"flag"
	"fmt"
	"github.com/monmohan/xferspdy/data"
)

func main() {
	var name string
	flag.StringVar(&name, "name", "unknown", "specify a name")
	flag.Parse()
	fmt.Printf("Hello %s\n", name)
	sign := data.NewSignature("/msingh/projects/genknow/gitcheatsheet", 16)
	fmt.Printf(" %v\n", *sign)
}
