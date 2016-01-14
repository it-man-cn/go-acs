package log

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

var Logger = logs.NewLogger(10000)

func init() {
	logfile := beego.AppConfig.String("logfile")
	Logger.EnableFuncCallDepth(true)
	loglevel := beego.AppConfig.String("loglevel")
	switch loglevel {
	case "debug":
		Logger.SetLevel(logs.LevelDebug)
	case "info":
		Logger.SetLevel(logs.LevelInfo)
	case "warn":
		Logger.SetLevel(logs.LevelWarn)
	case "error":
		Logger.SetLevel(logs.LevelError)
	default:
		Logger.SetLevel(logs.LevelInfo)
	}
	Logger.SetLevel(logs.LevelInfo)
	Logger.SetLogger("file", `{"filename":"`+logfile+`","daily":true}`)
}
