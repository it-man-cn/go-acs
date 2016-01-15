package main

import (
	"fmt"
	"go-acs/libs/encode"
)

const (
	MAPPEDADDRESS            = 0x0001
	RESPONSEADDRESS          = 0x0002
	CHANGEREQUEST            = 0x0003
	SOURCEADDRESS            = 0x0004
	CHANGEDADDRESS           = 0x0005
	USERNAME                 = 0x0006
	PASSWORD                 = 0x0007
	MESSAGEINTEGRITY         = 0x0008
	ERRORCODE                = 0x0009
	UNKNOWNATTRIBUTE         = 0x000a
	REFLECTEDFROM            = 0x000b
	CONNECTIONREQUESTBINDING = 0xC001
	BINDINGCHANGE            = 0xC002
	DUMMY                    = 0x0000
)

type attribute interface {
	//getType() uint16
	//getLength() uint16
	//getVal() string
	parse(buf []byte)
}

func ParseAttrs(buf []byte) (attributes []attribute) {
	var (
		attrType uint16
		length   uint16
	)
	if len(buf) > 0 {
		attributes = make([]attribute, 1)
		attrType = encode.Binary.Uint16(buf[0:2])
		length = encode.Binary.Uint16(buf[2:4])
		switch attrType {
		case PASSWORD:
			passwd := &password{}
			passwd.parse(buf[4 : 4+length])
			attributes = append(attributes, passwd)
		case CONNECTIONREQUESTBINDING:
		}

	}
	fmt.Printf("in parse attrs len: %d\n", len(attributes))
	return
}
