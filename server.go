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
	"io/ioutil"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"
)

//TODO support range request
func (xrpc *Provider) GetObject(preq *GetRequest, presp *Response) error {

	ofile, _ := os.Open(filepath.Join(xrpc.FileStorePath, preq.Key))
	defer ofile.Close()

	var fp Fingerprint
	if preq.Fingerprint {
		fpfile, _ := os.Open(filepath.Join(xrpc.FileStorePath, preq.Key+".fingerprint"))
		defer fpfile.Close()
		dec := gob.NewDecoder(fpfile)
		if e := dec.Decode(&fp); e != nil {
			return e
		}

	}
	buf, _ := ioutil.ReadAll(ofile)
	presp.Object = Object{Key: preq.Key, Fingerprint: &fp, Data: buf}
	glog.V(2).Infof("Request successfully processed, Returning file %s", ofile.Name())
	return nil
}

//TODO will not work for large files
func (xrpc *Provider) PutObject(preq *PutRequest, presp *Response) error {

	ofile, _ := os.OpenFile(filepath.Join(xrpc.FileStorePath, preq.Key), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	fpfile, _ := os.OpenFile(filepath.Join(xrpc.FileStorePath, preq.Key+".fingerprint"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	defer ofile.Close()
	defer fpfile.Close()
	n, e := ofile.Write(preq.Data)
	if e != nil {
		return fmt.Errorf("Failed to create file %s, error %s", preq.Key, e)

	}
	//No version support
	f := NewFingerprintFromReader(bytes.NewReader(preq.Data), preq.Blocksize)
	f.Source = preq.Key
	enc := gob.NewEncoder(fpfile)
	enc.Encode(*f)

	presp.Object = Object{Key: preq.Key, Fingerprint: f}
	glog.V(2).Infof("Request successfully processed, bytes written %d", n)

	return nil
}

//TODO Version support
func (xrpc *Provider) PatchObject(preq *PatchRequest, presp *Response) error {
	sourceFile := filepath.Join(xrpc.FileStorePath, preq.Key)
	patchedFile := filepath.Join(xrpc.FileStorePath, preq.Key+".patched")

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
	r, _ := os.Open(sourceFile)
	f := NewFingerprintFromReader(r, preq.Blocksize)
	f.Source = preq.Key
	presp.Object = Object{Key: preq.Key, Fingerprint: f}
	glog.V(2).Infof("Request successfully processed, Patch file created %s", patchedFile)

	return nil
}

func ServeRPC(useHTTP bool, listener net.Listener, provider *Provider) {
	rpc.Register(provider)
	if useHTTP {
		glog.V(2).Infof("Starting RPC over HTTP..")
		rpc.HandleHTTP()
		http.Serve(listener, nil)
	} else {
		glog.V(2).Infof("Starting RPC Server..")
		rpc.Accept(listener)
	}

}
