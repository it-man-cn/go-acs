package test

import (
	"encoding/json"
	"fmt"
	"go-acs/models/messages"
	"testing"
)

func TestCreateSetParameterValues(t *testing.T) {
	resp := new(messages.SetParameterValues)
	params := make(map[string]messages.ValueStruct)
	param := messages.ValueStruct{messages.XSD_STRING, "abc"}
	params["InternetGatewayDevice.DeviceInfo.Manufacturer"] = param
	resp.Params = params
	fmt.Println(string(resp.CreateXml()))
	jsonstr, _ := json.Marshal(&resp)
	fmt.Println(string(jsonstr))

}

func TestParseSetParameterValues(t *testing.T) {

}
