package rocketmq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/0DeOrg/gutils/dumputils"
	"github.com/0DeOrg/gutils/logutils"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"go.uber.org/zap"
	"os"
	"strconv"
	"sync"
	"time"
)

/**
 * @Author: lee
 * @Description:
 * @File: producer
 * @Date: 2023-04-26 6:35 下午
 */

type ProducerProxy struct {
	producer       rocketmq.Producer
	cfg            *RocketMQConfig
	chContent      chan *PublishContent
	idx            int
	batchChanCache map[string]chan *primitive.Message
	cancel         context.CancelFunc
}

func NewProducerProxy(cfg *RocketMQConfig, idx int) (*ProducerProxy, error) {
	instName := strconv.Itoa(os.Getpid()) + "_" + strconv.Itoa(idx)
	groupName := cfg.ProducerGroup
	p, err := rocketmq.NewProducer(
		producer.WithNameServer(cfg.NameServers),
		producer.WithGroupName(groupName),
		producer.WithInstanceName(instName),
	)

	if nil != err {
		return nil, fmt.Errorf("producer create err: %s", err.Error())
	}

	err = p.Start()
	if nil != err {
		return nil, fmt.Errorf("producer start err: %s", err.Error())
	}

	ret := &ProducerProxy{
		producer:       p,
		cfg:            cfg,
		chContent:      make(chan *PublishContent, 1000),
		batchChanCache: make(map[string]chan *primitive.Message),
		idx:            idx,
	}

	req := &ReqPing{
		InstName:  instName,
		GroupName: groupName,
	}
	if err := ret.Ping(req); nil != err {
		return nil, err
	}

	go ret.goDispatchThread()
	return ret, nil
}

func (proxy *ProducerProxy) Shutdown() error {
	if nil != proxy.cancel {
		proxy.cancel()
	}
	return proxy.producer.Shutdown()
}

func (proxy *ProducerProxy) PublishContent(content *PublishContent) {
	defer func() {
		//管道关闭时会调用到
		if r := recover(); nil != r {

		}
	}()
	proxy.chContent <- content
}

func (proxy *ProducerProxy) Ping(req *ReqPing) error {
	body, _ := json.Marshal(req)
	msg := &primitive.Message{
		Topic: TopicPing,
		Body:  body,
	}

	if req.GroupName != "" {
		msg.WithKeys([]string{req.GroupName})
	}

	_, err := proxy.producer.SendSync(context.TODO(), msg)

	if nil != err {
		return fmt.Errorf("ProducerProxy ping err: %s", err.Error())
	}

	return nil
}
func (proxy *ProducerProxy) goSendThread(ctx context.Context, ch chan []*primitive.Message) {
	defer dumputils.HandlePanic(func() {
		logutils.Warn("ProducerProxy|goSendThread finish")
	})

	success := 0
	failed := 0
	successSend := 0
	failSend := 0

	elapse := time.Duration(0)
	d := 60 * time.Second
	reportTicker := time.NewTicker(d)
	now := time.Now()
	for {
		select {
		case batch := <-ch:
			if len(batch) <= 0 {
				break
			}
			now = time.Now()
			for {
				_, err := proxy.producer.SendSync(ctx, batch...)
				if nil != err {
					failed += len(batch)
					failSend++
					time.Sleep(50 * time.Millisecond)
					logutils.Error("ProducerProxy|SendSync err", zap.Error(err), zap.Int("idx", proxy.idx), zap.String("topic", batch[0].Topic))
				} else {
					success += len(batch)
					successSend++
					break
				}
			}

			elapse += time.Since(now)

		case <-reportTicker.C:
			logutils.Info("ProducerProxy|SendThread report", zap.Int("idx", proxy.idx), zap.Int("success", success), zap.Int("succSend", successSend),
				zap.Int("failed", failed), zap.Int("failSend", failSend), zap.Duration("elapse", elapse), zap.Duration("duration", d))

			success = 0
			failed = 0
			elapse = 0
			successSend = 0
			failSend = 0
		case <-ctx.Done():
			logutils.Info("ProducerProxy|SendThread done", zap.Int("idx", proxy.idx), zap.Int("success", success), zap.Int("succSend", successSend),
				zap.Int("failed", failed), zap.Int("failSend", failSend))
			return
		}
	}
}

