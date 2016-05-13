package main

import (
	log "github.com/it-man-cn/log4go"
	"github.com/streadway/amqp"
)

//RabbitMQ conn info
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	done    chan error
}

var rabbit RabbitMQ

//Connect create a new amqp conn
func (r *RabbitMQ) Connect() (err error) {
	r.conn, err = amqp.Dial(Conf.AMQPUrl)
	if err != nil {
		log.Info("[amqp] connect error: %s\n", err)
		return err
	}
	r.channel, err = r.conn.Channel()
	if err != nil {
		log.Info("[amqp] get channel error: %s\n", err)
		return err
	}
	r.done = make(chan error)
	return nil
}

//Publish publish amqp msg
func (r *RabbitMQ) publish(exchange, key string, msg Message) (err error) {
	err = r.channel.Publish(exchange, key, false, false,
		amqp.Publishing{
			Headers:         msg.Headers,
			ContentType:     msg.ContentType,
			ContentEncoding: msg.ContentEncoding,
			CorrelationId:   msg.CorrelationID,
			//Expiration:      msg.Expiration,
			Body: msg.Body,
		},
	)
	if err != nil {
		log.Info("[amqp] publish message error: %s", err)
		return err
	}
	return nil
}

//PublishText publish text msg
func (r *RabbitMQ) publishText(exchange, key, text string) (err error) {
	err = r.channel.Publish(exchange, key, false, false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: "",
			ReplyTo:       "",
			Body:          []byte(text),
		},
	)
	if err != nil {
		log.Info("[amqp] publish message error: %s", err)
		return err
	}
	return nil
}

//ConsumeQueue bindding queue
func (r *RabbitMQ) consumeQueue(queue string) (message Message, err error) {
	deliver, ok, err := r.channel.Get(queue, true)
	if err != nil {
		log.Info("[amqp] consume queue error: %s", err)
		return message, err
	}
	if ok {
		return Message{
			MessageProperties: MessageProperties{
				Headers:         deliver.Headers,
				CorrelationID:   deliver.CorrelationId,
				ReplyTo:         deliver.ReplyTo,
				ContentEncoding: deliver.ContentEncoding,
				ContentType:     deliver.ContentType,
			},
			Body: deliver.Body,
		}, nil
	}
	return message, nil
}

//Close close conn
func (r *RabbitMQ) Close() (err error) {
	err = r.conn.Close()
	if err != nil {
		log.Info("[amqp] close error: %s", err)
		return err
	}
	return nil
}

//ReciveMsg query msg
func ReciveMsg(queue string) (message Message) {
	rabbit := new(RabbitMQ)
	if err := rabbit.Connect(); err != nil {
		log.Info("[amqp] connect error: %s", err)
		return message
	}
	defer rabbit.Close()
	msg, err := rabbit.consumeQueue(queue)
	if err != nil {
		log.Info("[amqp] get error: %s", err)
		return message
	}
	return msg
}

//Reply rpc reply
func Reply(msg Message, sendTo string) (err error) {
	rabbit := new(RabbitMQ)
	if err = rabbit.Connect(); err != nil {
		log.Info("[amqp] connect error: %s", err)
		return
	}
	defer rabbit.Close()
	err = rabbit.publish("", sendTo, msg)
	if err != nil {
		log.Info("[amqp] send error: %s", err)
		return
	}
	return
}

//SendMsg send msg
func SendMsg(msg Message, sendTo string) (err error) {
	rabbit := new(RabbitMQ)
	if err = rabbit.Connect(); err != nil {
		log.Info("[amqp] connect error: %s", err)
		return
	}
	defer rabbit.Close()
	_, err = rabbit.channel.QueueDeclare(
		sendTo, // name
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		amqp.Table{"x-message-ttl": int32(300000)}, // arguments in ms
	)
	if err != nil {
		log.Info("[amqp] QueueDeclare error: %s", err)
		return
	}
	err = rabbit.publish("", sendTo, msg)
	if err != nil {
		log.Info("[amqp] send error: %s", err)
		return
	}
	return
}

//SendText send text
func SendText(text, sendTo string) (err error) {
	rabbit := new(RabbitMQ)
	if err := rabbit.Connect(); err != nil {
		log.Info("[amqp] connect error: %s", err)
		return err
	}
	defer rabbit.Close()
	err = rabbit.publishText("", sendTo, text)
	if err != nil {
		log.Info("[amqp] send error: %s", err)
		return err
	}
	return nil
}

//Kick kick stun
func Kick(sn string) {
	_, err := rabbit.channel.QueueDeclare(
		"stun_queue", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,
		//amqp.Table{"x-message-ttl": int32(300000)}, // arguments in ms
	)

	if err != nil {
		log.Info("[amqp] QueueDeclare error: %s", err)
		return
	}
	//stun kick
	props := MessageProperties{
		CorrelationID:   "",
		ReplyTo:         "",
		ContentEncoding: "UTF-8",
		ContentType:     "text/plain",
	}

	kickmsg := Message{MessageProperties: props, Body: []byte(sn)}

	err = rabbit.publish("", "stun_queue", kickmsg)
	if err != nil {
		log.Info("[amqp] send error: %s", err)
		return
	}
	return
}
