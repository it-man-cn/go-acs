package messages

import (
	"encoding/xml"
	"fmt"
	"github.com/coraldane/godom"
	"strconv"
	"time"
)

type DownloadResponse struct {
	Id           string
	Name         string
	Status       int
	StartTime    string
	CompleteTime string
}

type DownloadResponseBodyStruct struct {
	DownResp DownloadResponseStruct `xml:"cwmp:DownloadResponse"`
}

type DownloadResponseStruct struct {
	Status       int
	StartTime    string
	CompleteTime string
}

func (msg *DownloadResponse) GetId() string {
	if len(msg.Id) < 1 {
		msg.Id = fmt.Sprintf("ID:intrnl.unset.id.%s%d.%d", msg.GetName(), time.Now().Unix(), time.Now().UnixNano())
	}
	return msg.Id
}

func (msg *DownloadResponse) GetName() string {
	return "DownloadResponse"
}

func (msg *DownloadResponse) CreateXml() []byte {
	env := Envelope{}
	env.XmlnsEnv = "http://schemas.xmlsoap.org/soap/envelope/"
	env.XmlnsEnc = "http://schemas.xmlsoap.org/soap/encoding/"
	env.XmlnsXsd = "http://www.w3.org/2001/XMLSchema"
	env.XmlnsXsi = "http://www.w3.org/2001/XMLSchema-instance"
	env.XmlnsCwmp = "urn:dslforum-org:cwmp-1-0"
	id := IdStruct{Attr: "1", Value: msg.GetId()}
	env.Header = HeaderStruct{ID: id}
	body := DownloadResponseStruct{
		StartTime:    msg.StartTime,
		CompleteTime: msg.CompleteTime,
		Status:       msg.Status,
	}
	env.Body = DownloadResponseBodyStruct{body}
	//output, err := xml.Marshal(env)
	output, err := xml.MarshalIndent(env, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	return output
}

func (msg *DownloadResponse) Parse(xmlstr string) {
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
	startTime := root.GetElementsByTagName("StartTime")
	if startTime.Length() > 0 {
		msg.StartTime = startTime.Item(0).FirstChild().NodeValue()
	}
	completeTime := root.GetElementsByTagName("CompleteTime")
	if completeTime.Length() > 0 {
		msg.CompleteTime = completeTime.Item(0).FirstChild().NodeValue()
	}
}
