package test

import (
	"fmt"
	"go-acs/models"
	"go-acs/models/messages"
	//"reflect"

	"testing"
)

func TestValueChange(t *testing.T) {
	sn := "15303174"
	inform := new(messages.ValueChange)
	inform.Sn = sn
	queuename := "WX_VALUE_CHANGE"
	props := models.MessageProperties{
		CorrelationId:   "",
		ReplyTo:         queuename,
		ContentEncoding: "UTF-8",
		ContentType:     "application/json",
		//Expiration:      "1000",
	}
	msg, err := models.CreateMessage(inform, props)
	if err != nil {
		fmt.Println(err)
	}
	models.SendMsg(msg, queuename)
}

/*
func TestClientChange(t *testing.T) {
	sn := "15303174"
	onlineInform := new(messages.OnlineInform)
	onlineInform.Sn = sn
	newMac := "71:08:13:14:33:77"
	newHostName := "chenbindeiPhone"
	onlineInform.Hosts = append(onlineInform.Hosts,
		messages.Host{Mac: newMac, HostName: newHostName})
	queuename := "WX_VALUE_CHANGE"
	props := models.MessageProperties{
		CorrelationId:   "",
		ReplyTo:         "",
		ContentEncoding: "UTF-8",
		ContentType:     "application/json",
	}
	msg, err := models.CreateMessage(onlineInform, props)
	if err != nil {
		fmt.Println(err)
	}
	models.SendMsg(msg, queuename)
}
*/
