package rocketmq

import (
	"context"
	"fmt"
	"github.com/0DeOrg/gutils/dumputils"
	"github.com/0DeOrg/gutils/logutils"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"go.uber.org/zap"
)

/**
 * @Author: lee
 * @Description:
 * @File: producer
 * @Date: 2023-04-26 6:35 下午
 */

type ProducerProxy struct {
	producer  rocketmq.Producer
	cfg       *RocketMQConfig
	chContent chan *PublishContent
	idx       int
	cancel    context.CancelFunc
}

func NewProducerProxy(cfg *RocketMQConfig, idx int) (*ProducerProxy, error) {
	p, err := rocketmq.NewProducer(
		producer.WithNameServer(cfg.NameServers),
		producer.WithGroupName(cfg.ProducerGroup),
	)

	if nil != err {
		return nil, fmt.Errorf("producer create err: %s", err.Error())
	}

	err = p.Start()
	if nil != err {
		return nil, fmt.Errorf("producer start err: %s", err.Error())
	}

	ret := &ProducerProxy{
		producer:  p,
		cfg:       cfg,
		chContent: make(chan *PublishContent, 1000),
		idx:       idx,
	}

	go ret.goSendThread()
	return ret, nil
}

func (proxy *ProducerProxy) Shutdown() error {
	if nil != proxy.cancel {
		proxy.cancel()
	}
	return proxy.producer.Shutdown()
}

func (proxy *ProducerProxy) PublishContent(content *PublishContent) {
	proxy.chContent <- content
}

func (proxy *ProducerProxy) goSendThread() {
	defer dumputils.HandlePanic(func() {
		logutils.Info("ProducerProxy|goSendThread finish")
	})
	ctx, cancel := context.WithCancel(context.TODO())
	proxy.cancel = cancel
	for {
		select {
		case content := <-proxy.chContent:
			msg := &primitive.Message{
				Topic: content.Topic,
				Body:  content.Body,
			}
			if "" != content.Tag {
				msg.WithTag(content.Tag)
			}
			_, err := proxy.producer.SendSync(ctx, msg)
			if nil != err {
				logutils.Error("ProducerProxy|SendSync err", zap.Error(err), zap.Int("idx", proxy.idx))
			}
		case <-ctx.Done():
			return
		}
	}
}
