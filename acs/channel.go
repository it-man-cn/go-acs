package main

import (
	"sync"
	"time"
)

// Channel used by message pusher send msg to write goroutine.
type Channel struct {
	signal   chan int
	CliProto Message
	SvrProto Ring
	cLock    sync.Mutex
}

//NewChannel init channel
func NewChannel(proto int) *Channel {
	c := new(Channel)
	c.signal = make(chan int, SignalNum)
	c.SvrProto.Init(proto)
	return c
}

// Reply server reply message.
func (c *Channel) Reply(p *Message) (err error) {
	var proto *Message
	c.cLock.Lock()
	if proto, err = c.SvrProto.Set(); err == nil { //get server-end writer buffer
		*proto = *p
		c.SvrProto.SetAdv()
	}
	c.cLock.Unlock()
	if err == nil {
		c.Signal()
	}
	return
}

// Push server push message.
func (c *Channel) Push(prop MessageProperties, body []byte) (err error) {
	var proto *Message
	c.cLock.Lock()
	if proto, err = c.SvrProto.Set(); err == nil { //get svr writer buffer
		//proto.Ver = ver
		//proto.Operation = operation
		//proto.Body = body
		//proto = &msg
		proto.MessageProperties = prop
		proto.Body = body
		c.SvrProto.SetAdv()
	}
	c.cLock.Unlock()
	if err == nil {
		c.Signal()
	}
	return
}

// Ready check the channel ready or close?
func (c *Channel) Ready() bool {
	return (<-c.signal) == ProtoReady
}

// ReadyWithTimeout check the channel ready or close?
func (c *Channel) ReadyWithTimeout(timeout time.Duration) bool {
	var s int
	select {
	case s = <-c.signal:
		return s == ProtoReady
	case <-time.After(timeout):
		return false
	}
}

// Signal send signal to the channel, protocol ready. 数据准备就绪
func (c *Channel) Signal() {
	// just ignore duplication signal
	select {
	case c.signal <- ProtoReady:
	default:
	}
}

// Close close the channel.
func (c *Channel) Close() {
	select { //等待其它操作结束
	case c.signal <- ProtoFinish:
	default:
	}
}
