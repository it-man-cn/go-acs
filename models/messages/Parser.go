package messages

import (
	"github.com/coraldane/godom"
)

func ParseXml(xmlstr string) (Message, error) {
	document, err := dom.ParseString(xmlstr)
	root := document.DocumentElement()
	body := root.GetElementsByTagName("Body")
	var msg Message
	if body.Length() > 0 {
		node := GetFirstChild(body.Item(0))
		if node != nil {
			name := node.NodeName()
			switch name {
			case "Inform":
				msg = &Inform{}
				msg.Parse(xmlstr)
			case "GetParameterValuesResponse":
				msg = &GetParameterValuesResponse{}
				msg.Parse(xmlstr)
			case "SetParameterValuesResponse":
				msg = &SetParameterValuesResponse{}
				msg.Parse(xmlstr)
			case "DownloadResponse":
				msg = &DownloadResponse{}
				msg.Parse(xmlstr)
			case "TransferComplete":
				msg = &TransferComplete{}
				msg.Parse(xmlstr)
			case "GetRPCMethodsResponse":
				msg = &GetRPCMethodsResponse{}
				msg.Parse(xmlstr)
			case "RebootResponse":
				msg = &RebootResponse{}
				msg.Parse(xmlstr)
			}
		}
	}
	return msg, err
}

func GetChildElementValue(pNode dom.Node, name string) string {
	if pNode.NodeType() == dom.ELEMENT_NODE && pNode.HasChildNodes() {
		nodes := pNode.ChildNodes()
		for i := uint(0); i < nodes.Length(); i++ {
			node := nodes.Item(i)
			if node.NodeName() == name {
				if node.HasChildNodes() {
					return node.FirstChild().NodeValue()
				} else {
					return ""
				}
			}
		}
	}
	return ""
}

func GetFirstChild(pNode dom.Node) dom.Node {
	nodes := pNode.ChildNodes()
	for i := uint(0); i < nodes.Length(); i++ {
		node := nodes.Item(i)
		if node.NodeType() == dom.ELEMENT_NODE {
			return node
		}
	}
	return nil
}
