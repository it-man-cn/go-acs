package test

import (
	"fmt"
	"go-acs/acs/models/messages"
	"testing"
)

func TestCreateGetRPCMethods(t *testing.T) {
	resp := new(messages.GetRPCMethods)
	fmt.Println(string(resp.CreateXml()))
}
