package test

import (
	//"fmt"
	"go-acs/models"
	//"go-acs/models/messages"
	//"reflect"

	"testing"
)

/*
func TestReciveMsg(t *testing.T) {
	msg := models.ReciveMsg("acs.1586025")
	fmt.Println(msg.Headers)
	fmt.Println(msg.Body)
	m := models.FromMessage(msg)
	if m != nil {
		fmt.Println(m.GetName(), m.GetId())
		inform := m.(*messages.Inform)
		fmt.Println(inform.OUI)
		for k, v := range inform.Params {
			fmt.Println(k, v)
		}
		//msg.Body = []byte("reponse from golang")
		//models.SendMsg(msg)
	}

	//obj := reflect.New(reflect.TypeOf(messages.Inform{}))
	//m := obj.Interface().(messages.Message)
	//fmt.Println(m.GetName(), m.GetId())
}
*/
func TestSendMsg(t *testing.T) {
	props := models.MessageProperties{
		CorrelationId:   "",
		ReplyTo:         "",
		ContentEncoding: "UTF-8",
		ContentType:     "text/plain",
	}

	msg := models.Message{MessageProperties: props, Body: []byte("15303174")}
	models.SendMsg(msg, "stun_queue")
}
