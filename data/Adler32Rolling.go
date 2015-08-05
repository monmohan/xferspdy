package data

import (
	"fmt"
)

const (
	// mod is the largest prime that is less than 65536.
	mod = 65521
	//number of bytes that can be added
	nmax = 5552
)

// The low 16 bits are s1, the high 16 bits are s2.
type checksum uint32

type State struct {
	window []byte
	s1     uint32
	s2     uint32
}

func Checksum(p []byte) (uint32, *State) {
	s1, s2 := uint32(1&0xffff), uint32(1>>16)
	fmt.Printf("Init: s1 %d s2 %d\n", s1, s2)
	orig := p
	for len(p) > 0 {
		var q []byte
		if len(p) > nmax {
			p, q = p[:nmax], p[nmax:]
		}
		for _, x := range p {
			s1 += uint32(x)
			s2 += s1
		}
		s1 %= mod
		s2 %= mod
		p = q
	}
	fmt.Printf("s1 %d s2 %d\n", s1, s2)
	return uint32(s2<<16 | s1), &State{orig, s1, s2}
}

func (s *State) UpdateWindow(nb byte) uint32 {
	fmt.Printf("Update window init : s1 %d s2 %d byte appended %d \n", s.s1, s.s2, nb)
	s.window = append(s.window, nb)
	x := s.window[0]
	s.window = s.window[1:]
	s.s1 = s.s1 + uint32(nb) - uint32(x)
	s.s1 %= mod
	b := (uint32(len(s.window)) * uint32(x)) + 1
	a := s.s2 + s.s1
	for b > a {
		fmt.Printf("b %d greater than a %d\n", b, a)
		a += mod
	}
	s.s2 = a - b
	s.s2 %= mod
	fmt.Printf("Update window: s1 %d s2 %d\n", s.s1, s.s2)
	return uint32(s.s2<<16 | s.s1)
}
