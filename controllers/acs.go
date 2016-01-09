package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	. "go-acs/log"
	"go-acs/models"
	"go-acs/models/messages"
	"time"
)

var (
	RedisClient *redis.Pool
	REDIS_HOST  string
	REDIS_DB    int
)

func init() {
	REDIS_HOST = beego.AppConfig.String("redis.host")
	REDIS_DB, _ = beego.AppConfig.Int("redis.db")
	RedisClient = &redis.Pool{
		MaxIdle:     beego.AppConfig.DefaultInt("redis.maxidle", 1),
		MaxActive:   beego.AppConfig.DefaultInt("redis.maxactive", 10),
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", REDIS_HOST)
			if err != nil {
				return nil, err
			}
			c.Do("SELECT", REDIS_DB)
			return c, nil
		},
	}
}

type MainController struct {
	beego.Controller
}

const (
	MIN_MAX_ENVELOPES int    = 1
	MY_MAX_ENVELOPES  int    = 1
	ATTR_LASTINFORM   string = "ATTR_LASTINFORM"
	ATTR_SN           string = "ATTR_SN"
	ATTR_CORR_ID      string = "ATTR_CORR_ID"
	ATTR_REPLY_TO     string = "ATTR_REPLY_TO"
)

func (c *MainController) Get() {
	c.processRequest()

}

func (c *MainController) Post() {
	c.processRequest()
}

