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
	"go.uber.org/zap"
	"gutils/dumputils"
	"gutils/logutils"
	"log"
	"time"
)

//channel 容量
const capicity_publish_ch = 10000

//最大堵塞数量
const max_traffic_count = capicity_publish_ch / 2

type RabbitMq struct {
	Id        int
	Url       string
	connProxy *connectionProxy
	publishCh chan *PublishContent
}

func NewRabbitMq(cfg *RabbitMQConfig) (*RabbitMq, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s/%s", cfg.User, cfg.Password, cfg.Address, cfg.VHost)
	ret := &RabbitMq{
		Url:       url,
		publishCh: make(chan *PublishContent, capicity_publish_ch),
	}

	conn := NewConnectionProxy(url)
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

func (rq *RabbitMq) getChannel() *amqp.Channel {

	return rq.connProxy.getUnusedChannel()
}

func (rq *RabbitMq) ExchangeDeclare(name string, kind ExchangeKind, durable bool) error {
	return rq.getChannel().ExchangeDeclare(name, string(kind), durable, false, false, false, nil)
}

func (rq *RabbitMq) QueueDeclare(name string) error {
	_, err := rq.getChannel().QueueDeclare(name, false, false, false, false, nil)
	return err
}

func (rq *RabbitMq) QueueBind(name, exchange, routingKey string) error {
	return rq.getChannel().QueueBind(name, routingKey, exchange, false, nil)
}

func (rq *RabbitMq) Consume(name string) (<-chan amqp.Delivery, error) {
	return rq.getChannel().Consume(name, "", false, false, false, false, nil)
}

func (rq *RabbitMq) Process() {
	defer dumputils.HandlePanic()
	idx := 0
	go func() {
		for {
			select {
			case content := <-rq.publishCh:
				idx++
				for {
					ch := rq.getChannel()
					if nil == ch {
						time.Sleep(time.Second)
						continue
					}

					//每100次打一次日志
					if 100 == idx {
						logutils.Info("mqpublish 100 times", zap.Int("ch len", len(rq.publishCh)),
							zap.String("content", string(content.Content)))
					}

					_, err := rq.Publish(content, false, ch)
					if nil != err {
						logutils.Error("RabbitMq Publish fatal", zap.Any("content", content), zap.Error(err))

						//当达到最大堵塞数量时，不堵塞了 防止影响正常流程，mq推送暂时就不保证了
						if len(rq.publishCh) > max_traffic_count {
							logutils.Warn("RabbitMq publish has reach max traffic count", zap.Int("traffic", len(rq.publishCh)))
							break
						}
						success := false
						for idx := 0; idx < maxChannelCountPerConnection-1; idx++ {
							ch = rq.getChannel()
							if nil == ch {
								time.Sleep(time.Second)
								continue
							}
							_, err = rq.Publish(content, false, ch)
							if nil == err {
								success = true
								break
							} else {
								logutils.Error("RabbitMq publish retry failed", zap.Int("retry", idx+1))
							}
						}

						if !success {
							time.Sleep(time.Second)
							continue
						} else {
							break
						}
					} else {
						break
					}
				}

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

func (rq *RabbitMq) Publish(content *PublishContent, reliable bool, amqpChannel *amqp.Channel) (confirmed bool, err error) {
	if nil == amqpChannel {
		return false, fmt.Errorf("channel is not ready")
	}
	contentType := content.ContentType
	if "" == contentType {
		contentType = "text/json"
	}

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
