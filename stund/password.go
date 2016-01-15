package main

import (
//"go-acs/libs/encode"
)

type password struct {
	Password string
}

func (p *password) parse(buf []byte) {
	p.Password = string(buf)
	return
}
