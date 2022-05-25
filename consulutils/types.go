package consulutils

/**
 * @Author: lee
 * @Description:
 * @File: consul_define
 * @Date: 2021/9/2 6:06 下午
 */
import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"strings"
)

const (
	ConsulSecure = "secure"
)

type ConsulConfig struct {
	Address    string `mapstructure:"address"        json:"address"       yaml:"address"`
	Port       int    `mapstructure:"port"           json:"port"          yaml:"port"`
	Name       string `mapstructure:"name"           json:"name"          yaml:"name"`
	HealthPath string `mapstructure:"health-path"     json:"health-path"    yaml:"health-path"`
}

type ConsulFindCfg struct {
	Name   string `mapstructure:"name"           json:"name"          yaml:"name"`
	Secure bool   `mapstructure:"secure"           json:"secure"          yaml:"secure"`
}

type ServiceInstance struct {
	InstanceId string
	Name       string
	Host       string
	Port       int
	Secure     bool
	Metadata   map[string]string
	HealthUrl  string
}

type consulServiceRegistry struct {
	Client               *consulapi.Client
	LocalServiceInstance *ServiceInstance
}

func (c *consulServiceRegistry) register(instance *ServiceInstance) error {
	reg := &consulapi.AgentServiceRegistration{}
	reg.ID = instance.InstanceId
	reg.Name = instance.Name
	reg.Address = instance.Host
	reg.Port = instance.Port
	tags := make([]string, 0)
	if instance.Secure {
		tags = append(tags, ConsulSecure+"=true")
	} else {
		tags = append(tags, ConsulSecure+"=false")
	}

	if nil != instance.Metadata {
		reg.Meta = instance.Metadata
	}

	reg.Tags = tags
	//reg.Meta = instance.Metadata

	check := &consulapi.AgentServiceCheck{}
	check.Timeout = "5s"
	check.Interval = "5s"
	check.DeregisterCriticalServiceAfter = "30s"
	scheme := "http"
	if instance.Secure {
		scheme = "https"
	}
	check.HTTP = fmt.Sprintf("%s://%s:%d%s", scheme, instance.Host, instance.Port, instance.HealthUrl)
	reg.Check = check

	err := c.Client.Agent().ServiceRegister(reg)
	if nil != err {
		return err
	}

	c.LocalServiceInstance = instance

	return nil
}

func (c *consulServiceRegistry) unregister(instance ServiceInstance) error {
	if nil == c.Client {
		return fmt.Errorf("service has not register")
	}
	err := c.Client.Agent().ServiceDeregister(instance.InstanceId)
	if err != nil {
		return err
	}
	return nil
}

// getInstanceByName
/* @Description: 查找consul上面的服务，可能多个服务名字一样
 * @param name string
 * @return []*ServiceInstance
 * @return error
 */
func (c *consulServiceRegistry) getInstanceByName(name string) ([]*ServiceInstance, error) {
	serviceEntry, _, err := c.Client.Health().Service(name, "", true, nil)
	if nil != err {
		return nil, err
	}

	ret := make([]*ServiceInstance, 0)
	for _, v := range serviceEntry {

		isSecure := CheckServiceSecure(v.Service)
		s := &ServiceInstance{
			InstanceId: v.Service.ID,
			Name:       v.Service.Namespace,
			Host:       v.Service.Address,
			Port:       v.Service.Port,
			Metadata:   v.Service.Meta,
			Secure:     isSecure,
		}
		ret = append(ret, s)
	}

	return ret, nil
}

// getServicesOnAgent
/* @Description: 获取当前注册node上面的服务，不是datacenter上面的所有服务，仅仅当前节点上面的所有服务
 * @return []*ServiceInstance
 * @return error
 */
func (c *consulServiceRegistry) getServicesOnAgent() ([]*ServiceInstance, error) {
	services, err := c.Client.Agent().Services()
	if nil != err {
		return nil, err
	}

	ret := make([]*ServiceInstance, 0)

	for _, v := range services {
		s := &ServiceInstance{
			InstanceId: v.ID,
			Name:       v.Service,
			Host:       v.Address,
			Port:       v.Port,
			Metadata:   v.Meta,
		}

		ret = append(ret, s)
	}

	return ret, nil
}
func CheckServiceSecure(service *consulapi.AgentService) bool {
	tags := service.Tags
	for _, v := range tags {
		list := strings.Split(v, "=")
		if len(list) == 2 && "secure" == list[0] {
			isSecure := list[1]
			if "true" == strings.ToLower(isSecure) {
				return true
			} else {
				return false
			}
		}
	}

	for key, value := range service.Meta {
		if "secure" == key {
			if "true" == strings.ToLower(value) {
				return true
			} else {
				return false
			}
		}
	}

	return false
}
