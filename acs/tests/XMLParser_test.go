package tests

import (
	"fmt"
	"go-acs/libs/xml"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	xmlstr := `<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:SOAP-ENC="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:cwmp="urn:dslforum-org:cwmp-1-0">
      <SOAP-ENV:Header>
          <cwmp:ID SOAP-ENV:mustUnderstand="1">ID:intrnl.unset.id.Inform958064891.1658176487</cwmp:ID>
          <cwmp:NoMoreRequests>0</cwmp:NoMoreRequests>
      </SOAP-ENV:Header>
      <SOAP-ENV:Body>
          <cwmp:Inform>
              <DeviceId xsi:type="cwmp:DeviceIdStruct">
                  <Manufacturer xsi:type="xsd:string">UTT</Manufacturer>
                  <OUI xsi:type="xsd:string">OO22AA</OUI>
                  <ProductClass xsi:type="xsd:string">D915W</ProductClass>
                  <SerialNumber xsi:type="xsd:string">14586025</SerialNumber>
              </DeviceId>
              <Event SOAP-ENC:arrayType="cwmp:EventStruct[1]">
                  <EventStruct>
                      <EventCode xsi:type="xsd:string">2 PERIODIC</EventCode>
                      <CommandKey></CommandKey>
                  </EventStruct>
              </Event>
              <MaxEnvelopes xsi:type="xsd:unsignedInt">0</MaxEnvelopes>
              <CurrentTime xsi:type="xsd:string">2015-01-21T14:46:07</CurrentTime>
              <RetryCount xsi:type="xsd:unsignedInt">0</RetryCount>
              <ParameterList SOAP-ENC:arrayType="cwmp:ParameterValueStruct[11]">
                  <ParameterValueStruct>
                      <Name xsi:type="string">InternetGatewayDevice.DeviceInfo.ProvisioningCode</Name>
                      <Value xsi:type="string">D915W</Value>
                  </ParameterValueStruct>
                  <ParameterValueStruct>
                      <Name xsi:type="string">InternetGatewayDevice.ManagementServer.URL</Name>
                      <Value xsi:type="string"></Value>
                  </ParameterValueStruct>
                  <ParameterValueStruct>
                      <Name xsi:type="string">InternetGatewayDevice.ManagementServer.Password</Name>
                      <Value xsi:type="string"></Value>
                  </ParameterValueStruct>
                  <ParameterValueStruct>
                      <Name xsi:type="string">InternetGatewayDevice.ManagementServer.ConnectionRequestURL</Name>
                      <Value xsi:type="string">http://192.168.16.68:5400</Value>
                  </ParameterValueStruct>
                  <ParameterValueStruct>
                      <Name xsi:type="string">InternetGatewayDevice.DeviceSummary</Name>
                      <Value xsi:type="string"></Value>
                  </ParameterValueStruct>
                  <ParameterValueStruct>
                      <Name xsi:type="string">InternetGatewayDevice.DeviceInfo.SpecVersion</Name>
                      <Value xsi:type="string">V1.0</Value>
                  </ParameterValueStruct>
                  <ParameterValueStruct>
                      <Name xsi:type="string">InternetGatewayDevice.ManagementServer.UDPConnectionRequestAddress</Name>
                      <Value xsi:type="string">200.200.202.68:1036</Value>
                  </ParameterValueStruct>
                  <ParameterValueStruct>
                      <Name xsi:type="string">InternetGatewayDevice.ManagementServer.ConnectionRequestPassword</Name>
                      <Value xsi:type="string"></Value>
                  </ParameterValueStruct>
                  <ParameterValueStruct>
                      <Name xsi:type="string">InternetGatewayDevice.DeviceInfo.SoftwareVersion</Name>
                      <Value xsi:type="string">nvD915wv2.2.0-150120</Value>
                  </ParameterValueStruct>
                  <ParameterValueStruct>
                      <Name xsi:type="string">InternetGatewayDevice.WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.ExternalIPAddress</Name>
                      <Value xsi:type="string">192.168.16.68</Value>
                  </ParameterValueStruct>
                  <ParameterValueStruct>
                      <Name xsi:type="string">InternetGatewayDevice.DeviceInfo.HardwareVersion</Name>
                      <Value xsi:type="string">V1.0</Value>
                  </ParameterValueStruct>
              </ParameterList>
          </cwmp:Inform>
      </SOAP-ENV:Body>
  </SOAP-ENV:Envelope>`
	doc := xml.New()
	doc.LoadString(xmlstr, nil)
	//fmt.Println(doc.String())
	//header := doc.SelectNode("*", "Header")
	//fmt.Println(header)
	id := doc.SelectNode("*", "ID")
	fmt.Println("id===:", id.GetValue())
	fmt.Println("id===:", id.Value)
	deviceNode := doc.SelectNode("*", "DeviceId")
	fmt.Println(deviceNode.SelectNode("", "Manufacturer").GetValue())
	fmt.Println(deviceNode.SelectNode("", "OUI").GetValue())
	fmt.Println(deviceNode.SelectNode("", "ProductClass").GetValue())
	fmt.Println(deviceNode.SelectNode("", "SerialNumber").GetValue())
	informNode := doc.SelectNode("*", "Inform")
	fmt.Println(informNode.SelectNode("", "MaxEnvelopes").GetValue())
	fmt.Println(informNode.SelectNode("", "CurrentTime").GetValue())
	fmt.Println(informNode.SelectNode("", "RetryCount").GetValue())
	eventNode := doc.SelectNode("*", "Event")
	fmt.Println("for events:")
	for _, event := range eventNode.Children {
		if len(strings.TrimSpace(event.String())) > 0 {
			fmt.Println(event.SelectNode("", "EventCode").GetValue())
			fmt.Println(event.SelectNode("", "CommandKey").GetValue())
		}
	}

	paramsNode := doc.SelectNode("*", "ParameterList")
	fmt.Println("for params:")
	for _, param := range paramsNode.Children {
		if len(strings.TrimSpace(param.String())) > 0 {
			fmt.Println(param.SelectNode("", "Name").GetValue())
			fmt.Println(param.SelectNode("", "Value").GetValue())
		}
	}
	bodyNode := doc.SelectNode("*", "Body")
	fmt.Println("body:")
	fmt.Println(bodyNode.Children[1].Name.Local)
}

