// Copyright 2015 Monmohan Singh. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xferspdy provides the basic interfaces around binary diff and patching process
package xferspdy

import (
	"github.com/golang/glog"
	"net/rpc"
)

type Client interface {
	PutObject(req PutRequest) (obj Object, err error)
}
type RPCClient struct {
	Client *rpc.Client
}

func NewRPCClient(useHTTP bool, network string, address string) *RPCClient {
	var client *rpc.Client
	var err error
	if useHTTP {
		client, err = rpc.DialHTTP(network, address)
	} else {
		client, err = rpc.Dial(network, address)
	}
	if err != nil {
		glog.Fatalf("dialing failed: %s", err)
	}
	rpcl := &RPCClient{Client: client}
	return rpcl

}

func (rpcl *RPCClient) PutObject(req PutRequest) (obj Object, err error) {
	var reply Response
	err = rpcl.Client.Call("Provider.PutObject", &req, &reply)
	if err != nil {
		glog.Fatal("Call errror:", err)
	}
	return reply.Object, err
}

func (rpcl *RPCClient) GetObject(req GetRequest) (obj Object, err error) {
	var reply Response
	err = rpcl.Client.Call("Provider.GetObject", &req, &reply)
	if err != nil {
		glog.Fatal("Call errror:", err)
	}
	return reply.Object, err
}

func (rpcl *RPCClient) PatchObject(req PatchRequest) (obj Object, err error) {
	var reply Response
	err = rpcl.Client.Call("Provider.PatchObject", &req, &reply)
	if err != nil {
		glog.Fatal("Call errror:", err)
	}
	return reply.Object, err
}
