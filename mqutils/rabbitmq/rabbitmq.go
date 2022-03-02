package rabbitmq

/**
 * @Author: lee
 * @Description:
 * @File: rabbitmq
 * @Date: 2022/1/11 5:56 下午
 */

import (
	"github.com/streadway/amqp"
)

type rabbitMq struct {
	Id          int
	Uris        []string
	Connection  *amqp.Connection
	amqpChannel *amqp.Channel
	//publishChan chan *PublishContent
}

func NewRabbitMq() *rabbitMq {

	return nil
}
