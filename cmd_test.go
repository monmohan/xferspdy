// Copyright 2015 Monmohan Singh. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package xferspdy

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"
)

var (
	filev1      = "testdata/SamplePPT_v1.pptx"
	fprint      = "testdata/SamplePPT_v1.pptx.fingerprint"
	filev2      = "testdata/SamplePPT_v2.pptx"
	fpatch      = "testdata/SamplePPT_v2.pptx.patch"
	patchedFile = "testdata/Patched_SamplePPT_v1.pptx"
)

func TestCmdUtilities(t *testing.T) {

	cleanup()
	defer cleanup()
	//fpgen
	var cmd = exec.Command(fmt.Sprintf("%s/bin/fpgen", os.ExpandEnv("$GOPATH")), "-file", filev1)
	handleCommand(cmd, t)
	//diff
	cmd = exec.Command(fmt.Sprintf("%s/bin/diff", os.ExpandEnv("$GOPATH")), "-file", filev2, "-fingerprint", fprint)
	handleCommand(cmd, t)
	//patch
	cmd = exec.Command(fmt.Sprintf("%s/bin/patch", os.ExpandEnv("$GOPATH")), "-base", filev1, "-patch", fpatch)
	handleCommand(cmd, t)

	sign1 := NewFingerprint(filev2, uint32(2048))
	sign2 := NewFingerprint(patchedFile, uint32(2048))
	if sign1.DeepEqual(sign2) {
		fmt.Printf("Signature matched %s %s \n", sign1.Source, sign2.Source)
	} else {
		t.Fail()
	}

}

func handleCommand(cmd *exec.Cmd, t *testing.T) {
	var out bytes.Buffer
	var eout bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &eout
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error %q\n", eout.String())
		log.Fatal(err)
		t.FailNow()

	}

	fmt.Printf("Output %q\n", out.String())
}

func cleanup() {
	os.Remove(fprint)
	os.Remove(fpatch)
	os.Remove(patchedFile)
}
