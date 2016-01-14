package test

import (
	"encoding/json"
	"fmt"
	"go-acs/acs/models/messages"
	"testing"
)

func TestConvertJson2Struct(t *testing.T) {
	jsonString := `{"CurrentTime":"2015-01-21T14:46:07","Manufacturer":"ACS","MaxEnvelopes":1,"ProductClass":"it-man","RetryCount":0,"events":[{"key":"2 PERIODIC","value":""}],"fault":false,"id":"ID:intrnl.unset.id.Inform958064891.1658176487","name":"Inform","noMore":false,"oui":"0011AB","params":{"InternetGatewayDevice.DeviceInfo.HardwareVersion":"V1.0","InternetGatewayDevice.DeviceInfo.ProvisioningCode":"it-man","InternetGatewayDevice.DeviceInfo.SoftwareVersion":"it-man.bin-150120","InternetGatewayDevice.DeviceInfo.SpecVersion":"V1.0","InternetGatewayDevice.DeviceSummary":"","InternetGatewayDevice.ManagementServer.ConnectionRequestPassword":"","InternetGatewayDevice.ManagementServer.ConnectionRequestURL":"http://192.168.16.68:5400","InternetGatewayDevice.ManagementServer.Password":"","InternetGatewayDevice.ManagementServer.UDPConnectionRequestAddress":"200.200.202.68:1036","InternetGatewayDevice.ManagementServer.URL":"","InternetGatewayDevice.WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.ExternalIPAddress":"192.168.16.68"},"sn":"1456789"}`
	inform := messages.Inform{}
	json.Unmarshal([]byte(jsonString), &inform)
	fmt.Println(inform.Manufacturer)
	fmt.Println(inform.OUI)
	fmt.Println(inform.Sn)
	fmt.Println(inform.ProductClass)
	fmt.Println(inform.Id)
	fmt.Println(inform.Name)
	fmt.Println("curTime", inform.CurrentTime)

	for k, v := range inform.Events {
		fmt.Println(k, v)
	}

	for k, v := range inform.Params {
		fmt.Println(k, v)
	}
}

func TestConvertStruct2Json(t *testing.T) {
	inform := messages.Inform{Id: "abc",
		Manufacturer: "ACS", OUI: "0011ab",
		ProductClass: "it-man",
		Sn:           "1456789",
		MaxEnvelopes: 1,
		CurrentTime:  "2015-02-12T13:40:07",
		RetryCount:   1}
	events := make(map[string]string)
	events["6 CONNECTION REQUEST"] = ""
	inform.Events = events
	params := make(map[string]string)
	params["InternetGatewayDevice.DeviceInfo.Manufacturer"] = "ACS"
	inform.Params = params
	jsonString, _ := json.Marshal(&inform)
	fmt.Println(string(jsonString))
}
