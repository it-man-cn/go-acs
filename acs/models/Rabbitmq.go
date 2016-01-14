package models

import (
	"github.com/astaxie/beego"
	"github.com/streadway/amqp"
	. "go-acs/acs/log"
)

var amqpUri string = beego.AppConfig.String("amqpurl")

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	done    chan error
}

type MessageProperties struct {
	Headers         amqp.Table
	CorrelationId   string
	ReplyTo         string
	ContentEncoding string
	ContentType     string
	Expiration      string
}

type Message struct {
	MessageProperties
	Body []byte
}

func (r *RabbitMQ) Connect() (err error) {
	r.conn, err = amqp.Dial(amqpUri)
	if err != nil {
		Logger.Info("[amqp] connect error: %s\n", err)
		return err
	}
	r.channel, err = r.conn.Channel()
	if err != nil {
		Logger.Info("[amqp] get channel error: %s\n", err)
		return err
	}
	r.done = make(chan error)
	return nil
}

func (r *RabbitMQ) Publish(exchange, key string, msg Message) (err error) {
	err = r.channel.Publish(exchange, key, false, false,
		amqp.Publishing{
			Headers:         msg.Headers,
			ContentType:     msg.ContentType,
			ContentEncoding: msg.ContentEncoding,
			CorrelationId:   msg.CorrelationId,
			//Expiration:      msg.Expiration,
			Body: msg.Body,
		},
	)
	if err != nil {
		Logger.Info("[amqp] publish message error: %s\n", err)
		return err
	}
	return nil
}

func (r *RabbitMQ) PublishText(exchange, key, text string) (err error) {
	err = r.channel.Publish(exchange, key, false, false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: "",
			ReplyTo:       "",
			Body:          []byte(text),
		},
	)
	if err != nil {
		Logger.Info("[amqp] publish message error: %s\n", err)
		return err
	}
	return nil
}

func (r *RabbitMQ) ConsumeQueue(queue string) (message Message, err error) {
	deliver, ok, err := r.channel.Get(queue, true)
	if err != nil {
		Logger.Info("[amqp] consume queue error: %s\n", err)
		return message, err
	}
	if ok {
		return Message{
			MessageProperties: MessageProperties{
				Headers:         deliver.Headers,
				CorrelationId:   deliver.CorrelationId,
				ReplyTo:         deliver.ReplyTo,
				ContentEncoding: deliver.ContentEncoding,
				ContentType:     deliver.ContentType,
			},
			Body: deliver.Body,
		}, nil
	}
	return message, nil
}

func (r *RabbitMQ) Close() (err error) {
	err = r.conn.Close()
	if err != nil {
		Logger.Info("[amqp] close error: %s\n", err)
		return err
	}
	return nil
}

func ReciveMsg(queue string) (message Message) {
	rabbit := new(RabbitMQ)
	if err := rabbit.Connect(); err != nil {
		Logger.Info("[amqp] connect error: %s\n", err)
		return message
	}
	defer rabbit.Close()
	msg, err := rabbit.ConsumeQueue(queue)
	if err != nil {
		Logger.Info("[amqp] get error: %s\n", err)
		return message
	}
	return msg
}

func Reply(msg Message, sendTo string) (err error) {
	rabbit := new(RabbitMQ)
	if err = rabbit.Connect(); err != nil {
		Logger.Info("[amqp] connect error: %s\n", err)
		return
	}
	defer rabbit.Close()
	err = rabbit.Publish("", sendTo, msg)
	if err != nil {
		Logger.Info("[amqp] send error: %s\n", err)
		return
	}
	return
}

func SendMsg(msg Message, sendTo string) (err error) {
	rabbit := new(RabbitMQ)
	if err = rabbit.Connect(); err != nil {
		Logger.Info("[amqp] connect error: %s\n", err)
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
		Logger.Info("[amqp] QueueDeclare error: %s\n", err)
		return
	}
	err = rabbit.Publish("", sendTo, msg)
	if err != nil {
		Logger.Info("[amqp] send error: %s\n", err)
		return
	}
	return
}

func SendText(text, sendTo string) (err error) {
	rabbit := new(RabbitMQ)
	if err := rabbit.Connect(); err != nil {
		Logger.Info("[amqp] connect error: %s\n", err)
		return err
	}
	defer rabbit.Close()
	err = rabbit.PublishText("", sendTo, text)
	if err != nil {
		Logger.Info("[amqp] send error: %s\n", err)
		return err
	}
	return nil
}
