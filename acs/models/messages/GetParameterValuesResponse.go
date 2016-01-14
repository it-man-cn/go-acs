package messages

import (
	"encoding/xml"
	"fmt"
	"github.com/coraldane/godom"
	"strconv"
	"time"
)

type GetParameterValuesResponse struct {
	Id     string
	Name   string
	Values map[string]string
}

func NewGetParameterValuesResponse() (m *GetParameterValuesResponse) {
	m = &GetParameterValuesResponse{}
	m.Id = m.GetId()
	m.Name = m.GetName()
	return m
}

type GetParameterValuesResponseBodyStruct struct {
	Body GetParameterValuesResponseStruct `xml:"cwmp:GetParameterValuesResponse"`
}

type GetParameterValuesResponseStruct struct {
	Params ParameterListStruct `xml:"ParameterList"`
}

func (msg *GetParameterValuesResponse) GetName() string {
	return "GetParameterValuesResponse"
}

func (msg *GetParameterValuesResponse) GetId() string {
	if len(msg.Id) < 1 {
		msg.Id = fmt.Sprintf("ID:intrnl.unset.id.%s%d.%d", msg.GetName(), time.Now().Unix(), time.Now().UnixNano())
	}
	return msg.Id
}

func (msg *GetParameterValuesResponse) CreateXml() []byte {
	env := Envelope{}
	id := IdStruct{"1", msg.GetId()}
	env.XmlnsEnv = "http://schemas.xmlsoap.org/soap/envelope/"
	env.XmlnsEnc = "http://schemas.xmlsoap.org/soap/encoding/"
	env.XmlnsXsd = "http://www.w3.org/2001/XMLSchema"
	env.XmlnsXsi = "http://www.w3.org/2001/XMLSchema-instance"
	env.XmlnsCwmp = "urn:dslforum-org:cwmp-1-0"
	env.Header = HeaderStruct{ID: id}

	paramLen := strconv.Itoa(len(msg.Values))
	params := ParameterListStruct{Type: "cwmp:ParameterValueStruct[" + paramLen + "]"}
	for k, v := range msg.Values {
		param := ParameterValueStruct{
			Name:  NodeStruct{Type: XSD_STRING, Value: k},
			Value: NodeStruct{Type: XSD_STRING, Value: v}}
		params.Params = append(params.Params, param)
	}
	info := GetParameterValuesResponseStruct{Params: params}
	env.Body = GetParameterValuesResponseBodyStruct{info}
	output, err := xml.MarshalIndent(env, "  ", "    ")
	//output, err := xml.Marshal(env)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	return output
}

func (msg *GetParameterValuesResponse) Parse(xmlstr string) {
	document, _ := dom.ParseString(xmlstr)
	root := document.DocumentElement()
	hdr := root.GetElementsByTagName("Header")
	if hdr.Length() > 0 {
		pNode := hdr.Item(0)
		msg.Id = GetChildElementValue(pNode, "ID")
	}

	paramList := root.GetElementsByTagName("ParameterValueStruct")
	if paramList.Length() > 0 {
		params := make(map[string]string)
		for i := uint(0); i < paramList.Length(); i++ {
			pi := paramList.Item(i)
			if pi.NodeType() == dom.ELEMENT_NODE {
				nodes := pi.ChildNodes()
				var name, value string
				for j := uint(0); j < nodes.Length(); j++ {
					node := nodes.Item(j)
					if node.NodeType() == dom.ELEMENT_NODE {
						if "Name" == node.NodeName() && node.HasChildNodes() {
							name = node.FirstChild().NodeValue()
						} else {
							if node.HasChildNodes() {
								value = node.FirstChild().NodeValue()
							} else {
								value = ""
							}
						}
					} else {
						continue
					}
				}
				params[name] = value
			}

		}
		msg.Values = params
	}
}