func (proxy *ProducerProxy) goBatchMessage(ctx context.Context, topic string, chBatch chan *primitive.Message, chSend chan []*primitive.Message, wg *sync.WaitGroup) {
	defer dumputils.HandlePanic(func() {
		logutils.Info("ProducerProxy|goBatchMessage finish")
	})

	batchCount := 0
	batchSize := 0
	batch := make([]*primitive.Message, 0, proxy.cfg.BatchCount)
	total := 0
	d := 60 * time.Second
	reportTicker := time.NewTicker(d)
	for {
		select {
		case msg, ok := <-chBatch:
			if !ok {
				wg.Done()
				return
			}
			batch = batch[0:0]
			batchCount = 0
			batchSize = 0
			batch = append(batch, msg)
			batchCount = len(chBatch)
			batchSize += len(msg.Body)
			if batchCount > 0 && batchSize < proxy.cfg.BatchSize {
				if batchCount > proxy.cfg.BatchCount-1 {
					batchCount = proxy.cfg.BatchCount - 1
				}

				for i := 0; i < batchCount; i++ {
					batch = append(batch, <-chBatch)
					batchSize += len(msg.Body)
					if batchSize >= proxy.cfg.BatchSize {
						break
					}
				}
			}

			total += len(batch)
			chSend <- batch
		case <-reportTicker.C:
			logutils.Info("ProducerProxy|goBatchMessage report", zap.String("topic", topic), zap.Int("batched", total), zap.Int("idx", proxy.idx),
				zap.Int("remain", len(chBatch)))
			total = 0
		case <-ctx.Done():
			close(chBatch)
			logutils.Info("ProducerProxy|goBatchMessage done", zap.String("topic", topic), zap.Int("batched", total), zap.Int("idx", proxy.idx),
				zap.Int("remain", len(chBatch)))
		}
	}
}

func (proxy *ProducerProxy) goDispatchThread() {
	defer dumputils.HandlePanic(func() {
		logutils.Info("ProducerProxy|goDispatchThread finish")
	})

	ctx, cancel := context.WithCancel(context.TODO())
	proxy.cancel = cancel

	chSend := make(chan []*primitive.Message)
	sendCtx, sendCancel := context.WithCancel(context.TODO())
	go proxy.goSendThread(sendCtx, chSend)

	batchCtx, batchCancel := context.WithCancel(context.TODO())

	wg := sync.WaitGroup{}
	total := 0
	d := 60 * time.Second
	reportTicker := time.NewTicker(d)
	for {
		select {
		case content, ok := <-proxy.chContent:
			if !ok {

				//关闭batch 协程
				batchCancel()
				wg.Add(len(proxy.batchChanCache))
				wg.Wait()

				//关闭send 协程
				sendCancel()
				logutils.Info("ProducerProxy|goDispatchThread close channel", zap.Int("idx", proxy.idx), zap.Int("remain", len(proxy.chContent)))
				return
			}
			msg := formatContent(content)

			batchChan := proxy.getBatchChan(msg.Topic, batchCtx, chSend, &wg)
			batchChan <- msg
			total++
		case <-reportTicker.C:
			logutils.Info("ProducerProxy|goDispatchThread report", zap.Int("total", total), zap.Int("idx", proxy.idx), zap.Int("remain", len(proxy.chContent)))
			total = 0
		case <-ctx.Done():
			//关闭管道
			close(proxy.chContent)
			logutils.Info("ProducerProxy|goDispatchThread done", zap.Int("total", total), zap.Int("idx", proxy.idx), zap.Int("remain", len(proxy.chContent)))
		}
	}
}

func (proxy *ProducerProxy) getBatchChan(topic string, ctx context.Context, chSend chan []*primitive.Message, wg *sync.WaitGroup) chan *primitive.Message {
	ret, ok := proxy.batchChanCache[topic]
	if ok {
		return ret
	}
	ret = make(chan *primitive.Message, 1000)

	//打包协程
	go proxy.goBatchMessage(ctx, topic, ret, chSend, wg)
	proxy.batchChanCache[topic] = ret
	return ret
}

func formatContent(content *PublishContent) *primitive.Message {
	msg := &primitive.Message{
		Topic: content.Topic,
		Body:  content.Body,
	}
	if "" != content.Tag {
		msg.WithTag(content.Tag)
	}
	if len(content.Keys) > 0 {
		msg.WithKeys(content.Keys)
	}

	return msg
}
