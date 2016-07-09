// Copyright 2015 Monmohan Singh. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xferspdy provides the basic interfaces around binary diff and patching process
package xferspdy

import (
	"bytes"
	"fmt"
	//"github.com/golang/glog"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"
)

type Object struct {
	Key         string
	VersionId   string
	Fingerprint *Fingerprint
}

type PutRequest struct {
	Data      []byte
	Key       string
	Blocksize uint32
}

type PutResponse struct {
	Object Object
}

type Provider struct {
	//TODO The store should be configurable
	FileStorePath string
}

func NewProvider(filestorepath string) *Provider {
	absPath, e := filepath.Abs(filestorepath)
	if e != nil {
		log.Fatalf("Error in setting file storage path %s", e)
	}
	return &Provider{FileStorePath: absPath}
}

//TODO will not work for large files
func (xrpc *Provider) PutObject(preq *PutRequest, presp *PutResponse) error {

	ofile, _ := os.OpenFile(filepath.Join(xrpc.FileStorePath, preq.Key), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	n, e := ofile.Write(preq.Data)
	if e != nil {
		return fmt.Errorf("Failed to create file %s, error %s", preq.Key)

	}
	//No version support
	f := NewFingerprintFromReader(bytes.NewReader(preq.Data), preq.Blocksize)
	f.Source = preq.Key
	presp.Object = Object{Key: preq.Key, Fingerprint: f}
	log.Printf("Request successfully processed, bytes written %d", n)

	return nil
}

func ServeRPC(useHTTP bool, listener net.Listener, provider *Provider) {
	rpc.Register(provider)
	if useHTTP {
		rpc.HandleHTTP()
		http.Serve(listener, nil)
	} else {
		fmt.Println("Starting RPC..")
		rpc.Accept(listener)
	}

}
