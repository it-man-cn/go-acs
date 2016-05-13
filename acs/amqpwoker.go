package main

import (
	log "github.com/it-man-cn/log4go"
	"github.com/streadway/amqp"
	"sync"
	"time"
)

//MsgQueue sn->message queue
type MsgQueue struct {
	cLock sync.RWMutex
	//queue map[string]models.Message
}

//NewMsgQueue init msg queue
func NewMsgQueue() (queue *MsgQueue) {
	queue = new(MsgQueue)
	//queue.queue = make(map[string]models.Message)
	go queue.readMsg()
	return
}

func (q *MsgQueue) readMsg() {
	for {
		q.fetchMsg()
		time.Sleep(Conf.RetryInterval)
	}
}

func (q *MsgQueue) fetchMsg() {
	log.Info("fetch")
	conn, err := amqp.Dial(Conf.AMQPUrl)
	if err != nil {
		log.Error("Failed to connect to RabbitMQ", err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Error("Failed to open a channel", err)
		return
	}
	defer ch.Close()

	qu, err := ch.QueueDeclare(
		Conf.ACSQueue, // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		log.Error("Failed to declare a queue ", err)
		return
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		log.Error("Failed to set QoS", err)
		return
	}

	msgs, err := ch.Consume(
		qu.Name, // queue
		"",      // consumer
		false,   // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		log.Error("Failed to register a consumer", err)
		return
	}
	forever := make(chan bool)
	go func(flag chan bool) {
		log.Error("closing", <-conn.NotifyClose(make(chan *amqp.Error)))
		flag <- true
	}(forever)
	go func() {
		for deliver := range msgs {
			if Debug {
				log.Debug("Received a message: %s", deliver.Body)
				log.Debug("Received a message replyTo: %s", deliver.ReplyTo)
			}
			deliver.Ack(false)
			if deliver.Headers["DeviceID"] != nil {
				sn := deliver.Headers["DeviceID"].(string)
				msg := Message{
					MessageProperties: MessageProperties{
						Headers:         deliver.Headers,
						CorrelationID:   deliver.CorrelationId,
						ReplyTo:         deliver.ReplyTo,
						ContentEncoding: deliver.ContentEncoding,
						ContentType:     deliver.ContentType,
					},
					Body: deliver.Body,
				}
				//q.Push(sn, msg)
				var channel *Channel
				if channel = DefaultBucket.Channel(sn); channel != nil { //根据sub key找到user对应的channel
					err := channel.Push(msg.MessageProperties, msg.Body)
					if err != nil {
						log.Warn("push msg into %s ch err:%s", sn, err)
					}
				} else {
					//new channel
					channel = NewChannel(Conf.RingSize)
					err := channel.Push(msg.MessageProperties, msg.Body)
					if err != nil {
						log.Warn("push msg into %s new ch err:%s", sn, err)
					}
					DefaultBucket.Put(sn, channel)
				}
				//send to stun
				//dot_count := bytes.Count(d.Body, []byte("."))
				//t := time.Duration(dot_count)
				//time.Sleep(t * time.Second)
				//models.Kick(sn)
				SendText(sn, "stun_queue")
			}
		}
	}()

	<-forever
}