func (c *MainController) processRequest() {
	requestBody := c.Ctx.Input.RequestBody
	var response []byte
	var countEnvelopes, maxEnvelopes int = 0, 0
	var lastInform *messages.Inform
	inform := c.Ctx.Input.CruSession.Get(ATTR_LASTINFORM)
	if inform != nil {
		lastInform = inform.(*messages.Inform)
	}
	if lastInform != nil {
		maxEnvelopes = lastInform.MaxEnvelopes
	} else {
		maxEnvelopes = MIN_MAX_ENVELOPES
	}

	if maxEnvelopes < MIN_MAX_ENVELOPES {
		Logger.Info("MaxEnvelopes are less then ", MIN_MAX_ENVELOPES, "(", maxEnvelopes, ").Setting to", MIN_MAX_ENVELOPES)
		maxEnvelopes = 1
	}

	var cpeSentEmptyReq bool = true
	var sn string
	if len(requestBody) > 0 {
		cpeSentEmptyReq = false
		msg, err := messages.ParseXml(string(c.Ctx.Input.RequestBody))
		if err == nil {
			reqType := msg.GetName()
			if "Inform" == reqType {
				lastInform = msg.(*messages.Inform)
				c.Ctx.Input.CruSession.Set(ATTR_LASTINFORM, lastInform)
				sn = lastInform.Sn
				c.Ctx.Input.CruSession.Set(ATTR_SN, sn)
				//lastInform (cpe info) store in memcache
				//TODO
				//response
				resp := new(messages.InformResponse)
				resp.Id = lastInform.Id
				resp.MaxEnvelopes = MY_MAX_ENVELOPES
				response = resp.CreateXml()
				countEnvelopes++
				maxEnvelopes = lastInform.MaxEnvelopes
				Logger.Debug(lastInform.GetId())
				//log inform
				if lastInform.IsEvent(messages.EVENT_PERIODIC) {
					bytesReceived := lastInform.GetParam("InternetGatewayDevice.WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.Stats.EthernetBytesReceived")
					if len(bytesReceived) > 0 {
						//go
						Logger.Debug("inform write sn:" + sn)
						select {
						case models.INFORM_CHANNEL <- models.Receive{sn, bytesReceived}:
							Logger.Debug("channel rwrite ok id: %s\n", sn)
						default:
							Logger.Debug("channel write block")
						}
					}
					//compare config version
					key := "inform_" + sn
					rc := RedisClient.Get()
					r, err := redis.String(rc.Do("GET", key))
					defer rc.Close()
					if len(r) > 0 && err == nil {
						informMessage := models.InformMessage{}
						err = json.Unmarshal([]byte(r), &informMessage)
						oldInform := informMessage.Inform
						oldVersion := oldInform.GetConfigVersion()
						newVersion := lastInform.GetConfigVersion()
						if oldVersion != newVersion {
							//need update config
							Logger.Debug("config version chagne,need change sn:%s\n", sn)
							inform := new(messages.ValueChange)
							inform.Sn = sn
							select {
							case models.CHANGE_CHANNEL <- inform:
								Logger.Debug("value change channel rwrite ok id: %s\n", sn)
							default:
								Logger.Debug("value change channel write block")
							}
						}
					}
				} else if lastInform.IsEvent(messages.EVENT_VALUE_CHANGE) {
					//update all config
					inform := new(messages.ValueChange)
					inform.Sn = sn
					select {
					case models.CHANGE_CHANNEL <- inform:
						Logger.Debug("value change channel rwrite ok id: %s\n", sn)
					default:
						Logger.Debug("value change channel write block")
					}
				}
				stamp := time.Now().Format("2006-01-02 15:04:05")
				values, err := json.Marshal(models.InformMessage{Inform: lastInform, Timestamp: stamp})
				if err == nil {
					key := "inform_" + sn
					value := string(values)
					rc := RedisClient.Get()
					rc.Do("SETEX", key, 600, value)
					defer rc.Close()
				} else {
					Logger.Error("Error %v", err)
				}
			} else if "TransferComplete" == reqType {
				tc := msg.(*messages.TransferComplete)
				//do something
				resp := new(messages.TransferCompleteResponse)
				resp.Id = tc.Id
				response = resp.CreateXml()
				countEnvelopes++
			} else if "GetRPCMethods" == reqType {
				gm := msg.(*messages.GetRPCMethods)
				resp := new(messages.GetRPCMethodsResponse)
				resp.Id = gm.Id
				response = resp.CreateXml()
				countEnvelopes++
			} else {
				c.SendRPCResponse(msg)
				//dosomethin ,update log status
				//has message any more
				inform := c.Ctx.Input.CruSession.Get(ATTR_LASTINFORM)
				if inform != nil {
					lastInform = inform.(*messages.Inform)
					msg := c.ReceiveMsg(lastInform.Sn, 1000)
					if msg != nil {
						response = msg.CreateXml()
						Logger.Debug("response:", string(response))
					}
				}

			}
		} else {
			response = []byte("error")
		}
	}

	if lastInform != nil && cpeSentEmptyReq {
		//response command (Set,Reboot etc)
		sn = lastInform.Sn
		idle, keepalive := 0, 0
		for countEnvelopes < maxEnvelopes {
			//Get command from MQ(RabbitMQ)
			Logger.Debug("sn:%s,Host:%s,IP:%s", sn, c.Ctx.Input.Host(), c.Ctx.Input.IP())
			msg := c.ReceiveMsg(sn, 1000)
			if msg == nil && lastInform.IsEvent(messages.EVENT_CONNECTION_REQUEST) {
				msg = c.ReceiveMsg(sn, 2000)
			}

			if msg == nil {
				if idle >= keepalive {
					idle++
					break
				} else {
					countEnvelopes++
					continue
				}
			}
			response = msg.CreateXml()
			Logger.Debug("response:", string(response))
			//need CPE do something
			countEnvelopes++
		}
	}
	//c.Ctx.Output.Header("ContentType", "text/xml")
	//c.Ctx.Output.ContentType("text/xml")
	c.Ctx.Output.Header("Content-Type", "application/xml; charset=utf-8")
	//c.Ctx.Output.Xml(data, hasIndent)
	c.Ctx.Output.Body(response)
}

func (c *MainController) ReceiveMsg(sn string, timeout int) messages.Message {
	msg := models.ReciveMsg("acs." + sn)
	Logger.Debug("recv body:", string(msg.Body))
	m := models.FromMessage(msg)
	if m != nil {
		//store info to session
		Logger.Debug("msgId:" + m.GetId())
		Logger.Debug("corrId:" + msg.CorrelationId)
		Logger.Debug("replyTo:" + msg.ReplyTo)
		c.Ctx.Input.CruSession.Set(ATTR_CORR_ID, msg.CorrelationId)
		c.Ctx.Input.CruSession.Set(ATTR_REPLY_TO, msg.ReplyTo)
	}
	return m
}

func (c *MainController) SendRPCResponse(msg messages.Message) error {
	corrId := c.Ctx.Input.CruSession.Get(ATTR_CORR_ID)
	replyTo := c.Ctx.Input.CruSession.Get(ATTR_REPLY_TO)
	if replyTo != nil && corrId != nil {
		props := models.MessageProperties{
			CorrelationId:   corrId.(string),
			ReplyTo:         replyTo.(string),
			ContentEncoding: "UTF-8",
			ContentType:     "application/json",
		}
		reply, err := models.CreateMessage(msg, props)
		if err != nil {
			return err
		}
		models.Reply(reply, replyTo.(string))
	}
	return nil
}
