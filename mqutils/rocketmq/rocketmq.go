package rocketmq

/**
 * @Author: lee
 * @Description:
 * @File: rocketmq
 * @Date: 2023-04-26 2:51 下午
 */
import (
	"context"
	"fmt"
	"github.com/0DeOrg/gutils/dumputils"
	"github.com/0DeOrg/gutils/logutils"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"go.uber.org/zap"
	"sync"
	"time"
)

type RocketMQ struct {
	producerCache []*ProducerProxy
	chContent     chan *PublishContent
	cfg           *RocketMQConfig
	mtx           sync.RWMutex
	mapProducer   map[string]*ProducerProxy
	idx           int
	cancel        context.CancelFunc
}

const default_producer_count = 2

func NewRocketMQ(cfg *RocketMQConfig) (*RocketMQ, error) {
	if 0 == cfg.ProducerCount {
		cfg.ProducerCount = default_producer_count
	}
	ret := &RocketMQ{
		producerCache: make([]*ProducerProxy, 0, cfg.ProducerCount),
		chContent:     make(chan *PublishContent, 1000),
		cfg:           cfg,
		mapProducer:   make(map[string]*ProducerProxy),
	}
	for i := 0; i < cfg.ProducerCount; i++ {
		p, err := NewProducerProxy(cfg, i)
		if nil != err {
			return nil, err
		}
		ret.producerCache = append(ret.producerCache, p)
	}

	if len(ret.producerCache) > 0 {
		go ret.goSendThread()
	}
	return ret, nil
}

func (mq *RocketMQ) Close() (ret []error) {
	var err error
	for idx, pro := range mq.producerCache {
		err = pro.Shutdown()
		if nil != err {
			ret = append(ret, fmt.Errorf("producer '%d' shutdown err: %s", idx, err.Error()))
		}
	}

	return nil
}

func (mq *RocketMQ) RegisterPushConsumer(groupName string, m consumer.MessageModel) (*ConsumerPushProxy, error) {
	proxy, err := NewConsumerPushProxy(groupName, m, mq.cfg)
	if nil != err {
		return nil, err
	}

	return proxy, nil
}

func (mq *RocketMQ) PublishContent(content *PublishContent) {
	mq.chContent <- content
}

func (mq *RocketMQ) goSendThread() {
	defer dumputils.HandlePanic()

	ctx, cancel := context.WithCancel(context.TODO())
	mq.cancel = cancel
	d := 60 * time.Second
	ticker := time.NewTicker(d)
	count := 0
	for {
		select {
		case <-ctx.Done():
			logutils.Warn("RocketMQ context cancel")
			return
		case content := <-mq.chContent:
			p := mq.getProducerWithTopic(content.Topic, content.Tag)
			p.chContent <- content
			count++
		case <-ticker.C:
			logutils.Info("RocketMQ publish report", zap.Int("count", count), zap.Duration("duration", d))
			count = 0
		}

	}

}

func (mq *RocketMQ) getProducerWithTopic(topic string, tag string) *ProducerProxy {
	key := topic + tag
	if proxy, ok := mq.mapProducer[key]; ok {
		return proxy
	}

	idx := mq.idx
	ret := mq.producerCache[idx]
	mq.mapProducer[key] = ret
	idx++
	if idx >= len(mq.producerCache) {
		idx = 0
	}
	mq.idx = idx

	return ret
}
