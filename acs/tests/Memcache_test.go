package test

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/bradfitz/gomemcache/memcache"
	"go-acs/acs/models"
	"go-acs/acs/models/messages"
	"testing"
	"time"
)

func TestMemcacheSet(t *testing.T) {
	stamp := time.Now().Format("2006-01-02 15:04:05")
	sn := "1456789"
	lastInform := &messages.Inform{Id: "abc",
		Manufacturer: "ACS", OUI: "0011ab",
		ProductClass: "it-man",
		Sn:           sn,
		MaxEnvelopes: 1,
		CurrentTime:  "2015-02-12T13:40:07",
		RetryCount:   1}
	values, err := json.Marshal(models.InformMessage{Inform: lastInform, Timestamp: stamp})
	if err == nil {
		mc := memcache.New(beego.AppConfig.String("memcache"))
		fmt.Println(string(values))
		key := "inform_" + sn
		fmt.Println(key)
		item := &memcache.Item{Key: key, Value: values, Flags: 32, Expiration: 600}
		mc.Set(item)
	} else {
		fmt.Println(err)
	}
}

func TestMemcacheGet(t *testing.T) {
	sn := "1456789"
	key := "inform_" + sn
	mc := memcache.New(beego.AppConfig.String("memcache"))
	fmt.Println(key)
	item, err := mc.Get(key)
	if err == nil {
		fmt.Println(string(item.Value))
	}

}
