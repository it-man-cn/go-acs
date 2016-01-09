package models

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"go-acs/models/messages"
	"reflect"
)

const (
	DEFAULT_CLASSID_FIELD_NAME string = "__TypeId__"
)

var idClassMapping map[string]reflect.Type
var classIdMapping map[reflect.Type]string

func init() {
	idClassMapping = map[string]reflect.Type{
		"GetParameterValues":         reflect.TypeOf(messages.GetParameterValues{}),
		"GetParameterValuesResponse": reflect.TypeOf(messages.GetParameterValuesResponse{}),
		"SetParameterValues":         reflect.TypeOf(messages.SetParameterValues{}),
		"SetParameterValuesResponse": reflect.TypeOf(messages.SetParameterValuesResponse{}),
		"Reboot":                     reflect.TypeOf(messages.Reboot{}),
		"RebootResponse":             reflect.TypeOf(messages.RebootResponse{}),
		"Download":                   reflect.TypeOf(messages.Download{}),
		"DownloadResponse":           reflect.TypeOf(messages.DownloadResponse{}),
		"GetRPCMethods":              reflect.TypeOf(messages.GetRPCMethods{}),
		"GetRPCMethodsResponse":      reflect.TypeOf(messages.GetRPCMethodsResponse{}),
		"Inform":                     reflect.TypeOf(messages.Inform{}),
		"InformResponse":             reflect.TypeOf(messages.InformResponse{}),
		"OnlineInform":               reflect.TypeOf(messages.OnlineInform{}),
		"ValueChange":                reflect.TypeOf(messages.ValueChange{}),
	}

	classIdMapping = map[reflect.Type]string{
		reflect.TypeOf(&messages.GetParameterValues{}):         "GetParameterValues",
		reflect.TypeOf(&messages.GetParameterValuesResponse{}): "GetParameterValuesResponse",
		reflect.TypeOf(&messages.SetParameterValues{}):         "SetParameterValues",
		reflect.TypeOf(&messages.SetParameterValuesResponse{}): "SetParameterValuesResponse",
		reflect.TypeOf(&messages.Reboot{}):                     "Reboot",
		reflect.TypeOf(&messages.RebootResponse{}):             "RebootResponse",
		reflect.TypeOf(&messages.Download{}):                   "Download",
		reflect.TypeOf(&messages.DownloadResponse{}):           "DownloadResponse",
		reflect.TypeOf(&messages.GetRPCMethods{}):              "GetRPCMethods",
		reflect.TypeOf(&messages.GetRPCMethodsResponse{}):      "GetRPCMethodsResponse",
		reflect.TypeOf(&messages.Inform{}):                     "Inform",
		reflect.TypeOf(&messages.InformResponse{}):             "InformResponse",
		reflect.TypeOf(&messages.OnlineInform{}):               "OnlineInform",
		reflect.TypeOf(&messages.ValueChange{}):                "ValueChange",
	}
}

func FromMessage(m Message) messages.Message {
	var msg messages.Message
	value := m.Headers[DEFAULT_CLASSID_FIELD_NAME]
	if value != nil {
		classid := fmt.Sprintf("%s", value)
		clazz := idClassMapping[classid]
		if clazz != nil {
			class := reflect.New(clazz)
			msg = class.Interface().(messages.Message)
			json.Unmarshal(m.Body, &msg)
		}
	}
	return msg
}

func CreateMessage(m messages.Message, properties MessageProperties) (Message, error) {
	typ := reflect.TypeOf(m)
	clazz := classIdMapping[typ]
	body, err := json.Marshal(m)
	if err != nil {
		return Message{}, err
	}
	return Message{
		MessageProperties: MessageProperties{
			Headers:         amqp.Table{DEFAULT_CLASSID_FIELD_NAME: clazz},
			CorrelationId:   properties.CorrelationId,
			ReplyTo:         properties.ReplyTo,
			ContentEncoding: properties.ContentEncoding,
			ContentType:     properties.ContentType,
		},
		Body: body}, nil
}
