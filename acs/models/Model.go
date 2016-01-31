package models

import (
	"go-acs/acs/models/messages"
)

//InformMessage used to store device inform
type InformMessage struct {
	Inform    *messages.Inform
	Timestamp string
}
