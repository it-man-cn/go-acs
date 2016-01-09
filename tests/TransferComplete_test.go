package test

import (
	"fmt"
	"go-acs/models/messages"
	"testing"
)

func TestCreateTransferComplete(t *testing.T) {
	resp := new(messages.TransferComplete)
	resp.CommandKey = "abc"
	resp.StartTime = "2015-02-12T13:40:07"
	resp.CompleteTime = "2015-02-12T13:40:07"
	resp.FaultCode = 1
	resp.FaultString = "error"
	fmt.Println(string(resp.CreateXml()))
}

func TestParseTransferComplete(t *testing.T) {

	data := `<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:SOAP-ENC="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:cwmp="urn:dslforum-org:cwmp-1-0">
	      <SOAP-ENV:Header>
	          <cwmp:ID SOAP-ENV:mustUnderstand="1">ID:intrnl.unset.id.TransferComplete1439547139.1439547139028717609</cwmp:ID>
	      </SOAP-ENV:Header>
	      <SOAP-ENV:Body>
	          <cwmp:TransferComplete>
	              <CommandKey>abc</CommandKey>
	              <StartTime>2015-02-12T13:40:07</StartTime>
	              <CompleteTime>2015-02-12T13:40:07</CompleteTime>
	              <FaultStruct>
	                  <FaultCode>1</FaultCode>
	                  <FaultString>error</FaultString>
	              </FaultStruct>
	          </cwmp:TransferComplete>
	      </SOAP-ENV:Body>
	  </SOAP-ENV:Envelope>`
	msg, _ := messages.ParseXml(data)
	resp := msg.(*messages.TransferComplete)
	fmt.Println(resp.CommandKey)
	fmt.Println(resp.StartTime)
	fmt.Println(resp.CompleteTime)
	fmt.Println(resp.FaultCode)
	fmt.Println(resp.FaultString)
}
