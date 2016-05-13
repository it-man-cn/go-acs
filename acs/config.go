package main

import (
	//"flag"
	"github.com/Terry-Mao/goconf"
	"runtime"
	"time"
)

var (
	gconf *goconf.Config
	//Conf global config
	Conf     *Config
	confFile string
)

func init() {
	//flag.StringVar(&confFile, "c", "./conf/acs.conf", " set config file path")
	confFile = "./conf/acs.conf"
}

//Config config
type Config struct {
	// base section
	PidFile   string   `goconf:"base:pidfile"`
	Dir       string   `goconf:"base:dir"`
	Log       string   `goconf:"base:log"`
	LogDir    string   `goconf:"base:logdir"`
	MaxProc   int      `goconf:"base:maxproc"`
	PprofBind []string `goconf:"base:pprof.bind:,"`
	Debug     bool     `goconf:"base:debug"`
	HTTPBind  []string `goconf:"base:http.bind:,"`

	//channel
	RingSize    int `goconf:"channel:RingSize"`
	ChannelSize int `goconf:"channel:channel"`

	//mq
	AMQPUrl       string        `goconf:"mq:amqpurl"`
	ACSQueue      string        `goconf:"mq:acsqueue"`
	RetryInterval time.Duration `goconf:"mq:retryinterval:time"`

	//redis
	RedisDB        int    `goconf:"redis:redis.db"`
	RedisHost      string `goconf:"redis:redis.host"`
	RedisMaxIdle   int    `goconf:"redis:redis.maxidle"`
	RedisMaxActive int    `goconf:"redis:redis.maxactive"`

	//work
	InformWorkNum  int `goconf:"work:informWorkNum"`
	InformChanSize int `goconf:"work:informChanSize"`
	NoticeWorkNum  int `goconf:"work:noticeWorkNum"`
	NoticeChanSize int `goconf:"work:noticeChanSize"`
	ConfigWorkNum  int `goconf:"work:configWorkNum"`
	ConfigChanSize int `goconf:"work:configChanSize"`
}

func newConfig() *Config {
	return &Config{
		// base section
		PidFile:   "/tmp/acs.pid",
		Dir:       "./",
		Log:       "./conf",
		MaxProc:   runtime.NumCPU(),
		PprofBind: []string{"localhost:6971"},
		Debug:     true,
		HTTPBind:  []string{"localhost:19001"},
		LogDir:    "./",

		// Channel
		RingSize:    3,
		ChannelSize: 1000,

		//mq
		AMQPUrl:       "amqp://guest:guest@localhost",
		ACSQueue:      "acs.msg",
		RetryInterval: 30 * time.Second,

		//redis
		RedisDB:        0,
		RedisHost:      "127.0.0.1:6379",
		RedisMaxActive: 20,
		RedisMaxIdle:   5,

		//work
		InformWorkNum:  10,
		InformChanSize: 10,
		NoticeWorkNum:  2,
		NoticeChanSize: 10,
		ConfigWorkNum:  2,
		ConfigChanSize: 10,
	}
}

// InitConfig init the global config.
func InitConfig() (err error) {
	Conf = newConfig()
	gconf = goconf.New()
	if err = gconf.Parse(confFile); err != nil {
		return err
	}

	if err := gconf.Unmarshal(Conf); err != nil {
		return err
	}
	return nil
}
