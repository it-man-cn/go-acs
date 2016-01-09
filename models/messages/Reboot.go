package messages

import (
	"encoding/xml"
	"fmt"
	"time"
)

type Reboot struct {
	Id         string
	Name       string
	NoMore     int
	CommandKey string
}

type RebootBodyStruct struct {
	Body RebootStruct `xml:"cwmp:Reboot"`
}

type RebootStruct struct {
	CommandKey string
}

func (msg *Reboot) GetId() string {
	if len(msg.Id) < 1 {
		msg.Id = fmt.Sprintf("ID:intrnl.unset.id.%s%d.%d", msg.GetName(), time.Now().Unix(), time.Now().UnixNano())
	}
	return msg.Id
}

func (msg *Reboot) GetName() string {
	return "Reboot"
}

func (msg *Reboot) CreateXml() []byte {
	env := Envelope{}
	env.XmlnsEnv = "http://schemas.xmlsoap.org/soap/envelope/"
	env.XmlnsEnc = "http://schemas.xmlsoap.org/soap/encoding/"
	env.XmlnsXsd = "http://www.w3.org/2001/XMLSchema"
	env.XmlnsXsi = "http://www.w3.org/2001/XMLSchema-instance"
	env.XmlnsCwmp = "urn:dslforum-org:cwmp-1-0"
	id := IdStruct{Attr: "1", Value: msg.GetId()}
	env.Header = HeaderStruct{ID: id, NoMore: msg.NoMore}
	body := RebootStruct{CommandKey: msg.CommandKey}
	env.Body = RebootBodyStruct{body}
	//output, err := xml.Marshal(env)
	output, err := xml.MarshalIndent(env, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	return output
}

func (msg *Reboot) Parse(xmlstr string) {

}
