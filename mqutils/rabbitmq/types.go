package rabbitmq

import (
	"context"
	"fmt"
	"github.com/streadway/amqp"
	"gitlab.qihangxingchen.com/qt/gutils/logutils"
	"go.uber.org/zap"
	"math/rand"
	"time"
)

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

const (
	// When reconnecting to the server after connection failure
	reconnectDelay = 3 * time.Second

	// When setting up the channel after a channel exception
	reInitDelay = 1 * time.Second

	// When resending messages the server didn't confirm
	resendDelay = 5 * time.Second
)

const (
	maxChannelCountPerConnection = 4
)

type RabbitMQConfig struct {
	User        string `json:"user"     yaml:"user"   mapstructure:"user"`
	Password    string `json:"password"     yaml:"password"   mapstructure:"password"`
	AddressList string `json:"address-list"     yaml:"address-list"   mapstructure:"address-list"`
	VHost       string `json:"vhost"     yaml:"vhost"   mapstructure:"vhost"`
}

type PublishContent struct {
	ExchangeName string
	RoutingKey   string
	Content      []byte
	ContentType  string
}

type connectionProxy struct {
	urlCache      []string
	conn          *amqp.Connection //connection
	channelPool   []*channelProxy  //信道pool
	notifyClose   chan *amqp.Error //错误chan
	done          chan struct{}    //
	ctx           context.Context
	cancel        context.CancelFunc
	channelInited bool
	chOffset      int
	reliable      bool
}

func NewConnectionProxy(urls []string, reliable bool) *connectionProxy {
	ctx, cancel := context.WithCancel(context.TODO())
	ret := &connectionProxy{
		urlCache:      urls,
		channelPool:   make([]*channelProxy, 0, maxChannelCountPerConnection),
		done:          make(chan struct{}, 1),
		ctx:           ctx,
		cancel:        cancel,
		channelInited: false,
		chOffset:      -1,
	}

	//需要将channel创建好存入channelPool
	ret.initChannel(reliable)

	go ret.handleReconnect()
	ret.waitForFirstConnect()
	return ret
}

func (p *connectionProxy) waitForFirstConnect() {
	for idx, chProxy := range p.channelPool {
		select {
		case <-chProxy.flagCh:
			logutils.Info("connectionProxy waitForFirstConnect success", zap.Int("id", idx))
		case <-time.After(10 * time.Second):
			logutils.Fatal("connectionProxy waitForFirstConnect timeout 30 seconds")
		}
	}
}

func (p *connectionProxy) chooseUrl() string {
	idx := rand.Intn(len(p.urlCache))
	return p.urlCache[idx]
}

func (p *connectionProxy) registerChannel() int {
	if !p.channelInited {
		return -1
	}
	for i := 0; i < maxChannelCountPerConnection; i++ {
		chProxy := p.channelPool[i]
		if -1 != chProxy.id {
			chProxy.id = i
			return i
		}
	}
	return -1
}

func (p *connectionProxy) getUnusedChannel() *channelProxy {
	idx := 0
	for i := 0; i < maxChannelCountPerConnection; i++ {
		idx = (p.chOffset + 1) % maxChannelCountPerConnection
		if p.channelPool[idx].inited && p.channelPool[idx].running {
			p.chOffset = idx
			return p.channelPool[idx]
		}
	}

	return nil
}

func (p *connectionProxy) connect() error {
	if nil != p.conn {
		err := p.conn.Close()
		if nil != err {
			logutils.Error("connectionProxy|connect close conn failed", zap.Error(err))
		}
	}

	url := p.chooseUrl()
	conn, err := amqp.Dial(url)
	if err != nil {
		return err
	}

	p.conn = conn
	notifyClose := make(chan *amqp.Error, 1)
	p.notifyClose = conn.NotifyClose(notifyClose)

	return nil
}

func (p *connectionProxy) initChannel(reliable bool) {
	if p.channelInited {
		return
	}
	for i := 0; i < maxChannelCountPerConnection; i++ {
		flagCh := make(chan struct{}, 1)
		chProxy := NewChannelProxy(p, flagCh, reliable)
		p.channelPool = append(p.channelPool, chProxy)
	}

	p.channelInited = true
}

func (p *connectionProxy) handleReconnect() {
	for {
		err := p.connect()
		if nil != err {
			select {
			case <-p.done:
				break
			case <-time.After(reconnectDelay):
				continue
			}
		}

		logutils.Warn("connectionProxy handleReconnect success")
		select {
		case <-p.notifyClose:
			logutils.Warn("connectionProxy handleReconnect notify ")
			continue
		case <-p.done:
			//断开连接时,终止所有的channel
			p.cancel()
			break
		}

	}
}

