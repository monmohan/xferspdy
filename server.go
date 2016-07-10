// Copyright 2015 Monmohan Singh. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xferspdy provides the basic interfaces around binary diff and patching process
package xferspdy

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/golang/glog"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"
)

//TODO will not work for large files
func (xrpc *Provider) PutObject(preq *PutRequest, presp *Response) error {

	ofile, _ := os.OpenFile(filepath.Join(xrpc.FileStorePath, preq.Key), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	fpfile, _ := os.OpenFile(filepath.Join(xrpc.FileStorePath, preq.Key+".fingerprint"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	defer ofile.Close()
	defer fpfile.Close()
	n, e := ofile.Write(preq.Data)
	if e != nil {
		return fmt.Errorf("Failed to create file %s, error %s", preq.Key)

	}
	//No version support
	f := NewFingerprintFromReader(bytes.NewReader(preq.Data), preq.Blocksize)

	enc := gob.NewEncoder(fpfile)
	enc.Encode(*f)

	f.Source = preq.Key
	presp.Object = Object{Key: preq.Key, Fingerprint: f}
	glog.V(2).Infof("Request successfully processed, bytes written %d", n)

	return nil
}

//TODO Version support
func (xrpc *Provider) PatchObject(preq *PatchRequest, presp *Response) error {
	sourceFile := filepath.Join(xrpc.FileStorePath, preq.Key)
	patchedFile := filepath.Join(xrpc.FileStorePath, preq.Key, ".patched")

	pfile, _ := os.OpenFile(patchedFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	defer pfile.Close()

	e := PatchFile(preq.Delta, sourceFile, pfile)

	if e != nil {
		return fmt.Errorf("Failed to an patch file %s, error %s", preq.Key, e)

	}
	e = os.Rename(patchedFile, sourceFile)

	if e != nil {
		return fmt.Errorf("Failed to an patch file %s, error %s", preq.Key, e)

	}
	f := NewFingerprint(sourceFile, preq.Blocksize)
	presp.Object = Object{Key: preq.Key, Fingerprint: f}
	glog.V(2).Infof("Request successfully processed, Patch file created %s", patchedFile)

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
