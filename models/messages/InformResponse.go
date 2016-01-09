package messages

import (
	"encoding/xml"
	"fmt"
	"time"
)

type InformResponse struct {
	Id           string
	Name         string
	NoMore       int
	MaxEnvelopes int
}

type InformResponseBodyStruct struct {
	Body InformResponseStruct `xml:"cwmp:InformResponse"`
}

type InformResponseStruct struct {
	MaxEnvelopes int `xml:"MaxEnvelopes"`
}

func (msg *InformResponse) GetId() string {
	if len(msg.Id) < 1 {
		msg.Id = fmt.Sprintf("ID:intrnl.unset.id.%s%d.%d", msg.GetName(), time.Now().Unix(), time.Now().UnixNano())
	}
	return msg.Id
}

func (msg *InformResponse) GetName() string {
	return "InformResponse"
}

func (msg *InformResponse) CreateXml() []byte {
	env := Envelope{}
	env.XmlnsEnv = "http://schemas.xmlsoap.org/soap/envelope/"
	env.XmlnsEnc = "http://schemas.xmlsoap.org/soap/encoding/"
	env.XmlnsXsd = "http://www.w3.org/2001/XMLSchema"
	env.XmlnsXsi = "http://www.w3.org/2001/XMLSchema-instance"
	env.XmlnsCwmp = "urn:dslforum-org:cwmp-1-0"
	id := IdStruct{Attr: "1", Value: msg.GetId()}
	env.Header = HeaderStruct{ID: id, NoMore: msg.NoMore}
	respBody := InformResponseStruct{MaxEnvelopes: msg.MaxEnvelopes}
	env.Body = InformResponseBodyStruct{respBody}
	//output, err := xml.Marshal(env)
	output, err := xml.MarshalIndent(env, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	return output
}

func (msg *InformResponse) Parse(xmlstr string) {

}
