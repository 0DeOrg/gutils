package rocketmq

/**
 * @Author: lee
 * @Description:
 * @File: rocketmq
 * @Date: 2023-04-26 2:51 下午
 */
import (
	"fmt"
	"github.com/0DeOrg/gutils/dumputils"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"sync"
)

type RocketMQ struct {
	producerCache []*ProducerProxy
	chContent     chan *PublishContent
	cfg           *RocketMQConfig
	mtx           sync.RWMutex
	mapTag        map[string]*ProducerProxy
	idx           int
}

const default_producer_count = 2

func NewRocketMQ(cfg *RocketMQConfig) (*RocketMQ, error) {
	if 0 == cfg.ProducerCount {
		cfg.ProducerCount = default_producer_count
	}
	ret := &RocketMQ{
		producerCache: make([]*ProducerProxy, 0, cfg.ProducerCount),
	}
	for i := 0; i < cfg.ProducerCount; i++ {
		p, err := NewProducerProxy(cfg, i)
		if nil != err {
			return nil, err
		}
		ret.producerCache = append(ret.producerCache, p)
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

func (mq *RocketMQ) DoJob() {
	if len(mq.producerCache) > 0 {
		go mq.goSendThread()
	}

}

func (mq *RocketMQ) goSendThread() {
	defer dumputils.HandlePanic()

	for {
		select {
		case content := <-mq.chContent:
			p := mq.getProducerWithTopic(content.Topic, content.Tag)
			p.chContent <- content
		}

	}

}

func (mq *RocketMQ) getProducerWithTopic(topic string, tag string) *ProducerProxy {
	key := topic + tag
	if proxy, ok := mq.mapTag[key]; ok {
		return proxy
	}

	idx := mq.idx
	ret := mq.producerCache[idx]
	mq.mapTag[key] = ret
	idx++
	if idx >= len(mq.producerCache) {
		idx = 0
	}

	return ret
}
