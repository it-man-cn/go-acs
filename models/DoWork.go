package models

import (
	"github.com/astaxie/beego"
	. "go-acs/log"
	"go-acs/models/messages"
	"os"
	"time"
)

var INFORM_CHANNEL chan Receive
var CHANGE_CHANNEL chan messages.Message
var queuename = beego.AppConfig.String("queuename")
var logdir = beego.AppConfig.String("logdir")

type Receive struct {
	Sn            string
	BytesReceived string
}

func init() {
	var num, length int
	var err error
	length, err = beego.AppConfig.Int("inform_chan_len")
	if err != nil {
		length = 10
	}
	INFORM_CHANNEL = make(chan Receive, length)
	num, err = beego.AppConfig.Int("inform_work_num")
	if err != nil {
		num = 10
	}
	for i := 0; i < num; i++ {
		go func() {
			for {
				select {
				case receive := <-INFORM_CHANNEL:
					Logger.Info("channel read ok id: %s\n", receive.Sn)
					Log(receive)
				}
			}
		}()
	}

	length, err = beego.AppConfig.Int("notice_chan_len")
	if err != nil {
		length = 10
	}
	CHANGE_CHANNEL = make(chan messages.Message)
	num, err = beego.AppConfig.Int("notice_work_num")
	if err != nil {
		num = 10
	}
	for i := 0; i < num; i++ {
		go func() {
			for {
				select {
				case receive := <-CHANGE_CHANNEL:
					Logger.Info("notice channel read ok id: %s\n", receive)
					ValueChange(receive)
				}
			}
		}()
	}
}

func ValueChange(inform messages.Message) {
	props := MessageProperties{
		CorrelationId:   "",
		ReplyTo:         "",
		ContentEncoding: "UTF-8",
		ContentType:     "application/json",
	}
	msg, err := CreateMessage(inform, props)
	if err != nil {
		Logger.Info("createMessae error:%s", err)
	}
	SendMsg(msg, queuename)
}

func checkDirIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func Log(recv Receive) {
	now := time.Now().Format("20060102150405")
	today := now[0:8]
	stamp := now[8:]
	path := logdir + today + "/"
	logfile := recv.Sn + ".log"
	if !checkDirIsExist(path) {
		err := os.Mkdir(path, 0766)
		if err != nil {
			Logger.Info("mkdir %s error %s", path, err)
			return
		}
	}
	f, err := os.OpenFile(path+logfile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		Logger.Info("open file %s error %s", path+logfile, err)
		return
	}
	defer f.Close()
	_, e := f.WriteString(stamp + ":" + recv.BytesReceived + ",")
	if e != nil {
		Logger.Info("wirte file error %s", e)
	}
}
