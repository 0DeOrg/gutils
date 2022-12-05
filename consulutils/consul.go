package consulutils

/**
 * @Author: lee
 * @Description:
 * @File: consul
 * @Date: 2021/9/2 6:05 下午
 */

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/0DeOrg/gutils/network"
	"strconv"
	"strings"
)

var (
	gConsulRegistry *consulServiceRegistry = nil
)

func RegisterConsul(consulHost string, consulPort int, token string, instance *ServiceInstance) error {
	var err error
	if nil == gConsulRegistry {
		gConsulRegistry, err = NewConsulServiceRegistry(consulHost, consulPort, token)
		if nil != err {
			return fmt.Errorf("create consul service registry fatal, err: %s", err.Error())
		}
	}

	err = gConsulRegistry.register(instance)
	if nil != err {
		return fmt.Errorf("register instance fatal, err: %s", err.Error())
	}

	return nil
}

func DeregisterSelf() {

}

// NewLocalServiceInstance
/* @Description: 构造本地service 实例
 * @param name string
 * @param host string
 * @param port int
 * @param secure bool
 * @param metaData map[string]string
 * @param instanceId string
 * @param health string
 * @return *ServiceInstance
 */
func NewLocalServiceInstance(name string, host string, port int, secure bool, metaData map[string]string, instanceId string, health string) *ServiceInstance {
	if "" == host {
		host = network.GetLocalIP()
	}

	if "" == instanceId {
		sections := strings.Split(host, ".")
		instanceId = name
		for _, v := range sections {
			instanceId += "-"
			instanceId += v
		}

		instanceId += "-" + strconv.Itoa(port)
	}

	ret := ServiceInstance{
		InstanceId: instanceId,
		Name:       name,
		Host:       host,
		Port:       port,
		Secure:     secure,
		Metadata:   metaData,
		HealthUrl:  health,
	}

	return &ret
}

// NewConsulServiceRegistry
/* @Description: 构造一个注册器
 * @param host string
 * @param port int
 * @param token string
 * @return *consulServiceRegistry
 * @return error
 */
func NewConsulServiceRegistry(host string, port int, token string) (*consulServiceRegistry, error) {
	config := consulapi.DefaultConfig()
	config.Address = host + ":" + strconv.Itoa(port)
	config.Token = token
	client, err := consulapi.NewClient(config)
	if nil != err {
		return nil, err
	}

	return &consulServiceRegistry{Client: client}, nil
}

// FindServiceByName
/* @Description: 根据名字获取服务，可能有多个可根据实例名区分
 * @param name string
 * @return []*ServiceInstance
 * @return error
 */
func FindServiceByName(name string) ([]*ServiceInstance, error) {
	if nil == gConsulRegistry {
		return nil, fmt.Errorf("has not regiseter consul")
	}

	return gConsulRegistry.getInstanceByName(name)

}

func GetServicesOnAgent() ([]*ServiceInstance, error) {
	if nil == gConsulRegistry {
		return nil, fmt.Errorf("has not regiseter consul")
	}

	return gConsulRegistry.getServicesOnAgent()
}
