package models

import (
	"github.com/astaxie/beego"
	"github.com/streadway/amqp"
	log "go-acs/acs/log"
)

var amqpURI = beego.AppConfig.String("amqpurl")

//RabbitMQ conn info
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	done    chan error
}

//MessageProperties amqp msg property
type MessageProperties struct {
	Headers         amqp.Table
	CorrelationID   string
	ReplyTo         string
	ContentEncoding string
	ContentType     string
	Expiration      string
}

//Message amqp msg
type Message struct {
	MessageProperties
	Body []byte
}

//Connect create a new amqp conn
func (r *RabbitMQ) Connect() (err error) {
	r.conn, err = amqp.Dial(amqpURI)
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
		log.Info("[amqp] publish message error: %s\n", err)
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
		log.Info("[amqp] publish message error: %s\n", err)
		return err
	}
	return nil
}

//ConsumeQueue bindding queue
func (r *RabbitMQ) consumeQueue(queue string) (message Message, err error) {
	deliver, ok, err := r.channel.Get(queue, true)
	if err != nil {
		log.Info("[amqp] consume queue error: %s\n", err)
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
		log.Info("[amqp] close error: %s\n", err)
		return err
	}
	return nil
}

//ReciveMsg query msg
func ReciveMsg(queue string) (message Message) {
	rabbit := new(RabbitMQ)
	if err := rabbit.Connect(); err != nil {
		log.Info("[amqp] connect error: %s\n", err)
		return message
	}
	defer rabbit.Close()
	msg, err := rabbit.consumeQueue(queue)
	if err != nil {
		log.Info("[amqp] get error: %s\n", err)
		return message
	}
	return msg
}

//Reply rpc reply
func Reply(msg Message, sendTo string) (err error) {
	rabbit := new(RabbitMQ)
	if err = rabbit.Connect(); err != nil {
		log.Info("[amqp] connect error: %s\n", err)
		return
	}
	defer rabbit.Close()
	err = rabbit.publish("", sendTo, msg)
	if err != nil {
		log.Info("[amqp] send error: %s\n", err)
		return
	}
	return
}

//SendMsg send msg
func SendMsg(msg Message, sendTo string) (err error) {
	rabbit := new(RabbitMQ)
	if err = rabbit.Connect(); err != nil {
		log.Info("[amqp] connect error: %s\n", err)
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
		log.Info("[amqp] QueueDeclare error: %s\n", err)
		return
	}
	err = rabbit.publish("", sendTo, msg)
	if err != nil {
		log.Info("[amqp] send error: %s\n", err)
		return
	}
	return
}

//SendText send text
func SendText(text, sendTo string) (err error) {
	rabbit := new(RabbitMQ)
	if err := rabbit.Connect(); err != nil {
		log.Info("[amqp] connect error: %s\n", err)
		return err
	}
	defer rabbit.Close()
	err = rabbit.publishText("", sendTo, text)
	if err != nil {
		log.Info("[amqp] send error: %s\n", err)
		return err
	}
	return nil
}
