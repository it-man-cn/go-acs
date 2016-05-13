package main

import (
	"flag"
	log "github.com/it-man-cn/log4go"
	"go-acs/libs/perf"
	"runtime"
)

var (
	//DefaultBucket bucket
	DefaultBucket *Bucket
	Debug         bool
	//DefaultMsgQueue *MsgQueue
)

func main() {
	flag.Parse()
	if err := InitConfig(); err != nil {
		panic(err)
	}
	Debug = Conf.Debug
	runtime.GOMAXPROCS(Conf.MaxProc)
	log.LoadConfiguration(Conf.Log)
	defer log.Close()
	//Pprof listen
	perf.Init(Conf.PprofBind)

	DefaultBucket = NewBucket(BucketOptions{
		ChannelSize: Conf.ChannelSize,
	})
	//DefaultMsgQueue = NewMsgQueue()
	NewMsgQueue()
	if err := initHTTP(Conf.HTTPBind); err != nil {
		panic(err)
	}
	// block until a signal is received.
	//InitSignal()
	forever := make(chan bool)
	<-forever

}
