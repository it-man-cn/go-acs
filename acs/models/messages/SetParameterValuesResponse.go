package messages

import (
	"encoding/xml"
	"fmt"
	"github.com/coraldane/godom"
	"strconv"
	"time"
)

//SetParameterValuesResponse set param reponse
type SetParameterValuesResponse struct {
	ID           string
	Name         string
	Status       int
	ParameterKey string
}

type setParameterValuesResponseBodyStruct struct {
	Body setParameterValuesResponseStruct `xml:"cwmp:SetParameterValuesResponse"`
}

type setParameterValuesResponseStruct struct {
	Status       int
	ParameterKey string
}

//GetID get msg id
func (msg *SetParameterValuesResponse) GetID() string {
	if len(msg.ID) < 1 {
		msg.ID = fmt.Sprintf("ID:intrnl.unset.id.%s%d.%d", msg.GetName(), time.Now().Unix(), time.Now().UnixNano())
	}
	return msg.ID
}

//GetName get msg type
func (msg *SetParameterValuesResponse) GetName() string {
	return "SetParameterValuesResponse"
}

//CreateXML encode into xml
func (msg *SetParameterValuesResponse) CreateXML() []byte {
	env := Envelope{}
	env.XmlnsEnv = "http://schemas.xmlsoap.org/soap/envelope/"
	env.XmlnsEnc = "http://schemas.xmlsoap.org/soap/encoding/"
	env.XmlnsXsd = "http://www.w3.org/2001/XMLSchema"
	env.XmlnsXsi = "http://www.w3.org/2001/XMLSchema-instance"
	env.XmlnsCwmp = "urn:dslforum-org:cwmp-1-0"
	id := IDStruct{Attr: "1", Value: msg.GetID()}
	env.Header = HeaderStruct{ID: id}
	body := setParameterValuesResponseStruct{
		Status:       msg.Status,
		ParameterKey: msg.ParameterKey,
	}
	env.Body = setParameterValuesResponseBodyStruct{body}
	//output, err := xml.Marshal(env)
	output, err := xml.MarshalIndent(env, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	return output
}

//Parse decode from xml
func (msg *SetParameterValuesResponse) Parse(xmlstr string) {
	document, _ := dom.ParseString(xmlstr)
	root := document.DocumentElement()
	hdr := root.GetElementsByTagName("Header")
	if hdr.Length() > 0 {
		pNode := hdr.Item(0)
		msg.ID = GetChildElementValue(pNode, "ID")
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
