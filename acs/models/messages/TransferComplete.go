package messages

import (
	"encoding/xml"
	"fmt"
	"github.com/coraldane/godom"
	"strconv"
	"time"
)

//TransferComplete download complete
type TransferComplete struct {
	ID           string
	Name         string
	CommandKey   string
	StartTime    string
	CompleteTime string
	FaultCode    int
	FaultString  string
}

type transferCompleteBodyStruct struct {
	Body transferCompleteStruct `xml:"cwmp:TransferComplete"`
}

type transferCompleteStruct struct {
	CommandKey   string
	StartTime    string
	CompleteTime string
	Fault        interface{} `xml:"FaultStruct,ommitempty"`
}

//GetID get msg id
func (msg *TransferComplete) GetID() string {
	if len(msg.ID) < 1 {
		msg.ID = fmt.Sprintf("ID:intrnl.unset.id.%s%d.%d", msg.GetName(), time.Now().Unix(), time.Now().UnixNano())
	}
	return msg.ID
}

//GetName get msg type
func (msg *TransferComplete) GetName() string {
	return "TransferComplete"
}

//CreateXML encode into mxl
func (msg *TransferComplete) CreateXML() []byte {
	env := Envelope{}
	env.XmlnsEnv = "http://schemas.xmlsoap.org/soap/envelope/"
	env.XmlnsEnc = "http://schemas.xmlsoap.org/soap/encoding/"
	env.XmlnsXsd = "http://www.w3.org/2001/XMLSchema"
	env.XmlnsXsi = "http://www.w3.org/2001/XMLSchema-instance"
	env.XmlnsCwmp = "urn:dslforum-org:cwmp-1-0"
	id := IDStruct{Attr: "1", Value: msg.GetID()}
	env.Header = HeaderStruct{ID: id}
	var body transferCompleteStruct
	if len(msg.FaultString) > 0 {
		fault := FaultStruct{FaultCode: msg.FaultCode, FaultString: msg.FaultString}
		body = transferCompleteStruct{
			CommandKey:   msg.CommandKey,
			StartTime:    msg.StartTime,
			CompleteTime: msg.CompleteTime,
			Fault:        fault,
		}
	} else {
		body = transferCompleteStruct{
			CommandKey:   msg.CommandKey,
			StartTime:    msg.StartTime,
			CompleteTime: msg.CompleteTime,
		}
	}

	env.Body = transferCompleteBodyStruct{body}
	//output, err := xml.Marshal(env)
	output, err := xml.MarshalIndent(env, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	return output
}

//Parse decode from xml
func (msg *TransferComplete) Parse(xmlstr string) {
	document, _ := dom.ParseString(xmlstr)
	root := document.DocumentElement()
	hdr := root.GetElementsByTagName("Header")
	if hdr.Length() > 0 {
		pNode := hdr.Item(0)
		msg.ID = GetChildElementValue(pNode, "ID")
	}
	cmdKey := root.GetElementsByTagName("CommandKey")
	if cmdKey.Length() > 0 {
		msg.CommandKey = cmdKey.Item(0).FirstChild().NodeValue()
	}
	startTime := root.GetElementsByTagName("StartTime")
	if startTime.Length() > 0 {
		msg.StartTime = startTime.Item(0).FirstChild().NodeValue()
	}
	completeTime := root.GetElementsByTagName("CompleteTime")
	if completeTime.Length() > 0 {
		msg.CompleteTime = completeTime.Item(0).FirstChild().NodeValue()
	}
	faultCode := root.GetElementsByTagName("FaultCode")
	if faultCode.Length() > 0 {
		msg.FaultCode, _ = strconv.Atoi(faultCode.Item(0).FirstChild().NodeValue())
	}
	faultString := root.GetElementsByTagName("FaultString")
	if faultString.Length() > 0 {
		msg.FaultString = faultString.Item(0).FirstChild().NodeValue()
	}
}
