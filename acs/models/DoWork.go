package models

import (
	"github.com/astaxie/beego"
	log "go-acs/acs/log"
	"go-acs/acs/models/messages"
	"os"
	"time"
)

//InformChan is a inform msg queue,used for device log
var InformChan chan Receive

//ChangeChan is a change msg queue
var ChangeChan chan messages.Message
var queuename = beego.AppConfig.String("queuename")
var logdir = beego.AppConfig.String("logdir")

//Receive device log
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
	InformChan = make(chan Receive, length)
	num, err = beego.AppConfig.Int("inform_work_num")
	if err != nil {
		num = 10
	}
	for i := 0; i < num; i++ {
		go func() {
			for {
				select {
				case receive := <-InformChan:
					log.Info("channel read ok id: %s\n", receive.Sn)
					logFile(receive)
				}
			}
		}()
	}

	length, err = beego.AppConfig.Int("notice_chan_len")
	if err != nil {
		length = 10
	}
	ChangeChan = make(chan messages.Message, length)
	num, err = beego.AppConfig.Int("notice_work_num")
	if err != nil {
		num = 10
	}
	for i := 0; i < num; i++ {
		go func() {
			for {
				select {
				case receive := <-ChangeChan:
					log.Info("notice channel read ok id: %s\n", receive)
					valueChange(receive)
				}
			}
		}()
	}
}

func valueChange(inform messages.Message) {
	props := MessageProperties{
		CorrelationID:   "",
		ReplyTo:         "",
		ContentEncoding: "UTF-8",
		ContentType:     "application/json",
	}
	msg, err := CreateMessage(inform, props)
	if err != nil {
		log.Info("createMessae error:%s", err)
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

func logFile(recv Receive) {
	now := time.Now().Format("20060102150405")
	today := now[0:8]
	stamp := now[8:]
	path := logdir + today + "/"
	logfile := recv.Sn + ".log"
	if !checkDirIsExist(path) {
		err := os.Mkdir(path, 0766)
		if err != nil {
			log.Info("mkdir %s error %s", path, err)
			return
		}
	}
	f, err := os.OpenFile(path+logfile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Info("open file %s error %s", path+logfile, err)
		return
	}
	defer f.Close()
	_, e := f.WriteString(stamp + ":" + recv.BytesReceived + ",")
	if e != nil {
		log.Info("wirte file error %s", e)
	}
}
