package main

import (
	"encoding/hex"
	"fmt"
	"go-acs/libs/encode"
)

type bindingRequest struct {
	Type       int16
	Length     int16
	Id         []byte
	Attributes []attribute
}

func (msg *bindingRequest) Decode(buf []byte) {
	//fmt.Printf("len%d\n", len(buf))
	fmt.Println("request decode :" + string(buf))
	msg.Type = encode.Binary.Int16(buf[0:2])
	fmt.Println(msg.Type)

	msg.Length = encode.Binary.Int16(buf[2:4])
	fmt.Printf("msg length %d\n", msg.Length)
	msg.Id = buf[4:20]
	//msg.Id = string(buf[4:20])
	fmt.Println("msg id:" + hex.EncodeToString(msg.Id))
	msg.Attributes = ParseAttrs(buf[20:msg.Length])
	fmt.Printf("out parse attrs len: %d\n", len(msg.Attributes))
	password := msg.Attributes[0].(*password)
	fmt.Println(password.Password)
}