/*
<?xml version="1.0" encoding="utf-8" standalone="yes"?><http://schemas.xmlsoap.org/soap/envelope/:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:SOAP-ENC="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:cwmp="urn:dslforum-org:cwmp-1-0">
      <http://schemas.xmlsoap.org/soap/envelope/:Header>
          <urn:dslforum-org:cwmp-1-0:ID SOAP-ENV:mustUnderstand="1">ID:intrnl.unset.id.Inform958064891.1658176487</urn:dslforum-org:cwmp-1-0:ID>
          <urn:dslforum-org:cwmp-1-0:NoMoreRequests>0</urn:dslforum-org:cwmp-1-0:NoMoreRequests>
      </http://schemas.xmlsoap.org/soap/envelope/:Header>
      <http://schemas.xmlsoap.org/soap/envelope/:Body>
          <urn:dslforum-org:cwmp-1-0:Inform>
              <DeviceId xsi:type="cwmp:DeviceIdStruct">
                  <Manufacturer xsi:type="xsd:string">UTT</Manufacturer>
                  <OUI xsi:type="xsd:string">OO22AA</OUI>
                  <ProductClass xsi:type="xsd:string">D915W</ProductClass>
                  <SerialNumber xsi:type="xsd:string">14586025</SerialNumber>
              </DeviceId>
              <Event SOAP-ENC:arrayType="cwmp:EventStruct[1]">
                  <EventStruct>
                      <EventCode xsi:type="xsd:string">2 PERIODIC</EventCode>
                      <CommandKey />
                  </EventStruct>
              </Event>
              <MaxEnvelopes xsi:type="xsd:unsignedInt">0</MaxEnvelopes>
              <CurrentTime xsi:type="xsd:string">2015-01-21T14:46:07</CurrentTime>
              <RetryCount xsi:type="xsd:unsignedInt">0</RetryCount>
              <ParameterList SOAP-ENC:arrayType="cwmp:ParameterValueStruct[11]">
                  <ParameterValueStruct>
                      <Name xsi:type="string">InternetGatewayDevice.DeviceInfo.ProvisioningCode</Name>
                      <Value xsi:type="string">D915W</Value>
                  </ParameterValueStruct>
                  <ParameterValueStruct>
                      <Name xsi:type="string">InternetGatewayDevice.ManagementServer.URL</Name>
                      <Value xsi:type="string" />
                  </ParameterValueStruct>
                  <ParameterValueStruct>
                      <Name xsi:type="string">InternetGatewayDevice.ManagementServer.Password</Name>
                      <Value xsi:type="string" />
                  </ParameterValueStruct>
                  <ParameterValueStruct>
                      <Name xsi:type="string">InternetGatewayDevice.ManagementServer.ConnectionRequestURL</Name>
                      <Value xsi:type="string">http://192.168.16.68:5400</Value>
                  </ParameterValueStruct>
                  <ParameterValueStruct>
                      <Name xsi:type="string">InternetGatewayDevice.DeviceSummary</Name>
                      <Value xsi:type="string" />
                  </ParameterValueStruct>
                  <ParameterValueStruct>
                      <Name xsi:type="string">InternetGatewayDevice.DeviceInfo.SpecVersion</Name>
                      <Value xsi:type="string">V1.0</Value>
                  </ParameterValueStruct>
                  <ParameterValueStruct>
                      <Name xsi:type="string">InternetGatewayDevice.ManagementServer.UDPConnectionRequestAddress</Name>
                      <Value xsi:type="string">200.200.202.68:1036</Value>
                  </ParameterValueStruct>
                  <ParameterValueStruct>
                      <Name xsi:type="string">InternetGatewayDevice.ManagementServer.ConnectionRequestPassword</Name>
                      <Value xsi:type="string" />
                  </ParameterValueStruct>
                  <ParameterValueStruct>
                      <Name xsi:type="string">InternetGatewayDevice.DeviceInfo.SoftwareVersion</Name>
                      <Value xsi:type="string">nvD915wv2.2.0-150120</Value>
                  </ParameterValueStruct>
                  <ParameterValueStruct>
                      <Name xsi:type="string">InternetGatewayDevice.WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.ExternalIPAddress</Name>
                      <Value xsi:type="string">192.168.16.68</Value>
                  </ParameterValueStruct>
                  <ParameterValueStruct>
                      <Name xsi:type="string">InternetGatewayDevice.DeviceInfo.HardwareVersion</Name>
                      <Value xsi:type="string">V1.0</Value>
                  </ParameterValueStruct>
              </ParameterList>
          </urn:dslforum-org:cwmp-1-0:Inform>
      </http://schemas.xmlsoap.org/soap/envelope/:Body>
  </http://schemas.xmlsoap.org/soap/envelope/:Envelope>
*/
