package messages

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"time"
)

type SetParameterValues struct {
	Id           string
	Name         string
	NoMore       int
	Params       map[string]ValueStruct
	ParameterKey string
}

type SetParameterValuesBodyStruct struct {
	Body SetParameterValuesStruct `xml:"cwmp:SetParameterValues"`
}

type SetParameterValuesStruct struct {
	ParamList    ParameterListStruct `xml:"ParameterList"`
	ParameterKey string
}

func (msg *SetParameterValues) GetName() string {
	return "SetParameterValues"
}

func (msg *SetParameterValues) GetId() string {
	if len(msg.Id) < 1 {
		msg.Id = fmt.Sprintf("ID:intrnl.unset.id.%s%d.%d", msg.GetName(), time.Now().Unix(), time.Now().UnixNano())
	}
	return msg.Id
}

func (msg *SetParameterValues) CreateXml() []byte {
	env := Envelope{}
	id := IdStruct{"1", msg.GetId()}
	env.XmlnsEnv = "http://schemas.xmlsoap.org/soap/envelope/"
	env.XmlnsEnc = "http://schemas.xmlsoap.org/soap/encoding/"
	env.XmlnsXsd = "http://www.w3.org/2001/XMLSchema"
	env.XmlnsXsi = "http://www.w3.org/2001/XMLSchema-instance"
	env.XmlnsCwmp = "urn:dslforum-org:cwmp-1-0"
	env.Header = HeaderStruct{ID: id, NoMore: msg.NoMore}

	paramLen := strconv.Itoa(len(msg.Params))
	paramList := ParameterListStruct{Type: "cwmp:ParameterValueStruct[" + paramLen + "]"}
	for k, v := range msg.Params {
		param := ParameterValueStruct{
			Name:  NodeStruct{Value: k},
			Value: NodeStruct{Type: v.Type, Value: v.Value}}
		paramList.Params = append(paramList.Params, param)
	}
	body := SetParameterValuesStruct{
		ParamList:    paramList,
		ParameterKey: msg.ParameterKey,
	}
	env.Body = SetParameterValuesBodyStruct{body}
	output, err := xml.MarshalIndent(env, "  ", "    ")
	//output, err := xml.Marshal(env)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	return output
}

func (msg *SetParameterValues) Parse(xmlstr string) {

}
