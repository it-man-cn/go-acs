package messages

type OnlineInform struct {
	Sn    string `json:"sn"`
	Hosts []Host
}

type Host struct {
	Mac      string `json:"mac"`
	HostName string `json:"hostname"`
}

func (msg *OnlineInform) GetName() string {
	return "OnlineInform"
}

func (msg *OnlineInform) GetId() string {
	return "OnlineInform"
}

func (msg *OnlineInform) CreateXml() (xml []byte) {
	return xml
}

func (msg *OnlineInform) Parse(xmlstr string) {

}
