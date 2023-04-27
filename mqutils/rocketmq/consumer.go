package rocketmq

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

/**
 * @Author: lee
 * @Description:
 * @File: consumer
 * @Date: 2023-04-26 7:42 下午
 */

type ConsumerPushProxy struct {
	cfg      *RocketMQConfig
	consumer rocketmq.PushConsumer
	ctx      context.Context
}

func NewConsumerPushProxy(groupName string, m consumer.MessageModel, cfg *RocketMQConfig) (*ConsumerPushProxy, error) {
	c, err := consumer.NewPushConsumer(
		consumer.WithGroupName(groupName),
		consumer.WithNameServer(cfg.NameServers),
		consumer.WithConsumerModel(m))
	if nil != err {
		return nil, fmt.Errorf("consumer create err: %s", err.Error())
	}

	ret := &ConsumerPushProxy{
		cfg:      cfg,
		consumer: c,
	}

	c.Rebalance()
	return ret, nil
}

func (proxy *ConsumerPushProxy) Subscribe(topic string, selector consumer.MessageSelector, handler func(context.Context, ...*primitive.MessageExt) (consumer.ConsumeResult, error)) error {
	err := proxy.consumer.Subscribe(topic, selector, handler)
	if nil != err {
		return fmt.Errorf("consumer subscribe err: %s", err.Error())
	}

	err = proxy.consumer.Start()

	if nil != err {
		return fmt.Errorf("consumer start err: %s", err.Error())
	}

	return nil
}
