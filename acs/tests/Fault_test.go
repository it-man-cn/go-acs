package test

import (
	"fmt"
	"go-acs/acs/models/messages"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestCreateFault(t *testing.T) {
	fault := new(messages.Fault)
	fault.ID = "ID:intrnl.unset.id.SetParameterValues1458897112421.1852997967"
	fault.FaultCode = "Client"
	fault.FaultString = "Client fault"
	fault.CwmpFaultCode = "9003"
	fault.CwmpFaultString = "Invalid arguments"
	fault.SetParameterValuesFault = messages.SetParameterValuesFaultStruct{
		ParameterName: "InternetGatewayDevice.config.webauthglobal.successUrl",
		FaultCode:     "9008",
		FaultString:   "Attempt to set a non-writable parameter",
		ParameterKey:  "006aecb92b10459b91321e32c9b4502f",
	}
	fmt.Println(string(fault.CreateXML()))
}

func TestParseFault(t *testing.T) {
	data := `<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/"
	       xmlns:SOAP-ENC="http://schemas.xmlsoap.org/soap/encoding/"
	       xmlns:xsi="http://www.w3.org/1999/XMLSchema-instance"
	       xmlns:xsd="http://www.w3.org/1999/XMLSchema"
	       xmlns:cwmp="urn:dslforum-org:cwmp-1-0">
	       <SOAP-ENV:Header>
	       <cwmp:ID SOAP-ENV:mustUnderstand="1">ID:intrnl.unset.id.SetParameterValues1458897112421.1852997967</cwmp:ID>
	       </SOAP-ENV:Header>
	       <SOAP-ENV:Body>
	       <SOAP-ENV:Fault>
	       <faultcode>Client</faultcode>
	       <faultstring>Client fault</faultstring>
	       <detail>
	       <cwmp:Fault>
	       <FaultCode>9003</FaultCode>
	       <FaultString>Invalid arguments</FaultString>
	       <SetParameterValuesFault>
	       <ParameterName>InternetGatewayDevice.config.webauthglobal.successUrl</ParameterName>
	       <FaultCode>9008</FaultCode>
	       <FaultString>Attempt to set a non-writable parameter</FaultString>
	       <ParameterKey>006aecb92b10459b91321e32c9b4502f</ParameterKey>
	       </SetParameterValuesFault>
	       </cwmp:Fault>
	       </detail>
	       </SOAP-ENV:Fault>
	       </SOAP-ENV:Body>
	       </SOAP-ENV:Envelope>`
	msg, _ := messages.ParseXML([]byte(data))
	fault := msg.(*messages.Fault)
	fmt.Println(fault.FaultCode)
	fmt.Println(fault.FaultString)
	fmt.Println(fault.CwmpFaultCode)
	fmt.Println(fault.CwmpFaultString)
	fmt.Println(fault.SetParameterValuesFault.FaultCode)
	fmt.Println(fault.SetParameterValuesFault.FaultString)
	fmt.Println(fault.SetParameterValuesFault.ParameterName)
	fmt.Println(fault.SetParameterValuesFault.ParameterKey)
}

func TestSendFault(t *testing.T) {
	body := `<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/"
	       xmlns:SOAP-ENC="http://schemas.xmlsoap.org/soap/encoding/"
	       xmlns:xsi="http://www.w3.org/1999/XMLSchema-instance"
	       xmlns:xsd="http://www.w3.org/1999/XMLSchema"
	       xmlns:cwmp="urn:dslforum-org:cwmp-1-0">
	       <SOAP-ENV:Header>
	       <cwmp:ID SOAP-ENV:mustUnderstand="1">ID:intrnl.unset.id.SetParameterValues1458897112421.1852997967</cwmp:ID>
	       </SOAP-ENV:Header>
	       <SOAP-ENV:Body>
	       <SOAP-ENV:Fault>
	       <faultcode>Client</faultcode>
	       <faultstring>Client fault</faultstring>
	       <detail>
	       <cwmp:Fault>
	       <FaultCode>9003</FaultCode>
	       <FaultString>Invalid arguments</FaultString>
	       <SetParameterValuesFault>
	       <ParameterName>InternetGatewayDevice.config.webauthglobal.successUrl</ParameterName>
	       <FaultCode>9008</FaultCode>
	       <FaultString>Attempt to set a non-writable parameter</FaultString>
	       <ParameterKey>006aecb92b10459b91321e32c9b4502f</ParameterKey>
	       </SetParameterValuesFault>
	       </cwmp:Fault>
	       </detail>
	       </SOAP-ENV:Fault>
	       </SOAP-ENV:Body>
	       </SOAP-ENV:Envelope>`

	//resp, err := http.Post("http://acs.greenwifi.com.cn/ACS/tr069",
	resp, err := http.Post("http://localhost:10091/ACS/tr069",
		"text/xml",
		strings.NewReader(body))
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		fmt.Println(err)
	}

	fmt.Println(string(respBody))
}
