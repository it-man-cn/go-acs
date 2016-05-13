package main

import (
	"errors"
)

var (
	//ErrRingEmpty ring empty
	ErrRingEmpty = errors.New("ring buffer empty")
	//ErrRingFull ring full
	ErrRingFull = errors.New("ring buffer full")
)
