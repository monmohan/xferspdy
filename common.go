// Copyright 2015 Monmohan Singh. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xferspdy provides the basic interfaces around binary diff and patching process
package xferspdy

import (
	"github.com/golang/glog"
	"path/filepath"
)

type Object struct {
	Key         string
	VersionId   string
	Fingerprint *Fingerprint
	Data        []byte
}

type PutRequest struct {
	Data      []byte
	Key       string
	Blocksize uint32
}

type PatchRequest struct {
	Delta     []Block
	Key       string
	Blocksize uint32 //block size to use when generating the patched file fingerprint
}

type GetRequest struct {
	Key         string
	Fingerprint bool
}

type Response struct {
	Object Object
}

type Provider struct {
	//TODO The store should be configurable
	FileStorePath string
}

func NewProvider(filestorepath string) *Provider {
	absPath, e := filepath.Abs(filestorepath)
	if e != nil {
		glog.Fatalf("Error in setting file storage path %s", e)
	}
	return &Provider{FileStorePath: absPath}
}
