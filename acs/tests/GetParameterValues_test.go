package test

import (
	"fmt"
	"go-acs/acs/models/messages"
	"testing"
)

func TestCreateGetParameterValues(t *testing.T) {
	resp := new(messages.GetParameterValues)
    var names []string
	names = append(names, "InternetGatewayDevice.DeviceInfo.Manufacturer", "InternetGatewayDevice.DeviceInfo.ProvisioningCode")
	resp.ParamNames = names
	fmt.Println(string(resp.CreateXml()))
}

func TestParseGetParameterValues(t *testing.T) {

}
