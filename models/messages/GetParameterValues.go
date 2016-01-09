package messages

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"time"
)

type GetParameterValues struct {
	Id             string
	Name           string
	NoMore         int
	ParameterNames []string
}

type GetParameterValuesBodyStruct struct {
	Body GetParameterValuesStruct `xml:"cwmp:GetParameterValues"`
}

type GetParameterValuesStruct struct {
	Params ParameterNamesStruct `xml:"ParameterNames"`
}

type ParameterNamesStruct struct {
	Type       string   `xml:"SOAP-ENC:arrayType,attr"`
	ParamNames []string `xml:"string"`
}

func (msg *GetParameterValues) GetName() string {
	return "GetParameterValues"
}

func (msg *GetParameterValues) GetId() string {
	if len(msg.Id) < 1 {
		msg.Id = fmt.Sprintf("ID:intrnl.unset.id.%s%d.%d", msg.GetName(), time.Now().Unix(), time.Now().UnixNano())
	}
	return msg.Id
}

func (msg *GetParameterValues) CreateXml() []byte {
	env := Envelope{}
	id := IdStruct{"1", msg.GetId()}
	env.XmlnsEnv = "http://schemas.xmlsoap.org/soap/envelope/"
	env.XmlnsEnc = "http://schemas.xmlsoap.org/soap/encoding/"
	env.XmlnsXsd = "http://www.w3.org/2001/XMLSchema"
	env.XmlnsXsi = "http://www.w3.org/2001/XMLSchema-instance"
	env.XmlnsCwmp = "urn:dslforum-org:cwmp-1-0"
	env.Header = HeaderStruct{ID: id, NoMore: msg.NoMore}
	paramLen := strconv.Itoa(len(msg.ParameterNames))
	paramNames := ParameterNamesStruct{
		Type: XSD_STRING + "[" + paramLen + "]",
	}
	for _, v := range msg.ParameterNames {
		paramNames.ParamNames = append(paramNames.ParamNames, v)
	}
	body := GetParameterValuesStruct{paramNames}
	env.Body = GetParameterValuesBodyStruct{body}
	output, err := xml.MarshalIndent(env, "  ", "    ")
	//output, err := xml.Marshal(env)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	return output
}

func (msg *GetParameterValues) Parse(xmlstr string) {

}
