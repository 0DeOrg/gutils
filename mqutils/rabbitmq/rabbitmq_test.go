package rabbitmq

/**
 * @Author: lee
 * @Description:
 * @File: rabbitmq_test
 * @Date: 2022/2/22 10:56 上午
 */
import (
	"github.com/0DeOrg/gutils/logutils"
	"go.uber.org/zap"
	"log"
	"testing"
)

func Test_Publish(t *testing.T) {
	cfg := &RabbitMQConfig{
		User:      "admin",
		Password:  "QGgRlqdMXdGu",
		Addresses: []string{"192.168.10.45:5672"},
		VHost:     "/",
	}
	mq, err := NewRabbitMq(cfg, false)
	if nil != err {
		log.Fatal(err.Error())
	}
	err = mq.ExchangeDeclare("test.kline", ExchangeFanout, true)
	if nil != err {
		log.Fatal(err.Error())
	}

	//body := `{"key": "test123"}`
	//reqBody, _ := json.Marshal(body)
	//content := &PublishContent{
	//	ExchangeName: "test",
	//	Content:      reqBody,
	//	RoutingKey:   "fanout no need",
	//}
	//confirmed, err := mq.Publish(content, false)
	//if nil != err {
	//	log.Fatal(err.Error())
	//}
	//
	//log.Println("confirmed:", confirmed)
}

func Test_Consumer(t *testing.T) {
	logutils.InitLogger(logutils.DefaultZapConfig)
	cfg := &RabbitMQConfig{
		User:      "admin",
		Password:  "QGgRlqdMXdGu",
		Addresses: []string{"192.168.10.45:5672"},
		VHost:     "/",
	}

	mq, err := NewRabbitMq(cfg, false)
	if nil != err {
		logutils.Fatal("NewRabbitMq:", zap.Error(err))
	}

	err = mq.QueueDeclare("test.index.quotation", true, true, false, false)
	if nil != err {
		logutils.Fatal("QueueDeclare:", zap.Error(err))
	}
	queueName := "test.index.quotation"
	err = mq.QueueBind(queueName, "qt-property.index.quotation", "#")
	if nil != err {
		logutils.Fatal("QueueBind:", zap.Error(err))
	}

	msgCh, err := mq.Consume(queueName)
	if nil != err {
		logutils.Fatal("Consume: ", zap.Error(err))
	}

	for {
		select {
		case msg := <-msgCh:
			logutils.Info("msg: ", zap.String("msg", string(msg.Body)))
		}
	}
}
