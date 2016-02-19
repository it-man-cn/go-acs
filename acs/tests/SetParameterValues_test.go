package test

import (
	"encoding/json"
	"fmt"
	"go-acs/acs/models/messages"
	"testing"
)

func TestCreateSetParameterValues(t *testing.T) {
	resp := new(messages.SetParameterValues)
	params := make(map[string]messages.ValueStruct)
	param := messages.ValueStruct{messages.XsdString, "abc"}
	params["InternetGatewayDevice.DeviceInfo.Manufacturer"] = param
	resp.Params = params
	fmt.Println(string(resp.CreateXML()))
	jsonstr, _ := json.Marshal(&resp)
	fmt.Println(string(jsonstr))

}

func TestParseSetParameterValues(t *testing.T) {

}
