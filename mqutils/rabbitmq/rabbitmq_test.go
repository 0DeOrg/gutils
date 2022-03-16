package rabbitmq

/**
 * @Author: lee
 * @Description:
 * @File: rabbitmq_test
 * @Date: 2022/2/22 10:56 上午
 */
import (
	"encoding/json"
	"log"
	"testing"
)

func Test_Publish(t *testing.T) {
	mq, err := NewRabbitMq("lee", "956443", "localhost:5672", "/")
	if nil != err {
		log.Fatal(err.Error())
	}
	err = mq.ExchangeDeclare("test", ExchangeFanout, true)
	if nil != err {
		log.Fatal(err.Error())
	}

	body := `{"key": "test123"}`
	reqBody, _ := json.Marshal(body)
	content := &PublishContent{
		ExchangeName: "test",
		Content:      reqBody,
		RoutingKey:   "fanout no need",
	}
	confirmed, err := mq.Publish(content, false)
	if nil != err {
		log.Fatal(err.Error())
	}

	log.Println("confirmed:", confirmed)
}