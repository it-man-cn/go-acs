package main

import (
	"fmt"
	"github.com/astaxie/beego/session"
	log "github.com/it-man-cn/go-acs/acs/log"
	"github.com/it-man-cn/go-acs/acs/models/messages"
	"io/ioutil"
	"net"
	"net/http"
)

var globalSessions *session.Manager

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
	globalSessions, _ = session.NewManager("memory", `{"cookieName":"gosessionid","gclifetime":3600}`)
	go globalSessions.GC()
	httpServeMux.HandleFunc("/tr069", tr069)
	for _, bind = range addrs {
		if addr, err = net.ResolveTCPAddr("tcp4", bind); err != nil {
			log.Error("net.ResolveTCPAddr(\"tcp4\", \"%s\") error(%v)", bind, err)
			return
		}
		if listener, err = net.ListenTCP("tcp4", addr); err != nil {
			log.Error("net.ListenTCP(\"tcp4\", \"%s\") error(%v)", bind, err)
			return
		}
		server = &http.Server{Handler: httpServeMux}
		log.Debug("start http listen: \"%s\"", bind)
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
	fmt.Println("tr069")
	sess, err := globalSessions.SessionStart(w, r)
	if err != nil {
		fmt.Println(err)
	}
	//defer sess.SessionRelease(w)
	fmt.Println("sessionId:", sess.SessionID)
	var requestBody []byte
	if r.Method == "POST" {
		// receive posted data
		requestBody, err = ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(nil)
		}
		//fmt.Println(string(requestBody))
	}

	var response []byte
	//var lastInform *messages.Inform

	cpeSentEmptyReq := true
	var sn string
	isInform := false
	if len(requestBody) > 0 {
		cpeSentEmptyReq = false
		msg, err := messages.ParseXML(requestBody)
		if err == nil {
			switch msg.GetName() {
			case "Inform":
				isInform = true
				curInform := msg.(*messages.Inform)
				sn = curInform.Sn
				//sess.Set(attrLastInform, lastInform)
				sess.Set(attrSN, sn)
				//response
				resp := new(messages.InformResponse)
				resp.MaxEnvelopes = maxEnvelopes
				response = resp.CreateXML()
				//dosomething, update config
				//TODO
			case "TransferComplete":
				tc := msg.(*messages.TransferComplete)
				//do something
				resp := new(messages.TransferCompleteResponse)
				resp.ID = tc.ID
				response = resp.CreateXML()
			case "GetRPCMethods":
				gm := msg.(*messages.GetRPCMethods)
				resp := new(messages.GetRPCMethodsResponse)
				resp.ID = gm.ID
				response = resp.CreateXML()
			default:
				fmt.Println("rpc reply")
				//rpc reply
				//read replyto and send reponse msg to mq
				//TODO

				//has msg anymore,read from mq
				//query msg by sn
				//replyto save in session
				//TODO
			}
		}
	}

	if !isInform {
		val := sess.Get(attrSN)
		if val != nil {
			sn = val.(string)
		}
	}
	//read tr069 job
	fmt.Println("sn:", sn)
	if len(sn) > 0 && cpeSentEmptyReq {
		//query msg by sn
		fmt.Println("query msg")
		//replyto save in session
	}

	//write response
	w.Header().Add("Content-Type", "application/xml; charset=utf-8")
	w.Write(response)
}
