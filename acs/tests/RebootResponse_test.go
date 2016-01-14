package test

import (
	"fmt"
	"go-acs/acs/models/messages"
	"testing"
)

func TestCreateRebootResponse(t *testing.T) {

}

func TestParseRebootResponse(t *testing.T) {
	data := `<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:SOAP-ENC="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsi="http://www.w3.org/1999/XMLSchema-instance" xmlns:xsd="http://www.w3.org/1999/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0">
<SOAP-ENV:Header>
<cwmp:ID SOAP-ENV:mustUnderstand="1">ID:intrnl.unset.id.Reboot1439974339781.14572016</cwmp:ID>
</SOAP-ENV:Header>
<SOAP-ENV:Body>
<cwmp:RebootResponse></cwmp:RebootResponse>
</SOAP-ENV:Body>
</SOAP-ENV:Envelope>`
	msg, _ := messages.ParseXml(data)
	fmt.Println(msg)
	//resp := msg.(*messages.RebootResponse)
}
