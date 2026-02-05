package main

import "slices"

type (
	opt string

	parameters struct {
		daID      string
		cID       string
		opt       opt
		idx       int
		readLimit int
		hostname  string
	}
)

const (
	optIss opt = "iss"
	optEat opt = "eat"
	optSub opt = "sub"

	hostname  = "http://localhost:8080"
	readLimit = 64 * 1024
)

const (
	codeOK = iota
	codeErr
)

var allowedOpts = []opt{"", optIss, optEat, optSub}

func (o opt) valid() bool {
	return slices.Contains(allowedOpts, o)
}
