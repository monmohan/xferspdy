#!/bin/sh
$GOPATH/bin/fpgen -file SamplePPT_v1.pptx
$GOPATH/bin/diff -fingerprint SamplePPT_v1.pptx.fingerprint -file SamplePPT_v2.pptx
$GOPATH/bin/patch -patch SamplePPT_v2.pptx.patch -base SamplePPT_v1.pptx

