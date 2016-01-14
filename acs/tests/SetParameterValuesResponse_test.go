package test

import (
	"fmt"
	"go-acs/acs/models/messages"
	"testing"
)

func TestCreateSetParameterValuesResponse(t *testing.T) {
	resp := new(messages.SetParameterValuesResponse)
	fmt.Println(string(resp.CreateXml()))
}

func TestParseSetParameterValuesResponse(t *testing.T) {

	data := `<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:SOAP-ENC="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:cwmp="urn:dslforum-org:cwmp-1-0">
	      <SOAP-ENV:Header>
	          <cwmp:ID SOAP-ENV:mustUnderstand="1">ID:intrnl.unset.id.SetParameterValuesResponse1439546279.1439546279630058749</cwmp:ID>
	      </SOAP-ENV:Header>
	      <SOAP-ENV:Body>
	          <cwmp:SetParameterValuesResponse>
	              <Status>1</Status>
	              <ParameterKey>test</ParameterKey>
	          </cwmp:SetParameterValuesResponse>
	      </SOAP-ENV:Body>
	  </SOAP-ENV:Envelope>`
	msg, _ := messages.ParseXml(data)
	resp := msg.(*messages.SetParameterValuesResponse)
	fmt.Println(resp.Status)
	fmt.Println(resp.ParameterKey)

}
