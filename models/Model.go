package models

import (
	"go-acs/models/messages"
)

type InformMessage struct {
	Inform    *messages.Inform
	Timestamp string
}
