package test

import (
	"go-acs/acs/models"
	"testing"
)

func TestLog(t *testing.T) {
	models.Log("1458121345", "50")
}
