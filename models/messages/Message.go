package messages

import (
	"encoding/xml"
)

const (
	XSD_STRING      string = "xsd:string"
	XSD_UNSIGNEDINT string = "xsd:unsignedInt"
)

const (
	SOAP_ARRAY string = "SOAP-ENC:Array"
)

const (
	EVENT_BOOT_STRAP         string = "0 BOOTSTRAP"
	EVENT_BOOT               string = "1 BOOT"
	EVENT_PERIODIC           string = "2 PERIODIC"
	EVENT_SCHEDULED          string = "3 SCHEDULED"
	EVENT_VALUE_CHANGE       string = "4 VALUE CHANGE"
	EVENT_KICKED             string = "5 KICKED"
	EVENT_CONNECTION_REQUEST string = "6 CONNECTION REQUEST"
	EVENT_TRANSFER_COMPLETE  string = "7 TRANSFER COMPLETE"
	EVENT_CLIENT_CHANGE      string = "8 CLIENT CHANGE"
)

type Message interface {
	Parse(xmlstr string)
	CreateXml() []byte
	GetName() string
	GetId() string
}

type Envelope struct {
	XMLName   xml.Name    `xml:"SOAP-ENV:Envelope"`
	XmlnsEnv  string      `xml:"xmlns:SOAP-ENV,attr"`
	XmlnsEnc  string      `xml:"xmlns:SOAP-ENC,attr"`
	XmlnsXsd  string      `xml:"xmlns:xsd,attr"`
	XmlnsXsi  string      `xml:"xmlns:xsi,attr"`
	XmlnsCwmp string      `xml:"xmlns:cwmp,attr"`
	Header    interface{} `xml:"SOAP-ENV:Header"`
	Body      interface{} `xml:"SOAP-ENV:Body"`
}

type HeaderStruct struct {
	ID     IdStruct    `xml:"cwmp:ID"`
	NoMore interface{} `xml:"cwmp:NoMoreRequests,ommitempty"`
}

type IdStruct struct {
	Attr  string `xml:"SOAP-ENV:mustUnderstand,attr,ommitempty"`
	Value string `xml:",chardata"`
}

type NodeStruct struct {
	Type  interface{} `xml:"xsi:type,attr"`
	Value string      `xml:",chardata"`
}

type EventStruct struct {
	Type   string            `xml:"SOAP-ENC:arrayType,attr"`
	Events []EventNodeStruct `xml:"EventStruct"`
}

type EventNodeStruct struct {
	EventCode  NodeStruct `xml:"EventCode"`
	CommandKey string     `xml:"CommandKey"`
}

type ParameterListStruct struct {
	Type   string                 `xml:"SOAP-ENC:arrayType,attr"`
	Params []ParameterValueStruct `xml:"ParameterValueStruct"`
}

type ParameterValueStruct struct {
	Name  NodeStruct `xml:"Name"`
	Value NodeStruct `xml:"Value"`
}

type FaultStruct struct {
	FaultCode   int
	FaultString string
}

type ValueStruct struct {
	Type  string
	Value string
}
