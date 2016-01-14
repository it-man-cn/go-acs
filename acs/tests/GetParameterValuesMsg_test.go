package test

import (
	"encoding/json"
	"fmt"
	"go-acs/acs/models/messages"
	"testing"
)

/*
func TestRPCCall(t *testing.T) {
	msg := models.ReciveMsg("acs.1586025")
	fmt.Println(msg.Headers)
	fmt.Println(msg.Body)
	m := models.FromMessage(msg)
	if m != nil {
		fmt.Println(m.GetName(), m.GetId())
		request := m.(*messages.GetParameterValues)
		fmt.Println(request.GetName())
		fmt.Println(request.GetName())
		for k, v := range request.ParamNames {
			fmt.Println(k, v)
		}
		resp := messages.NewGetParameterValuesResponse()
		//resp := new(messages.GetParameterValuesResponse)
		params := make(map[string]string)
		params["InternetGatewayDevice.DeviceInfo.Manufacturer"] = "ACS"
		params["InternetGatewayDevice.DeviceInfo.OUI"] = "0011AB"
		resp.Values = params
		props := models.MessageProperties{
			CorrelationId:   msg.CorrelationId,
			ReplyTo:         msg.ReplyTo,
			ContentEncoding: msg.ContentEncoding,
			ContentType:     msg.ContentType,
		}
		reply, err := models.CreateMessage(resp, props)
		//msg.Body = []byte("reponse from golang")
		fmt.Println(err)
		models.SendMsg(reply)
	}

	//obj := reflect.New(reflect.TypeOf(messages.Inform{}))
	//m := obj.Interface().(messages.Message)
	//fmt.Println(m.GetName(), m.GetId())
}
*/
func TestConverter(t *testing.T) {
	//data := `{"fault":false,"id":"ID:intrnl.unset.id.GetParameterValues1439954709189.1134973336","name":"GetParameterValues","noMore":false,"parameterNames":["InternetGatewayDevice.DeviceInfo.Manufacturer","InternetGatewayDevice.DeviceInfo.OUI"]}`
	data := `{"id":"ID:intrnl.unset.id.GetParameterValues1439956067715.344209010","name":"GetParameterValues","noMore":1,"parameterNames":["InternetGatewayDevice.DeviceInfo.Manufacturer","InternetGatewayDevice.DeviceInfo.OUI"]}`
	msg := new(messages.GetParameterValues)
	json.Unmarshal([]byte(data), &msg)
	fmt.Println(msg.Id)
	fmt.Println(msg.Name)
	fmt.Println(msg.ParameterNames)
	fmt.Println(string(msg.CreateXml()))
}
