package rocketmq

/**
 * @Author: lee
 * @Description:
 * @File: types
 * @Date: 2023-04-26 2:51 下午
 */

type RocketMQConfig struct {
	NameServers   []string `json:"name-servers"     yaml:"name-servers"   mapstructure:"name-servers"`
	ProducerCount int      `json:"producer-count"     yaml:"producer-count"   mapstructure:"producer-count"`
	ProducerGroup string   `json:"producer-group"     yaml:"producer-group"   mapstructure:"producer-group"`
}

type PublishContent struct {
	Topic string
	Tag   string
	Body  []byte
}
