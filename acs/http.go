package main

import (
	log "github.com/it-man-cn/log4go"
	"go-acs/acs/models/messages"
	"go-acs/acs/session"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
)

var (
	globalSessions *session.Manager
)

const (
	minMaxEnvelopes int    = 1
	maxEnvelopes    int    = 1
	attrLastInform  string = "ATTR_LASTINFORM"
	attrSN          string = "ATTR_SN"
	attrCorrID      string = "ATTR_CORR_ID"
	attrReplyTo     string = "ATTR_REPLY_TO"
)

func initHTTP(addrs []string) (err error) {
	var (
		bind         string
		listener     *net.TCPListener
		addr         *net.TCPAddr
		server       *http.Server
		httpServeMux = http.NewServeMux()
	)
	globalSessions = session.NewManager()

	httpServeMux.HandleFunc("/ACS/tr069", tr069)

	for _, bind = range addrs {
		if addr, err = net.ResolveTCPAddr("tcp4", bind); err != nil {
			log.Error("net.ResolveTCPAddr(\"tcp4\", \"", bind, "\") error:", err)
			return
		}
		if listener, err = net.ListenTCP("tcp4", addr); err != nil {
			log.Error("net.ListenTCP(\"tcp4\", \"", bind, "\") error:", err)
			return
		}
		server = &http.Server{Handler: httpServeMux}
		log.Info("start http listen: \"%s\"", bind)
		go func() {
			if err = server.Serve(listener); err != nil {
				log.Error("server.Serve(\"%s\") error(%v)", bind, err)
				panic(err)
			}
		}()
	}
	return
}

func tr069(w http.ResponseWriter, r *http.Request) {
	var requestBody []byte
	var err error
	if r.Method == "POST" {
		// receive posted data
		requestBody, err = ioutil.ReadAll(r.Body)
		if err != nil {
			log.Warn("tr069 read body error(%v)", err)
		}
	}

	var response []byte
	cpeSentEmptyReq := true
	var sn string
	if len(requestBody) > 0 {
		cpeSentEmptyReq = false
		msg, err := messages.ParseXML(requestBody)
		if err == nil {
			switch msg.GetName() {
			case "Inform":
				curInform := msg.(*messages.Inform)
				sn = curInform.Sn
				globalSessions.Put(w, r, sn)
				//response
				resp := new(messages.InformResponse)
				resp.MaxEnvelopes = maxEnvelopes
				response = resp.CreateXML()
				if curInform.IsEvent(messages.EventPeriodic) {
					//TODO
					//sync device info
				} else if curInform.IsEvent(messages.EventValueChange) {
					//TODO
					//sync config
				}
			case "TransferComplete":
				tc := msg.(*messages.TransferComplete)
				resp := new(messages.TransferCompleteResponse)
				resp.ID = tc.ID
				//TODO
				//rpc reply
				response = resp.CreateXML()
			case "GetRPCMethods":
				gm := msg.(*messages.GetRPCMethods)
				resp := new(messages.GetRPCMethodsResponse)
				resp.ID = gm.ID
				response = resp.CreateXML()
			default:
				//rpc reply
				//read replyto and send reponse msg to mq
				//TODO
				//根据sessionid读取replytTo,corrID，session需要保证replyTo的一致性，不是每条消息都有replyTo
				session := globalSessions.Get(w, r)
				if session != nil {
					//sn := session.SN
					if len(session.ReplyTo) > 0 {
						//send rpc reply
						sendRPCResponse(msg, session.ReplyTo, session.CorrelationID)
					}
					//has msg anymore
					var channel *Channel
					if channel = DefaultBucket.Channel(session.SN); channel != nil {
						//read msg
						var msg *Message
						var err error
						if msg, err = channel.SvrProto.Get(); err != nil {
							// must be empty error
							err = nil
							session.ReplyTo = ""
							session.CorrelationID = ""
						} else {
							// just forward the message
							if Debug {
								log.Debug("channel msg:", string(msg.Body))
							}
							m := FromMessage(*msg)
							session.ReplyTo = msg.ReplyTo
							session.CorrelationID = msg.CorrelationID
							response = m.CreateXML()
							channel.SvrProto.GetAdv() //读计算器累加
						}
					} else {
						session.ReplyTo = ""
						session.CorrelationID = ""
					}
				}

			}
		}
	}

	if cpeSentEmptyReq {
		session := globalSessions.Get(w, r)
		if session != nil {
			//根据sn取消息
			var channel *Channel
			if channel = DefaultBucket.Channel(session.SN); channel != nil {
				//read msg
				var msg *Message
				var err error
				if msg, err = channel.SvrProto.Get(); err != nil {
					// must be empty error
					err = nil
					session.ReplyTo = ""
					session.CorrelationID = ""
				} else {
					// just forward the message
					if Debug {
						log.Debug("channel msg:", string(msg.Body))
					}
					m := FromMessage(*msg)
					session.ReplyTo = msg.ReplyTo
					session.CorrelationID = msg.CorrelationID
					response = m.CreateXML()
					channel.SvrProto.GetAdv() //读计算器累加
				}
			} else {
				session.ReplyTo = ""
				session.CorrelationID = ""
			}
		}
	}
	//write response
	w.Header().Add("Content-Type", "application/xml; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(response)))
	w.Write(response)
}

func sendRPCResponse(msg messages.Message, replyTo, corrID string) error {
	props := MessageProperties{
		CorrelationID:   corrID,
		ReplyTo:         replyTo,
		ContentEncoding: "UTF-8",
		ContentType:     "application/json",
	}
	reply, err := CreateMessage(msg, props)
	if err != nil {
		return err
	}
	Reply(reply, replyTo)
	return nil
}
