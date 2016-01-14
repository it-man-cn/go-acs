package test

import (
	"fmt"
	"go-acs/acs/models/messages"
	"testing"
)

func TestCreateGetRPCMethodsResponse(t *testing.T) {
	resp := new(messages.GetRPCMethodsResponse)
	methods := make([]string, 0)
	methods = append(methods, "GetRPCMethods", "GetParameterNames")
	resp.Methods = methods
	fmt.Println(string(resp.CreateXml()))
}

func TestParseGetRPCMethodsResponse(t *testing.T) {
	data := `<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:SOAP-ENC="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:cwmp="urn:dslforum-org:cwmp-1-0">
	      <SOAP-ENV:Header>
	          <cwmp:ID SOAP-ENV:mustUnderstand="1">ID:intrnl.unset.id.GetRPCMethodsResponse1439556667.1439556667543994313</cwmp:ID>
	      </SOAP-ENV:Header>
	      <SOAP-ENV:Body>
	          <cwmp:GetRPCMethodsResponse>
	              <cwmp:MethodList xsi:type="SOAP-ENC:Array" SOAP-ENC:arrayType="xsd:string[2]">
	                  <string>GetRPCMethods</string>
	                  <string>GetParameterNames</string>
	              </cwmp:MethodList>
	          </cwmp:GetRPCMethodsResponse>
	      </SOAP-ENV:Body>
	  </SOAP-ENV:Envelope>`
	msg, _ := messages.ParseXml(data)
	resp := msg.(*messages.GetRPCMethodsResponse)
	for _, v := range resp.Methods {
		fmt.Println(v)
	}
}
