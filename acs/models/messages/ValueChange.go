package messages

type ValueChange struct {
	Sn    string `json:"sn"`
	Names []string
}

func (msg *ValueChange) GetName() string {
	return "ValueChange"
}

func (msg *ValueChange) GetId() string {
	return "ValueChange"
}

func (msg *ValueChange) CreateXml() (xml []byte) {
	return xml
}

func (msg *ValueChange) Parse(xmlstr string) {

}
