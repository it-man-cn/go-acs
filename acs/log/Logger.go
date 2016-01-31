package log

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

//Logger object
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

//Info call beego log info
func Info(format string, v ...interface{}) {
	Logger.Info(format, v)
}

//Debug call beego
func Debug(format string, v ...interface{}) {
	Logger.Debug(format, v)
}

//Error call beego
func Error(format string, v ...interface{}) {
	Logger.Error(format, v)
}
