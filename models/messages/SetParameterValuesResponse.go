package messages

import (
	"encoding/xml"
	"fmt"
	"github.com/coraldane/godom"
	"strconv"
	"time"
)

type SetParameterValuesResponse struct {
	Id           string
	Name         string
	Status       int
	ParameterKey string
}

type SetParameterValuesResponseBodyStruct struct {
	Body SetParameterValuesResponseStruct `xml:"cwmp:SetParameterValuesResponse"`
}

type SetParameterValuesResponseStruct struct {
	Status       int
	ParameterKey string
}

func (msg *SetParameterValuesResponse) GetId() string {
	if len(msg.Id) < 1 {
		msg.Id = fmt.Sprintf("ID:intrnl.unset.id.%s%d.%d", msg.GetName(), time.Now().Unix(), time.Now().UnixNano())
	}
	return msg.Id
}

func (msg *SetParameterValuesResponse) GetName() string {
	return "SetParameterValuesResponse"
}

func (msg *SetParameterValuesResponse) CreateXml() []byte {
	env := Envelope{}
	env.XmlnsEnv = "http://schemas.xmlsoap.org/soap/envelope/"
	env.XmlnsEnc = "http://schemas.xmlsoap.org/soap/encoding/"
	env.XmlnsXsd = "http://www.w3.org/2001/XMLSchema"
	env.XmlnsXsi = "http://www.w3.org/2001/XMLSchema-instance"
	env.XmlnsCwmp = "urn:dslforum-org:cwmp-1-0"
	id := IdStruct{Attr: "1", Value: msg.GetId()}
	env.Header = HeaderStruct{ID: id}
	body := SetParameterValuesResponseStruct{
		Status:       msg.Status,
		ParameterKey: msg.ParameterKey,
	}
	env.Body = SetParameterValuesResponseBodyStruct{body}
	//output, err := xml.Marshal(env)
	output, err := xml.MarshalIndent(env, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	return output
}

func (msg *SetParameterValuesResponse) Parse(xmlstr string) {
	document, _ := dom.ParseString(xmlstr)
	root := document.DocumentElement()
	hdr := root.GetElementsByTagName("Header")
	if hdr.Length() > 0 {
		pNode := hdr.Item(0)
		msg.Id = GetChildElementValue(pNode, "ID")
	}
	status := root.GetElementsByTagName("Status")
	if status.Length() > 0 {
		msg.Status, _ = strconv.Atoi(status.Item(0).FirstChild().NodeValue())
	}
	key := root.GetElementsByTagName("ParameterKey")
	fmt.Println("key", key)
	if key.Length() > 0 {
		fmt.Println("key item0", key.Item(0).HasChildNodes())
		if key.Item(0).HasChildNodes() {
			msg.ParameterKey = key.Item(0).FirstChild().NodeValue()
		}
	}
}
