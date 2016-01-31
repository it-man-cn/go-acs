package messages

import (
	"encoding/xml"
	"fmt"
	"github.com/coraldane/godom"
	"strconv"
	"time"
)

//GetParameterValuesResponse getParameterValues response
type GetParameterValuesResponse struct {
	ID     string
	Name   string
	Values map[string]string
}

//NewGetParameterValuesResponse create GetParameterValuesResponse object
func NewGetParameterValuesResponse() (m *GetParameterValuesResponse) {
	m = &GetParameterValuesResponse{}
	m.ID = m.GetID()
	m.Name = m.GetName()
	return m
}

type getParameterValuesResponseBodyStruct struct {
	Body getParameterValuesResponseStruct `xml:"cwmp:GetParameterValuesResponse"`
}

type getParameterValuesResponseStruct struct {
	Params ParameterListStruct `xml:"ParameterList"`
}

//GetName get type name
func (msg *GetParameterValuesResponse) GetName() string {
	return "GetParameterValuesResponse"
}

//GetID get msg id
func (msg *GetParameterValuesResponse) GetID() string {
	if len(msg.ID) < 1 {
		msg.ID = fmt.Sprintf("ID:intrnl.unset.id.%s%d.%d", msg.GetName(), time.Now().Unix(), time.Now().UnixNano())
	}
	return msg.ID
}

//CreateXML encode into xml
func (msg *GetParameterValuesResponse) CreateXML() []byte {
	env := Envelope{}
	id := IDStruct{"1", msg.GetID()}
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
			Name:  NodeStruct{Type: XsdString, Value: k},
			Value: NodeStruct{Type: XsdString, Value: v}}
		params.Params = append(params.Params, param)
	}
	info := getParameterValuesResponseStruct{Params: params}
	env.Body = getParameterValuesResponseBodyStruct{info}
	output, err := xml.MarshalIndent(env, "  ", "    ")
	//output, err := xml.Marshal(env)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	return output
}

//Parse decode from xml
func (msg *GetParameterValuesResponse) Parse(xmlstr string) {
	document, _ := dom.ParseString(xmlstr)
	root := document.DocumentElement()
	hdr := root.GetElementsByTagName("Header")
	if hdr.Length() > 0 {
		pNode := hdr.Item(0)
		msg.ID = GetChildElementValue(pNode, "ID")
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
