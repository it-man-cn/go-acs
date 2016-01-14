package main

import (
	"fmt"
	"github.com/astaxie/beego/config"
	"strconv"
)

var (
	Conf  *Config
	Debug bool
)

type Config struct {
	PrimaryPort int
	PrimaryIp   string
	Log         string
	Bind        string
	MaxProc     int
}

func init() {
	Conf = &Config{}
	iniconf, err := config.NewConfig("ini", "./stund.conf")
	if err != nil {
		panic(err)
	}
	Conf.PrimaryPort, _ = iniconf.Int("primary.port")
	Conf.PrimaryIp = iniconf.String("primary.ip")
	Conf.Log = iniconf.String("log")
	fmt.Println("config log " + Conf.Log)
	Debug, _ = iniconf.Bool("debug")
	Conf.Bind = Conf.PrimaryIp + ":" + strconv.Itoa(Conf.PrimaryPort)
	Conf.MaxProc, _ = iniconf.Int("maxproc")
}
