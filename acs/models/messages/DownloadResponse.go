package messages

import (
	"encoding/xml"
	"fmt"
	"github.com/coraldane/godom"
	"strconv"
	"time"
)

//DownloadResponse download response
type DownloadResponse struct {
	ID           string
	Name         string
	Status       int
	StartTime    string
	CompleteTime string
}

type downloadResponseBodyStruct struct {
	DownResp downloadResponseStruct `xml:"cwmp:DownloadResponse"`
}

type downloadResponseStruct struct {
	Status       int
	StartTime    string
	CompleteTime string
}

//GetID tr069 msg id
func (msg *DownloadResponse) GetID() string {
	if len(msg.ID) < 1 {
		msg.ID = fmt.Sprintf("ID:intrnl.unset.id.%s%d.%d", msg.GetName(), time.Now().Unix(), time.Now().UnixNano())
	}
	return msg.ID
}

//GetName msg type name
func (msg *DownloadResponse) GetName() string {
	return "DownloadResponse"
}

//CreateXML encode into xml
func (msg *DownloadResponse) CreateXML() []byte {
	env := Envelope{}
	env.XmlnsEnv = "http://schemas.xmlsoap.org/soap/envelope/"
	env.XmlnsEnc = "http://schemas.xmlsoap.org/soap/encoding/"
	env.XmlnsXsd = "http://www.w3.org/2001/XMLSchema"
	env.XmlnsXsi = "http://www.w3.org/2001/XMLSchema-instance"
	env.XmlnsCwmp = "urn:dslforum-org:cwmp-1-0"
	id := IDStruct{Attr: "1", Value: msg.GetID()}
	env.Header = HeaderStruct{ID: id}
	body := downloadResponseStruct{
		StartTime:    msg.StartTime,
		CompleteTime: msg.CompleteTime,
		Status:       msg.Status,
	}
	env.Body = downloadResponseBodyStruct{body}
	//output, err := xml.Marshal(env)
	output, err := xml.MarshalIndent(env, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	return output
}

//Parse decode into struct
func (msg *DownloadResponse) Parse(xmlstr string) {
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
	startTime := root.GetElementsByTagName("StartTime")
	if startTime.Length() > 0 {
		msg.StartTime = startTime.Item(0).FirstChild().NodeValue()
	}
	completeTime := root.GetElementsByTagName("CompleteTime")
	if completeTime.Length() > 0 {
		msg.CompleteTime = completeTime.Item(0).FirstChild().NodeValue()
	}
}
