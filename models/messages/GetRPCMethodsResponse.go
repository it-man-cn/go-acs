package messages

import (
	"encoding/xml"
	"fmt"
	"github.com/coraldane/godom"
	"strconv"
	"time"
)

type GetRPCMethodsResponse struct {
	Id      string
	Name    string
	Methods []string
}

type GetRPCMethodsResponseBodyStruct struct {
	Body GetRPCMethodsResponseStruct `xml:"cwmp:GetRPCMethodsResponse"`
}

type GetRPCMethodsResponseStruct struct {
	MethodList MethodListStruct `xml:"cwmp:MethodList"`
}

type MethodListStruct struct {
	Type      string   `xml:"xsi:type,attr"`
	ArrayType string   `xml:"SOAP-ENC:arrayType,attr"`
	Methods   []string `xml:"string"`
}

func (msg *GetRPCMethodsResponse) GetName() string {
	return "GetRPCMethodsResponse"
}

func (msg *GetRPCMethodsResponse) GetId() string {
	if len(msg.Id) < 1 {
		msg.Id = fmt.Sprintf("ID:intrnl.unset.id.%s%d.%d", msg.GetName(), time.Now().Unix(), time.Now().UnixNano())
	}
	return msg.Id
}

func (msg *GetRPCMethodsResponse) CreateXml() []byte {
	env := Envelope{}
	id := IdStruct{"1", msg.GetId()}
	env.XmlnsEnv = "http://schemas.xmlsoap.org/soap/envelope/"
	env.XmlnsEnc = "http://schemas.xmlsoap.org/soap/encoding/"
	env.XmlnsXsd = "http://www.w3.org/2001/XMLSchema"
	env.XmlnsXsi = "http://www.w3.org/2001/XMLSchema-instance"
	env.XmlnsCwmp = "urn:dslforum-org:cwmp-1-0"
	env.Header = HeaderStruct{ID: id}
	methodsLen := strconv.Itoa(len(msg.Methods))
	methodList := MethodListStruct{
		Type:      SOAP_ARRAY,
		ArrayType: XSD_STRING + "[" + methodsLen + "]",
	}
	for _, v := range msg.Methods {
		methodList.Methods = append(methodList.Methods, v)
	}
	body := GetRPCMethodsResponseStruct{methodList}
	env.Body = GetRPCMethodsResponseBodyStruct{body}
	output, err := xml.MarshalIndent(env, "  ", "    ")
	//output, err := xml.Marshal(env)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	return output
}

func (msg *GetRPCMethodsResponse) Parse(xmlstr string) {
	document, _ := dom.ParseString(xmlstr)
	root := document.DocumentElement()
	hdr := root.GetElementsByTagName("Header")
	if hdr.Length() > 0 {
		pNode := hdr.Item(0)
		msg.Id = GetChildElementValue(pNode, "ID")
	}

	methodList := root.GetElementsByTagName("MethodList")
	if methodList.Length() > 0 {
		methods := make([]string, 0)
		for i := uint(0); i < methodList.Length(); i++ {
			pi := methodList.Item(i)
			if pi.NodeType() == dom.ELEMENT_NODE {
				nodes := pi.ChildNodes()
				var name string
				for j := uint(0); j < nodes.Length(); j++ {
					node := nodes.Item(j)
					name = ""
					if node.NodeType() == dom.ELEMENT_NODE {
						if "string" == node.NodeName() && node.HasChildNodes() {
							name = node.FirstChild().NodeValue()
							if len(name) > 0 {
								methods = append(methods, name)
							}
						}
					}
				}

			}
		}
		msg.Methods = methods
	}
}
