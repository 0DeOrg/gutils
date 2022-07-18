package rabbitmq

import (
	"context"
	"fmt"
	"github.com/streadway/amqp"
	"gitlab.qihangxingchen.com/qt/gutils/logutils"
	"go.uber.org/zap"
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

type connectionProxy struct {
	url           string
	conn          *amqp.Connection //connection
	channelPool   []*channelProxy  //信道pool
	notifyClose   chan *amqp.Error //错误chan
	done          chan struct{}    //
	ctx           context.Context
	cancel        context.CancelFunc
	channelInited bool
	chOffset      int
}

func NewConnectionProxy(url string) *connectionProxy {
	ctx, cancel := context.WithCancel(context.TODO())
	ret := &connectionProxy{
		url:           url,
		channelPool:   make([]*channelProxy, 0, maxChannelCountPerConnection),
		done:          make(chan struct{}, 1),
		ctx:           ctx,
		cancel:        cancel,
		channelInited: false,
		chOffset:      -1,
	}

	ret.initChannel()

	go ret.handleReconnect()
	return ret
}

func (p *connectionProxy) waitForFirstConnect() {
	for idx, chProxy := range p.channelPool {
		select {
		case <-chProxy.flagCh:
			logutils.Info("connectionProxy waitForFirstConnect success", zap.Int("id", idx))
		case <-time.After(30 * time.Second):
			logutils.Fatal("connectionProxy waitForFirstConnect timeout 30 seconds")
		}
	}
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

func (p *connectionProxy) getUnusedChannel() *amqp.Channel {
	idx := 0
	for i := 0; i < maxChannelCountPerConnection; i++ {
		idx = (p.chOffset + 1) % maxChannelCountPerConnection
		if p.channelPool[idx].inited {
			p.chOffset = idx
			return p.channelPool[idx].ch
		}
	}

	return nil
}

func (p *connectionProxy) connect() error {
	if nil != p.conn {
		//p.conn.Close()
	}

	conn, err := amqp.Dial(p.url)
	if err != nil {
		return err
	}

	p.conn = conn
	p.notifyClose = make(chan *amqp.Error, 1)
	conn.NotifyClose(p.notifyClose)

	return nil
}

func (p *connectionProxy) initChannel() {
	if p.channelInited {
		return
	}
	for i := 0; i < maxChannelCountPerConnection; i++ {
		flagCh := make(chan struct{}, 1)
		ch := NewChannelProxy(p.ctx, p, flagCh)
		p.channelPool = append(p.channelPool, ch)
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
	flagCh        chan struct{}
}

func NewChannelProxy(ctx context.Context, connProxy *connectionProxy, flagCh chan struct{}) *channelProxy {
	ret := &channelProxy{
		id:        -1,
		connProxy: connProxy,
		ctx:       ctx,
		flagCh:    flagCh,
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

		if !p.inited {
			p.flagCh <- struct{}{}
		}
		p.inited = true

		logutils.Warn("channelProxy handleReconnect success")

		select {
		case <-p.ctx.Done():
			break
		case <-p.notifyChClose:
			logutils.Warn("channelProxy handleReconnect notify ")
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

	p.ch = ch
	p.notifyChClose = make(chan *amqp.Error, 1)
	p.ch.NotifyClose(p.notifyChClose)

	return ch, nil
}