type channelProxy struct {
	id            int
	ch            *amqp.Channel
	connProxy     *connectionProxy
	notifyChClose chan *amqp.Error
	ctx           context.Context
	inited        bool
	running       bool
	flagCh        chan struct{}
	reliable      bool
	confirmCh     chan amqp.Confirmation
}

func NewChannelProxy(connProxy *connectionProxy, flagCh chan struct{}, reliable bool) *channelProxy {
	ret := &channelProxy{
		id:        -1,
		connProxy: connProxy,
		ctx:       connProxy.ctx,
		flagCh:    flagCh,
		reliable:  reliable,
	}

	go ret.handleReconnect()

	return ret
}

func (p *channelProxy) handleReconnect() {
	for {
		_, err := p.init()
		if nil != err {
			select {
			case <-p.ctx.Done():
				break
			case <-time.After(reInitDelay):
				continue
			}
		}

		logutils.Warn("channelProxy handleReconnect success")

		select {
		case <-p.ctx.Done():
			break
		case <-p.notifyChClose:
			p.running = false
			logutils.Warn("channelProxy handleReconnect notifyChClose ")
			continue
		}
	}
}

func (p *channelProxy) init() (*amqp.Channel, error) {
	if nil == p.connProxy.conn {
		return nil, fmt.Errorf("connection is nil")
	}
	ch, err := p.connProxy.conn.Channel()
	if nil != err {
		return nil, err
	}
	notifyClose := make(chan *amqp.Error, 1)
	p.notifyChClose = ch.NotifyClose(notifyClose)
	if p.reliable {
		confirmCh := make(chan amqp.Confirmation, 1)
		err = ch.Confirm(false)
		if nil != err {
			return nil, fmt.Errorf("channelProxy|init|Confirm err: %s", err.Error())
		} else {
			p.confirmCh = ch.NotifyPublish(confirmCh)
		}
	}

	p.ch = ch
	p.running = true
	if !p.inited {
		p.flagCh <- struct{}{}
	}
	p.inited = true

	return ch, nil
}

func (p *channelProxy) Publish(content *PublishContent) (confirmed bool, deliveryTag uint64, err error) {
	amqpCh := p.ch
	if nil == amqpCh {
		return false, 0, fmt.Errorf("channel is not ready")
	}
	contentType := content.ContentType
	if "" == contentType {
		contentType = "text/json"
	}

	if p.reliable {
		defer confirmOne(&err, p.confirmCh, &confirmed, &deliveryTag)
	}

	err = amqpCh.Publish(content.ExchangeName, content.RoutingKey, false, false,
		amqp.Publishing{
			ContentType: contentType,
			Timestamp:   time.Now(),
			Body:        content.Content,
		})
	return
}

func (p *channelProxy) ExchangeDeclare(name string, kind ExchangeKind, durable bool) error {
	if nil == p.ch {
		return fmt.Errorf("channelProxy|ExchangeDeclare channel is nil")
	}
	if p.running {
		return fmt.Errorf("channelProxy|ExchangeDeclare channel is not running")
	}

	err := p.ch.ExchangeDeclare(name, string(kind), durable, false, false, false, nil)
	if nil != err {
		return fmt.Errorf("channelProxy|ExchangeDeclare failed, err: %s", err.Error())
	}
	return nil
}

func (p *channelProxy) QueueDeclare(name string) error {
	if nil == p.ch {
		return fmt.Errorf("channelProxy|QueueDeclare channel is nil")
	}
	if p.running {
		return fmt.Errorf("channelProxy|QueueDeclare channel is not running")
	}

	_, err := p.ch.QueueDeclare(name, false, false, false, false, nil)
	if nil != err {
		return fmt.Errorf("channelProxy|QueueDeclare failed, err: %s", err.Error())
	}

	return nil
}

func (p *channelProxy) QueueBind(name, exchange, routingKey string) error {
	if nil == p.ch {
		return fmt.Errorf("channelProxy|QueueBind channel is nil")
	}
	if p.running {
		return fmt.Errorf("channelProxy|QueueBind channel is not running")
	}

	err := p.ch.QueueBind(name, routingKey, exchange, false, nil)
	if nil != err {
		return fmt.Errorf("channelProxy|QueueBind failed, err: %s", err.Error())
	}

	return nil
}
func (p *channelProxy) Consume(name string) (<-chan amqp.Delivery, error) {
	if nil == p.ch {
		return nil, fmt.Errorf("channelProxy|Consume channel is nil")
	}
	if p.running {
		return nil, fmt.Errorf("channelProxy|Consume channel is not running")
	}

	ret, err := p.ch.Consume(name, "", false, false, false, false, nil)
	if nil != err {
		return nil, fmt.Errorf("channelProxy|Consume failed, err: %s", err.Error())
	}

	return ret, nil
}
