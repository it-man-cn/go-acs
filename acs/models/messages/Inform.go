package messages

import (
	"encoding/xml"
	"fmt"
	"github.com/coraldane/godom"
	"strconv"
	"time"
)

type Inform struct {
	Id           string            `json:"id"`
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

type InformBodyStruct struct {
	Body InformStruct `xml:"cwmp:Inform"`
}

type InformStruct struct {
	DeviceId     DeviceIdStruct      `xml:"DeviceId"`
	Event        EventStruct         `xml:"Event"`
	MaxEnvelopes NodeStruct          `xml:"MaxEnvelopes"`
	CurrentTime  NodeStruct          `xml:"CurrentTime"`
	RetryCount   NodeStruct          `xml:"RetryCount"`
	Params       ParameterListStruct `xml:"ParameterList"`
}

type DeviceIdStruct struct {
	Type         string     `xml:"xsi:type,attr"`
	Manufacturer NodeStruct `xml:"Manufacturer"`
	OUI          NodeStruct `xml:"OUI"`
	ProductClass NodeStruct `xml:"ProductClass"`
	SerialNumber NodeStruct `xml:"SerialNumber"`
}

func (msg *Inform) GetName() string {
	return "Inform"
}

func (msg *Inform) GetId() string {
	if len(msg.Id) < 1 {
		msg.Id = fmt.Sprintf("ID:intrnl.unset.id.%s%d.%d", msg.GetName(), time.Now().Unix(), time.Now().UnixNano())
	}
	return msg.Id
}

func (msg *Inform) CreateXml() []byte {
	env := Envelope{}
	id := IdStruct{"1", msg.GetId()}
	env.XmlnsEnv = "http://schemas.xmlsoap.org/soap/envelope/"
	env.XmlnsEnc = "http://schemas.xmlsoap.org/soap/encoding/"
	env.XmlnsXsd = "http://www.w3.org/2001/XMLSchema"
	env.XmlnsXsi = "http://www.w3.org/2001/XMLSchema-instance"
	env.XmlnsCwmp = "urn:dslforum-org:cwmp-1-0"
	env.Header = HeaderStruct{ID: id}
	manufacturer := NodeStruct{Type: XSD_STRING, Value: msg.Manufacturer}
	oui := NodeStruct{Type: XSD_STRING, Value: msg.OUI}
	productClass := NodeStruct{Type: XSD_STRING, Value: msg.ProductClass}
	serialNumber := NodeStruct{Type: XSD_STRING, Value: msg.Sn}
	deviceId := DeviceIdStruct{Type: "cwmp:DeviceIdStruct", Manufacturer: manufacturer, OUI: oui, ProductClass: productClass, SerialNumber: serialNumber}
	eventLen := strconv.Itoa(len(msg.Events))
	event := EventStruct{Type: "cwmp:EventStruct[" + eventLen + "]"}
	for k, v := range msg.Events {
		eventCode := NodeStruct{Type: XSD_STRING, Value: k}
		event.Events = append(event.Events, EventNodeStruct{EventCode: eventCode, CommandKey: v})
	}

	maxEnv := strconv.Itoa(msg.MaxEnvelopes)
	maxEnvelopes := NodeStruct{Type: XSD_UNSIGNEDINT, Value: maxEnv}
	currentTime := NodeStruct{Type: XSD_STRING, Value: msg.CurrentTime}
	trys := strconv.Itoa(msg.RetryCount)
	retryCount := NodeStruct{Type: XSD_UNSIGNEDINT, Value: trys}
	paramLen := strconv.Itoa(len(msg.Params))
	paramList := ParameterListStruct{Type: "cwmp:ParameterValueStruct[" + paramLen + "]"}
	for k, v := range msg.Params {
		param := ParameterValueStruct{
			Name:  NodeStruct{Type: XSD_STRING, Value: k},
			Value: NodeStruct{Type: XSD_STRING, Value: v}}
		paramList.Params = append(paramList.Params, param)
	}
	info := InformStruct{DeviceId: deviceId, Event: event, MaxEnvelopes: maxEnvelopes, CurrentTime: currentTime, RetryCount: retryCount, Params: paramList}
	env.Body = InformBodyStruct{info}
	output, err := xml.MarshalIndent(env, "  ", "    ")
	//output, err := xml.Marshal(env)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	return output
}

func (msg *Inform) Parse(xmlstr string) {
	document, _ := dom.ParseString(xmlstr)
	root := document.DocumentElement()
	hdr := root.GetElementsByTagName("Header")
	if hdr.Length() > 0 {
		pNode := hdr.Item(0)
		msg.Id = GetChildElementValue(pNode, "ID")
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

func (msg *Inform) IsEvent(event string) bool {
	for k, _ := range msg.Events {
		if k == event {
			return true
		}
	}
	return false
}

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

func (msg *Inform) GetConfigVersion() (version string) {
	version = msg.GetParam("InternetGatewayDevice.DeviceConfig.ConfigVersion")
	return
}
