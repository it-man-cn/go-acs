package test

import (
	"fmt"
	"go-acs/acs/models/messages"
	"testing"
)

func TestCreateGetParameterValuesResponse(t *testing.T) {
	resp := new(messages.GetParameterValuesResponse)
	params := make(map[string]string)
	params["InternetGatewayDevice.DeviceInfo.Manufacturer"] = "ACS"
	params["InternetGatewayDevice.DeviceInfo.OUI"] = "0011AB"
	resp.Params = params
	fmt.Println(string(resp.CreateXml()))
}

func TestParseGetParameterValuesResponse(t *testing.T) {
	data := `<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:SOAP-ENC="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:cwmp="urn:dslforum-org:cwmp-1-0">
	      <SOAP-ENV:Header>
	          <cwmp:ID SOAP-ENV:mustUnderstand="1">ID:intrnl.unset.id.GetParameterValuesResponse1439545078.1439545078160785199</cwmp:ID>
	      </SOAP-ENV:Header>
	      <SOAP-ENV:Body>
	          <cwmp:GetParameterValuesResponse>
	              <ParameterList SOAP-ENC:arrayType="cwmp:ParameterValueStruct[2]">
	                  <ParameterValueStruct>
	                      <Name xsi:type="xsd:string">InternetGatewayDevice.DeviceInfo.Manufacturer</Name>
	                      <Value xsi:type="xsd:string">ACS</Value>
	                  </ParameterValueStruct>
	                  <ParameterValueStruct>
	                      <Name xsi:type="xsd:string">InternetGatewayDevice.DeviceInfo.OUI</Name>
	                      <Value xsi:type="xsd:string">0011AB</Value>
	                  </ParameterValueStruct>
	              </ParameterList>
	          </cwmp:GetParameterValuesResponse>
	      </SOAP-ENV:Body>
	  </SOAP-ENV:Envelope>`
	msg, _ := messages.ParseXml(data)
	resp := msg.(*messages.GetParameterValuesResponse)
	for k, v := range resp.Params {
		fmt.Println(k, v)
	}
}
