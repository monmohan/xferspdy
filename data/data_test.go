package data

import (
	"fmt"
	"testing"
)

func TestSignatureCreate(t *testing.T) {
	sign := NewSignature("/msingh/projects/genknow/gitcheatsheet", 16)
	fmt.Printf(" %v\n", *sign)
}
