package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	log "go-acs/acs/log"
	"go-acs/acs/models"
	"go-acs/acs/models/messages"
	"time"
)

var (
	redisClient *redis.Pool
)

func init() {
	redisHost := beego.AppConfig.String("redis.host")
	redisDB, _ := beego.AppConfig.Int("redis.db")
	redisClient = &redis.Pool{
		MaxIdle:     beego.AppConfig.DefaultInt("redis.maxidle", 1),
		MaxActive:   beego.AppConfig.DefaultInt("redis.maxactive", 10),
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisHost)
			if err != nil {
				return nil, err
			}
			c.Do("SELECT", redisDB)
			return c, nil
		},
	}
}

//MainController process tr069
type MainController struct {
	beego.Controller
}

const (
	minMaxEnvelopes int    = 1
	maxEnvelopes    int    = 1
	attrLastInform  string = "ATTR_LASTINFORM"
	attrSN          string = "ATTR_SN"
	attrCorrID      string = "ATTR_CORR_ID"
	attrReplyTo     string = "ATTR_REPLY_TO"
)

//Get tr069 msg
func (c *MainController) Get() {
	c.processRequest()
}

//Post tr069 msg
func (c *MainController) Post() {
	c.processRequest()
}

func (c *MainController) processRequest() {
	requestBody := c.Ctx.Input.RequestBody
	var response []byte
	var countEnvelopes, maxEnvelopes int = 0, 0
	var lastInform *messages.Inform
	inform := c.Ctx.Input.CruSession.Get(attrLastInform)
	if inform != nil {
		lastInform = inform.(*messages.Inform)
	}
	if lastInform != nil {
		maxEnvelopes = lastInform.MaxEnvelopes
	} else {
		maxEnvelopes = minMaxEnvelopes
	}

	if maxEnvelopes < minMaxEnvelopes {
		log.Info("MaxEnvelopes are less then ", minMaxEnvelopes, "(", maxEnvelopes, ").Setting to", minMaxEnvelopes)
		maxEnvelopes = 1
	}

	cpeSentEmptyReq := true
	var sn string
	if len(requestBody) > 0 {
		cpeSentEmptyReq = false
		msg, err := messages.ParseXML(c.Ctx.Input.RequestBody)
		if err == nil {
			reqType := msg.GetName()
			if "Inform" == reqType {
				lastInform = msg.(*messages.Inform)
				c.Ctx.Input.CruSession.Set(attrLastInform, lastInform)
				sn = lastInform.Sn
				c.Ctx.Input.CruSession.Set(attrSN, sn)
				//lastInform (cpe info) store in memcache
				//TODO
				//response
				resp := new(messages.InformResponse)
				resp.ID = lastInform.ID
				resp.MaxEnvelopes = maxEnvelopes
				response = resp.CreateXML()
				countEnvelopes++
				maxEnvelopes = lastInform.MaxEnvelopes
				log.Debug(lastInform.GetID())
				//log inform
				if lastInform.IsEvent(messages.EventPeriodic) {
					bytesReceived := lastInform.GetParam("InternetGatewayDevice.WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.Stats.EthernetBytesReceived")
					if len(bytesReceived) > 0 {
						//go
						log.Debug("inform write sn:" + sn)
						select {
						case models.InformChan <- models.Receive{Sn: sn, BytesReceived: bytesReceived}:
							log.Debug("channel rwrite ok id: %s\n", sn)
						default:
							log.Debug("channel write block")
						}
					}
					//compare config version
					key := "inform_" + sn
					rc := redisClient.Get()
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
							log.Debug("config version chagne,need change sn:%s\n", sn)
							inform := new(messages.ValueChange)
							inform.Sn = sn
							select {
							case models.ChangeChan <- inform:
								log.Debug("value change channel rwrite ok id: %s\n", sn)
							default:
								log.Debug("value change channel write block")
							}
						}
					}
				} else if lastInform.IsEvent(messages.EventValueChange) {
					//update all config
					inform := new(messages.ValueChange)
					inform.Sn = sn
					select {
					case models.ChangeChan <- inform:
						log.Debug("value change channel rwrite ok id: %s\n", sn)
					default:
						log.Debug("value change channel write block")
					}
				}
				stamp := time.Now().Format("2006-01-02 15:04:05")
				values, err := json.Marshal(models.InformMessage{Inform: lastInform, Timestamp: stamp})
				if err != nil {
					log.Error("%s json encode Error %v\n", sn, err)
				} else {
					key := "inform_" + sn
					value := string(values)
					rc := redisClient.Get()
					rc.Do("SETEX", key, 600, value)
					defer rc.Close()
				}
			} else if "TransferComplete" == reqType {
				tc := msg.(*messages.TransferComplete)
				//do something
				resp := new(messages.TransferCompleteResponse)
				resp.ID = tc.ID
				response = resp.CreateXML()
				countEnvelopes++
			} else if "GetRPCMethods" == reqType {
				gm := msg.(*messages.GetRPCMethods)
				resp := new(messages.GetRPCMethodsResponse)
				resp.ID = gm.ID
				response = resp.CreateXML()
				countEnvelopes++
			} else {
				c.sendRPCResponse(msg)
				//dosomethin ,update log status
				//has message any more
				inform := c.Ctx.Input.CruSession.Get(attrLastInform)
				if inform != nil {
					lastInform = inform.(*messages.Inform)
					msg := c.receiveMsg(lastInform.Sn, 1000)
					if msg != nil {
						response = msg.CreateXML()
						log.Debug("response:", string(response))
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
			log.Debug("sn:%s,Host:%s,IP:%s", sn, c.Ctx.Input.Host(), c.Ctx.Input.IP())
			msg := c.receiveMsg(sn, 1000)
			if msg == nil && lastInform.IsEvent(messages.EventConnectionRequest) {
				msg = c.receiveMsg(sn, 2000)
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
			response = msg.CreateXML()
			log.Debug("response:", string(response))
			//need CPE do something
			countEnvelopes++
		}
	}
	//c.Ctx.Output.Header("ContentType", "text/xml")
	//c.Ctx.Output.ContentType("text/xml")
	c.Ctx.Output.Header("Content-Type", "application/xml; charset=utf-8")
	//c.Ctx.Output.XML(data, hasIndent)
	c.Ctx.Output.Body(response)
}

func (c *MainController) receiveMsg(sn string, timeout int) messages.Message {
	msg := models.ReciveMsg("acs." + sn)
	log.Debug("recv body:", string(msg.Body))
	m := models.FromMessage(msg)
	if m != nil {
		//store info to session
		log.Debug("msgId:" + m.GetID())
		log.Debug("corrId:" + msg.CorrelationID)
		log.Debug("replyTo:" + msg.ReplyTo)
		c.Ctx.Input.CruSession.Set(attrCorrID, msg.CorrelationID)
		c.Ctx.Input.CruSession.Set(attrReplyTo, msg.ReplyTo)
	}
	return m
}

func (c *MainController) sendRPCResponse(msg messages.Message) error {
	corrID := c.Ctx.Input.CruSession.Get(attrCorrID)
	replyTo := c.Ctx.Input.CruSession.Get(attrReplyTo)
	if replyTo != nil && corrID != nil {
		props := models.MessageProperties{
			CorrelationID:   corrID.(string),
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
