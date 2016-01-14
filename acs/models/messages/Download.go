package messages

import (
	"encoding/xml"
	"fmt"
	"time"
)

const (
	FTFireware   string = "1 Firmware Upgrade Image"
	FTWebContent string = "2 Web Content"
	FTConfig     string = "3 Vendor Configuration File"
)

type Download struct {
	Id             string
	Name           string
	NoMore         int
	CommandKey     string
	FileType       string
	URL            string
	Username       string
	Password       string
	FileSize       int
	TargetFileName string
	DelaySeconds   int
	SuccessURL     string
	FailureURL     string
}

type DownloadBodyStruct struct {
	Body DownloadStruct `xml:"cwmp:Download"`
}

type DownloadStruct struct {
	CommandKey     string
	FileType       string
	URL            string
	Username       string
	Password       string
	FileSize       int
	TargetFileName string
	DelaySeconds   int
	SuccessURL     string
	FailureURL     string
}

func (msg *Download) GetId() string {
	if len(msg.Id) < 1 {
		msg.Id = fmt.Sprintf("ID:intrnl.unset.id.%s%d.%d", msg.GetName(), time.Now().Unix(), time.Now().UnixNano())
	}
	return msg.Id
}

func (msg *Download) GetName() string {
	return "Download"
}

func (msg *Download) CreateXml() []byte {
	env := Envelope{}
	env.XmlnsEnv = "http://schemas.xmlsoap.org/soap/envelope/"
	env.XmlnsEnc = "http://schemas.xmlsoap.org/soap/encoding/"
	env.XmlnsXsd = "http://www.w3.org/2001/XMLSchema"
	env.XmlnsXsi = "http://www.w3.org/2001/XMLSchema-instance"
	env.XmlnsCwmp = "urn:dslforum-org:cwmp-1-0"
	id := IdStruct{Attr: "1", Value: msg.GetId()}
	env.Header = HeaderStruct{ID: id, NoMore: msg.NoMore}
	body := DownloadStruct{
		CommandKey:     msg.CommandKey,
		FileType:       msg.FileType,
		URL:            msg.URL,
		Username:       msg.Username,
		Password:       msg.Password,
		FileSize:       msg.FileSize,
		TargetFileName: msg.TargetFileName,
		DelaySeconds:   msg.DelaySeconds,
		SuccessURL:     msg.SuccessURL,
		FailureURL:     msg.FailureURL}
	env.Body = DownloadBodyStruct{body}
	//output, err := xml.Marshal(env)
	output, err := xml.MarshalIndent(env, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	return output
}

func (msg *Download) Parse(xmlstr string) {

}
