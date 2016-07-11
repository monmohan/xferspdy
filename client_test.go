// Copyright 2015 Monmohan Singh. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package xferspdy

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"io/ioutil"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"testing"
)

const (
	network, address, useHTTP, storedir = "tcp", "localhost:9999", false, ".xferspdystore"
)

func TestClientPutRequest(t *testing.T) {

	fmt.Println("Test RPC put..")
	fmt.Printf("log level %v\n", *logLevel)
	flag.Lookup("v").Value.Set(fmt.Sprint(*logLevel))
	fname := "testdata/26bytefile"
	r, _ := os.Open(fname)
	buf, e := ioutil.ReadAll(r)
	client := NewRPCClient(useHTTP, network, address)
	o, e := client.PutObject(PutRequest{Data: buf, Key: "TestClientPutRequestKey", Blocksize: 8})
	if e != nil {
		fmt.Printf("error ..%s", e)
		t.Fail()
	}
	glog.V(4).Infof("Generated fingerprint %v", o.Fingerprint)
	fo := NewFingerprint(fname, 8)
	fo.Source = "TestClientPutRequestKey"
	if !o.Fingerprint.DeepEqual(fo) {
		t.Fail()
	}
}

func TestClientPatchRequest(t *testing.T) {

	fmt.Println("Test RPC patch a word document")
	fmt.Printf("log level %v\n", *logLevel)
	flag.Lookup("v").Value.Set(fmt.Sprint(*logLevel))

	fmt.Println("log v value ", flag.Lookup("v").Value)
	blksz := 2048
	key := "TestClientPatchRequestKey"

	ofname := "testdata/doc_v1.docx"
	r, _ := os.Open(ofname)
	buf, e := ioutil.ReadAll(r)
	client := NewRPCClient(useHTTP, network, address)
	o, e := client.PutObject(PutRequest{Data: buf, Key: key, Blocksize: uint32(blksz)})

	glog.V(4).Infof("Fingerprint for v1 %v\n %v\n", ofname, o.Fingerprint)

	nfname := "testdata/doc_v2.docx"
	delta := NewDiff(nfname, *o.Fingerprint)
	glog.V(4).Infof("Delta = %v ", delta)

	patched, e := client.PatchObject(PatchRequest{Delta: delta, Blocksize: uint32(blksz), Key: key})
	if e != nil {
		fmt.Printf("error ..%s", e)
		t.Fail()
	}
	v2fingerprint := NewFingerprint(nfname, uint32(blksz))
	glog.V(4).Infof("Generated fingerprint for version 2 %v\n", v2fingerprint)
	//update the source to key
	v2fingerprint.Source = key

	glog.V(4).Infof("Generated fingerprint after patch %v\n", patched.Fingerprint)
	if !v2fingerprint.DeepEqual(patched.Fingerprint) {
		t.Fail()
	}
}

func TestClientGetRequest(t *testing.T) {

	fmt.Println("Test RPC Get..")
	fmt.Printf("log level %v\n", *logLevel)
	flag.Lookup("v").Value.Set(fmt.Sprint(*logLevel))
	key := "TestClientGetRequestKey"
	fname := "testdata/26bytefile"
	r, _ := os.Open(fname)
	buf, e := ioutil.ReadAll(r)
	client := NewRPCClient(useHTTP, network, address)
	o, e := client.PutObject(PutRequest{Data: buf, Key: key, Blocksize: 8})
	if e != nil {
		fmt.Printf("error ..%s", e)
		t.Fail()
	}

	fp := o.Fingerprint

	//Now do a get
	o, e = client.GetObject(GetRequest{Key: key, Fingerprint: false})
	if e != nil {
		fmt.Printf("error ..%s", e)
		t.Fail()
	}
	if !reflect.DeepEqual(buf, o.Data) {
		fmt.Printf("Data doesn't match %s\n %v\n %v\n", o.Key, o.Data, buf)
		t.Fail()
	}

	o, e = client.GetObject(GetRequest{Key: key, Fingerprint: true})
	if !o.Fingerprint.DeepEqual(fp) {
		fmt.Printf("Fingerprint doesn't match %s, source 1 %s, source 2 %s\n", o.Key, o.Fingerprint.Source, fp.Source)
		t.Fail()
	}
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
	glog.V(2).Infof("Provider %v", *p)
	go ServeRPC(useHTTP, l, p)
}

func TestMain(m *testing.M) {
	runServer()
	flag.Parse()
	os.Exit(m.Run())
}
