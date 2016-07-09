// Copyright 2015 Monmohan Singh. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package xferspdy

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"testing"
)

const (
	network, address, useHTTP, storedir = "tcp", "localhost:9999", false, ".xferspdystore"
)

func TestRpcClient1(t *testing.T) {
	runServer()
	fmt.Println("Simple RPC test with small file")
	fname := "testdata/26bytefile"
	r, _ := os.Open(fname)
	buf, e := ioutil.ReadAll(r)
	client := NewRPCClient(useHTTP, network, address)
	o, e := client.PutObject(PutRequest{Data: buf, Key: "testkey", Blocksize: 8})
	if e != nil {
		fmt.Printf("error ..%s", e)
		t.Fail()
	}
	fmt.Printf("Returned object %v", o)
}

func getStorageDir() string {
	u, _ := user.Current()
	return filepath.Join(u.HomeDir, storedir)

}

func runServer() {

	l, e := net.Listen(network, address)
	if e != nil {
		fmt.Errorf("listen error:", e)
	}
	p := NewProvider(getStorageDir())
	fmt.Printf("Provider %v", *p)
	go ServeRPC(useHTTP, l, p)
}
