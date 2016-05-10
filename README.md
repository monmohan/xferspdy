# xferspdy

Fast binary diffs -
* Command line utilities to diff and patch binary files
* Library for fingerprint generation, rolling hash and block matching

Reference :
[Rsync Algorithm] (https://rsync.samba.org/tech_report/node2.html)

### Using the fpgen, diff and patch CLI utilities:
* Need to have go installed, [golang downloads] (https://golang.org/dl/)
* Clone the project 

    `git clone https://github.com/monmohan/xferspdy.git`
OR
* Do go get

     `go get github.com/monmohan/xferspdy`
* Install dependencies

     `go get github.com/golang/glog`

* Install the binaries (example from from xferspdy directory)

	`go install ./cmd/diff`
	    
  `go install ./cmd/patch`
    	
  `go install ./cmd/fpgen`

### Try it out
* Lets say you have a binary file  (e.g. power point presentation MyPrezVersion1.pptx).
* First generate a fingerprint of version 1

  `$GOPATH/bin/fpgen -file <path>/MyPrezVersion1.pptx`

* This will generate the fingerprint file <path>/MyPrezVersion1.pptx.fingerprint.
* Lets say that the file was changed now (for example add a slide or image) and saved as MyPrezVersion2.pptx.pptx
* Now Generate a diff (doesn't require original file)

   `$GOPATH/bin/diff -fingerprint <path>/MyPrezVersion1.pptx.fingerprint -file <path>/MyPrezVersion2.pptx`

 It will create a patch file <path>/MyPrezVersion2.pptx.patch

* Now patch the Version 1 file to get the Version 2
 
   `$GOPATH/bin/patch -patch <path>/MyPrezVersion2.pptx.patch -base <path>/MyPrezVersion1.pptx`

* This will generate <path>/Patched_MyPrezVersion1.pptx. This file would exactly be same as MyPrezVersion2.pptx.

NOTE: diff and patch are also common utilities present on most distributions so its better to give explicit path to these binaries. It would be $GOPATH/bin/diff and $GOPATH/bin/patch

