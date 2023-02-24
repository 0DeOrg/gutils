package rabbitmq

/**
 * @Author: lee
 * @Description:
 * @File: rabbitmq
 * @Date: 2022/1/11 5:56 下午
 */

import (
	"fmt"
	"github.com/0DeOrg/gutils/dumputils"
	"github.com/0DeOrg/gutils/logutils"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"time"
)

//channel 容量
const capicity_publish_ch = 10000

//最大堵塞数量
const max_traffic_count = capicity_publish_ch / 2

type RabbitMq struct {
	Id        int
	connProxy *connectionProxy
	publishCh chan *PublishContent
}

func NewRabbitMq(cfg *RabbitMQConfig, reliable bool) (*RabbitMq, error) {
	urls := make([]string, 0, len(cfg.Addresses))
	for _, addr := range cfg.Addresses {
		url := fmt.Sprintf("amqp://%s:%s@%s/%s", cfg.User, cfg.Password, addr, cfg.VHost)
		urls = append(urls, url)
	}

	ret := &RabbitMq{
		publishCh: make(chan *PublishContent, capicity_publish_ch),
	}

	conn := NewConnectionProxy(urls, reliable)
	ret.connProxy = conn

	//time.AfterFunc(time.Minute, func() {
	//	ret.connProxy.conn.Close()
	//})
	return ret, nil
}

func (rq *RabbitMq) Close() error {
	rq.connProxy.done <- struct{}{}
	return nil
}

func (rq *RabbitMq) WaitForFirstConnect() {
	rq.connProxy.waitForFirstConnect()
}

func (rq *RabbitMq) getChannelProxy() *channelProxy {
	return rq.connProxy.getUnusedChannel()
}

func (rq *RabbitMq) ExchangeDeclare(name string, kind ExchangeKind, durable bool) error {
	chProxy := rq.getChannelProxy()
	if nil == chProxy {
		return fmt.Errorf("RabbitMq|ExchangeDeclare chProxy is nil")
	}

	return chProxy.ExchangeDeclare(name, kind, durable)
}

func (rq *RabbitMq) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool) error {
	chProxy := rq.getChannelProxy()
	if nil == chProxy {
		return fmt.Errorf("RabbitMq|QueueDeclare chProxy is nil")
	}

	return chProxy.QueueDeclare(name, durable, autoDelete, exclusive, noWait)
}

func (rq *RabbitMq) QueueBind(name, exchange, routingKey string) error {
	chProxy := rq.getChannelProxy()
	if nil == chProxy {
		return fmt.Errorf("RabbitMq|QueueBind chProxy is nil")
	}

	return chProxy.QueueBind(name, exchange, routingKey)
}

func (rq *RabbitMq) Consume(name string) (<-chan amqp.Delivery, error) {
	chProxy := rq.getChannelProxy()
	if nil == chProxy {
		return nil, fmt.Errorf("RabbitMq|Consume chProxy is nil")
	}

	return chProxy.Consume(name)
}

func (rq *RabbitMq) Process() {
	defer dumputils.HandlePanic()
	ticker := time.NewTicker(60 * time.Second)
	cnSuccess := 0
	cnFailed := 0
	go func() {
		for {
			select {
			case content := <-rq.publishCh:
				for {
					_, _, err := rq.Publish(content)
					if nil != err {
						logutils.Error("RabbitMq|Process Publish fatal", zap.Any("content", content), zap.Error(err))
						//当达到最大堵塞数量时，不堵塞了 防止影响正常流程，mq推送暂时就不保证了
						if len(rq.publishCh) > max_traffic_count {
							logutils.Warn("RabbitMq publish has reach max traffic count", zap.Int("traffic", len(rq.publishCh)))
							break
						}
						time.Sleep(time.Second)
						cnFailed++
						continue
					}

					cnSuccess++
					break
				}
			case <-ticker.C:
				logutils.Info("RabbitMq|Process mq publish report", zap.Int("traffic", len(rq.publishCh)),
					zap.Int("success", cnSuccess), zap.Int("failed", cnFailed))
				cnSuccess = 0
				cnFailed = 0
			}
		}
	}()

}

func (rq *RabbitMq) PublishContent(exchangeName string, routingKey string, content []byte) {
	publishContent := &PublishContent{
		ExchangeName: exchangeName,
		RoutingKey:   routingKey,
		Content:      content,
	}

	rq.publishCh <- publishContent
}

func (rq *RabbitMq) Publish(content *PublishContent) (confirmed bool, deliveryTag uint64, err error) {
	chProxy := rq.getChannelProxy()
	if nil == chProxy {
		return false, 0, fmt.Errorf("RabbitMq|Publish chProxy is nil")
	}

	return chProxy.Publish(content)
}

func confirmOne(err *error, confirm <-chan amqp.Confirmation, confirmed *bool, deliveryTag *uint64) {
	if nil != *err {
		*confirmed = false
		return
	}
	select {
	case ack := <-confirm:
		{
			if ack.Ack {
				*confirmed = true

				*deliveryTag = ack.DeliveryTag
			} else {
				*confirmed = false
				*err = fmt.Errorf("confirmOne failed")
				*deliveryTag = ack.DeliveryTag
			}
		}
	case <-time.After(3 * time.Second):
		{
			*err = fmt.Errorf("confirmOne timeout")
		}
	}
}
