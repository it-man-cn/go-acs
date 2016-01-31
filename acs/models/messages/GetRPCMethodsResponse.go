package messages

import (
	"encoding/xml"
	"fmt"
	"github.com/coraldane/godom"
	"strconv"
	"time"
)

//GetRPCMethodsResponse getRPCMethods reponse
type GetRPCMethodsResponse struct {
	ID      string
	Name    string
	Methods []string
}

type getRPCMethodsResponseBodyStruct struct {
	Body getRPCMethodsResponseStruct `xml:"cwmp:GetRPCMethodsResponse"`
}

type getRPCMethodsResponseStruct struct {
	MethodList methodListStruct `xml:"cwmp:MethodList"`
}

type methodListStruct struct {
	Type      string   `xml:"xsi:type,attr"`
	ArrayType string   `xml:"SOAP-ENC:arrayType,attr"`
	Methods   []string `xml:"string"`
}

//GetName get msg type
func (msg *GetRPCMethodsResponse) GetName() string {
	return "GetRPCMethodsResponse"
}

//GetID get msg id
func (msg *GetRPCMethodsResponse) GetID() string {
	if len(msg.ID) < 1 {
		msg.ID = fmt.Sprintf("ID:intrnl.unset.id.%s%d.%d", msg.GetName(), time.Now().Unix(), time.Now().UnixNano())
	}
	return msg.ID
}

//CreateXML encode into xml
func (msg *GetRPCMethodsResponse) CreateXML() []byte {
	env := Envelope{}
	id := IDStruct{"1", msg.GetID()}
	env.XmlnsEnv = "http://schemas.xmlsoap.org/soap/envelope/"
	env.XmlnsEnc = "http://schemas.xmlsoap.org/soap/encoding/"
	env.XmlnsXsd = "http://www.w3.org/2001/XMLSchema"
	env.XmlnsXsi = "http://www.w3.org/2001/XMLSchema-instance"
	env.XmlnsCwmp = "urn:dslforum-org:cwmp-1-0"
	env.Header = HeaderStruct{ID: id}
	methodsLen := strconv.Itoa(len(msg.Methods))
	methodList := methodListStruct{
		Type:      SoapArray,
		ArrayType: XsdString + "[" + methodsLen + "]",
	}
	for _, v := range msg.Methods {
		methodList.Methods = append(methodList.Methods, v)
	}
	body := getRPCMethodsResponseStruct{methodList}
	env.Body = getRPCMethodsResponseBodyStruct{body}
	output, err := xml.MarshalIndent(env, "  ", "    ")
	//output, err := xml.Marshal(env)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	return output
}

//Parse decode from xml
func (msg *GetRPCMethodsResponse) Parse(xmlstr string) {
	document, _ := dom.ParseString(xmlstr)
	root := document.DocumentElement()
	hdr := root.GetElementsByTagName("Header")
	if hdr.Length() > 0 {
		pNode := hdr.Item(0)
		msg.ID = GetChildElementValue(pNode, "ID")
	}

	methodList := root.GetElementsByTagName("MethodList")
	if methodList.Length() > 0 {
		var methods []string
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
