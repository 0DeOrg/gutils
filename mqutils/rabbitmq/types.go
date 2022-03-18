package rabbitmq

import "github.com/streadway/amqp"

/**
 * @Author: lee
 * @Description:
 * @File: types
 * @Date: 2022/1/11 5:57 下午
 */

type ExchangeKind string

const (
	ExchangeFanout  = amqp.ExchangeFanout
	ExchangeDirect  = amqp.ExchangeDirect
	ExchangeTopic   = amqp.ExchangeTopic
	ExchangeHeaders = amqp.ExchangeHeaders
)

type RabbitMQConfig struct {
	User     string `json:"user"     yaml:"user"   mapstructure:"user"`
	Password string `json:"password"     yaml:"password"   mapstructure:"password"`
	Address  string `json:"address"     yaml:"address"   mapstructure:"address"`
	VHost    string `json:"vhost"     yaml:"vhost"   mapstructure:"vhost"`
}

type PublishContent struct {
	ExchangeName string
	RoutingKey   string
	Content      []byte
	ContentType  string
}
