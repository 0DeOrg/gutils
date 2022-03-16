package rabbitmq

/**
 * @Author: lee
 * @Description:
 * @File: rabbitmq
 * @Date: 2022/1/11 5:56 下午
 */

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"math/rand"
	"time"
)

const MAX_POOL_LENGTH = 4

type RabbitMq struct {
	Id          int
	Url         string
	Connection  *amqp.Connection
	channelPool []*amqp.Channel
}

func NewRabbitMq(user, pwd, address, vHost string) (*RabbitMq, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s/%s", user, pwd, address, vHost)
	ret := &RabbitMq{
		Url:         url,
		channelPool: make([]*amqp.Channel, MAX_POOL_LENGTH),
	}

	err := ret.connect()
	if nil != err {
		return nil, err
	}
	return ret, nil
}

func (rq *RabbitMq) connect() error {
	if nil != rq.Connection {
		rq.Connection.Close()
	}

	conn, err := amqp.Dial(rq.Url)
	if err != nil {
		return err
	}

	for i := 0; i < cap(rq.channelPool); i++ {
		ch, err := conn.Channel()
		if err != nil {
			return err
		}

		rq.channelPool = append(rq.channelPool, ch)
	}
	rq.Connection = conn

	return nil
}

func (rq *RabbitMq) Close() error {
	return rq.Connection.Close()
}

func (rq *RabbitMq) getChannel() *amqp.Channel {
	i := rand.Intn(MAX_POOL_LENGTH)
	ch := rq.channelPool[i]
	return ch
}

func (rq *RabbitMq) ExchangeDeclare(name string, kind ExchangeKind, durable bool) error {
	return rq.getChannel().ExchangeDeclare(name, string(kind), durable, false, false, false, nil)
}

func (rq *RabbitMq) Publish(content *PublishContent, reliable bool) (confirmed bool, err error) {
	contentType := content.ContentType
	if "" == contentType {
		contentType = "text/json"
	}

	amqpChannel := rq.getChannel()

	if reliable {
		err = amqpChannel.Confirm(false)
		confirm := make(chan amqp.Confirmation, 1)
		if nil != err {
			//不支持confirm
			close(confirm)
		} else {
			confirm = amqpChannel.NotifyPublish(confirm)
			defer confirmOne(err, confirm, &confirmed)
		}
	}

	err = amqpChannel.Publish(content.ExchangeName, content.RoutingKey, true, false,
		amqp.Publishing{
			ContentType: contentType,
			Timestamp:   time.Now(),
			Body:        content.Content,
		})
	return
}

func confirmOne(err error, confirm <-chan amqp.Confirmation, confirmed *bool) {
	if nil != err {
		*confirmed = false
		return
	}
	select {
	case ack := <-confirm:
		{
			if ack.Ack {
				*confirmed = true
				log.Println("confirmOne success tag:", ack.DeliveryTag)
			} else {
				*confirmed = false
				log.Println("confirmOne fatal tag:", ack.DeliveryTag)
			}
		}
	case <-time.After(5 * time.Second):
		{
			log.Println("confirmOne timeout")
		}
	}
}