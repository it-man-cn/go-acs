package messages

import (
	"encoding/xml"
	"fmt"
	"github.com/coraldane/godom"
	"strconv"
	"time"
)

//Inform tr069 inform (heartbeat)
type Inform struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Manufacturer string            `json:"manufacturer"`
	OUI          string            `json:"oui"`
	ProductClass string            `json:"productClass"`
	Sn           string            `json:"sn"`
	Events       map[string]string `json:"events"`
	MaxEnvelopes int               `json:"maxEnvelopes"`
	CurrentTime  string            `json:"currentTime"`
	RetryCount   int               `json:"retryCount"`
	Params       map[string]string `json:"params"`
}

type informBodyStruct struct {
	Body informStruct `xml:"cwmp:Inform"`
}

type informStruct struct {
	DeviceID     deviceIDStruct      `xml:"DeviceId"`
	Event        EventStruct         `xml:"Event"`
	MaxEnvelopes NodeStruct          `xml:"MaxEnvelopes"`
	CurrentTime  NodeStruct          `xml:"CurrentTime"`
	RetryCount   NodeStruct          `xml:"RetryCount"`
	Params       ParameterListStruct `xml:"ParameterList"`
}

type deviceIDStruct struct {
	Type         string     `xml:"xsi:type,attr"`
	Manufacturer NodeStruct `xml:"Manufacturer"`
	OUI          NodeStruct `xml:"OUI"`
	ProductClass NodeStruct `xml:"ProductClass"`
	SerialNumber NodeStruct `xml:"SerialNumber"`
}

//GetName get msg type
func (msg *Inform) GetName() string {
	return "Inform"
}

//GetID get msg id
func (msg *Inform) GetID() string {
	if len(msg.ID) < 1 {
		msg.ID = fmt.Sprintf("ID:intrnl.unset.id.%s%d.%d", msg.GetName(), time.Now().Unix(), time.Now().UnixNano())
	}
	return msg.ID
}

//CreateXML encode into xml
func (msg *Inform) CreateXML() []byte {
	env := Envelope{}
	id := IDStruct{"1", msg.GetID()}
	env.XmlnsEnv = "http://schemas.xmlsoap.org/soap/envelope/"
	env.XmlnsEnc = "http://schemas.xmlsoap.org/soap/encoding/"
	env.XmlnsXsd = "http://www.w3.org/2001/XMLSchema"
	env.XmlnsXsi = "http://www.w3.org/2001/XMLSchema-instance"
	env.XmlnsCwmp = "urn:dslforum-org:cwmp-1-0"
	env.Header = HeaderStruct{ID: id}
	manufacturer := NodeStruct{Type: XsdString, Value: msg.Manufacturer}
	oui := NodeStruct{Type: XsdString, Value: msg.OUI}
	productClass := NodeStruct{Type: XsdString, Value: msg.ProductClass}
	serialNumber := NodeStruct{Type: XsdString, Value: msg.Sn}
	deviceID := deviceIDStruct{Type: "cwmp:DeviceIdStruct", Manufacturer: manufacturer, OUI: oui, ProductClass: productClass, SerialNumber: serialNumber}
	eventLen := strconv.Itoa(len(msg.Events))
	event := EventStruct{Type: "cwmp:EventStruct[" + eventLen + "]"}
	for k, v := range msg.Events {
		eventCode := NodeStruct{Type: XsdString, Value: k}
		event.Events = append(event.Events, EventNodeStruct{EventCode: eventCode, CommandKey: v})
	}

	maxEnv := strconv.Itoa(msg.MaxEnvelopes)
	maxEnvelopes := NodeStruct{Type: XsdString, Value: maxEnv}
	currentTime := NodeStruct{Type: XsdString, Value: msg.CurrentTime}
	trys := strconv.Itoa(msg.RetryCount)
	retryCount := NodeStruct{Type: XsdString, Value: trys}
	paramLen := strconv.Itoa(len(msg.Params))
	paramList := ParameterListStruct{Type: "cwmp:ParameterValueStruct[" + paramLen + "]"}
	for k, v := range msg.Params {
		param := ParameterValueStruct{
			Name:  NodeStruct{Type: XsdString, Value: k},
			Value: NodeStruct{Type: XsdString, Value: v}}
		paramList.Params = append(paramList.Params, param)
	}
	info := informStruct{DeviceID: deviceID, Event: event, MaxEnvelopes: maxEnvelopes, CurrentTime: currentTime, RetryCount: retryCount, Params: paramList}
	env.Body = informBodyStruct{info}
	output, err := xml.MarshalIndent(env, "  ", "    ")
	//output, err := xml.Marshal(env)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	return output
}

//Parse decode from xml
func (msg *Inform) Parse(xmlstr string) {
	document, _ := dom.ParseString(xmlstr)
	root := document.DocumentElement()
	hdr := root.GetElementsByTagName("Header")
	if hdr.Length() > 0 {
		pNode := hdr.Item(0)
		msg.ID = GetChildElementValue(pNode, "ID")
	}
	deviceid := root.GetElementsByTagName("DeviceId")
	if deviceid.Length() > 0 {
		msg.Manufacturer = GetChildElementValue(deviceid.Item(0), "Manufacturer")
		msg.OUI = GetChildElementValue(deviceid.Item(0), "OUI")
		msg.ProductClass = GetChildElementValue(deviceid.Item(0), "ProductClass")
		msg.Sn = GetChildElementValue(deviceid.Item(0), "SerialNumber")
	}
	info := root.GetElementsByTagName("Inform")
	if info.Length() > 0 {
		msg.MaxEnvelopes, _ = strconv.Atoi(GetChildElementValue(info.Item(0), "MaxEnvelopes"))
		msg.CurrentTime = GetChildElementValue(info.Item(0), "CurrentTime")
		msg.RetryCount, _ = strconv.Atoi(GetChildElementValue(info.Item(0), "RetryCount"))
	}
	event := root.GetElementsByTagName("Event")
	if event.Length() > 0 && event.Item(0).HasChildNodes() {
		events := make(map[string]string)
		nodes := event.Item(0).ChildNodes()
		for i := uint(0); i < nodes.Length(); i++ {
			code := GetChildElementValue(nodes.Item(i), "EventCode")
			key := GetChildElementValue(nodes.Item(i), "CommandKey")
			events[code] = key
		}
		msg.Events = events
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
		msg.Params = params
	}

}

//IsEvent is a connect request or others
func (msg *Inform) IsEvent(event string) bool {
	/*
		for k,_:= range msg.Events {
			if k == event {
				return true
			}
		}
	*/
	if _, ok := msg.Events[event]; ok {
		return true
	}
	return false
}

//GetParam get param in inform
func (msg *Inform) GetParam(name string) (value string) {
	/*
		for k, v := range msg.Params {
			if k == name {
				value = v
				break
			}
		}
	*/
	value = msg.Params[name]
	return
}

//GetConfigVersion get current config version
func (msg *Inform) GetConfigVersion() (version string) {
	version = msg.GetParam("InternetGatewayDevice.DeviceConfig.ConfigVersion")
	return
}
