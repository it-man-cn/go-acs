package main

import (
	"fmt"
	"go-acs/libs/encode"
)

const (
	//MappedAddress MAPPEDADDRESS
	MappedAddress = 0x0001
	//ResponseAddress RESPONSEADDRESS
	ResponseAddress = 0x0002
	//ChangeRequest CHANGEREQUEST
	ChangeRequest = 0x0003
	//SourceAddress SOURCEADDRESS
	SourceAddress = 0x0004
	//ChangedAddress CHANGEDADDRESS
	ChangedAddress = 0x0005
	//Username USERNAME
	Username = 0x0006
	//Password PASSWORD
	Password = 0x0007
	//MessageIntegrity MESSAGEINTEGRITY
	MessageIntegrity = 0x0008
	//ErrorCode ERRORCODE
	ErrorCode = 0x0009
	//UnknownAttribute UNKNOWNATTRIBUTE
	UnknownAttribute = 0x000a
	//ReflectedFrom REFLECTEDFROM
	ReflectedFrom = 0x000b
	//ConnectionRequestBinding CONNECTIONREQUESTBINDING
	ConnectionRequestBinding = 0xC001
	//BindingChange BINDINGCHANGE
	BindingChange = 0xC002
	//Dummy DUMMY
	Dummy = 0x0000
)

//Attribute stun msg attribute
type Attribute interface {
	//getType() uint16
	//getLength() uint16
	//getVal() string
	parse(buf []byte)
}

//ParseAttrs parse attributes
func ParseAttrs(buf []byte) (attributes []Attribute) {
	var (
		attrType uint16
		length   uint16
	)
	if len(buf) > 0 {
		attributes = make([]Attribute, 1)
		attrType = encode.Binary.Uint16(buf[0:2])
		length = encode.Binary.Uint16(buf[2:4])
		switch attrType {
		case Password:
			passwd := &password{}
			passwd.parse(buf[4 : 4+length])
			attributes = append(attributes, passwd)
		case ConnectionRequestBinding:
		}

	}
	fmt.Printf("in parse attrs len: %d\n", len(attributes))
	return
}
