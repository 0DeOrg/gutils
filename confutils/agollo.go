package confutils

/**
 * @Author: lee
 * @Description:
 * @File: appolo
 * @Date: 2022/2/18 5:37 下午
 */

import (
	"bytes"
	"fmt"
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/spf13/viper"
	"log"
	"reflect"
	"strings"
)

const (
	propertiesType   = "properties"
	propertiesFormat = "%s:%v\n" // NOTE: 请勿做任何修改, 否则极有可能引发服务雪崩!
)

type CommonProperties struct {
	properties []byte
}

// ApolloGetServerConfig
/* @Description: Apollo配置动态加载（服务整体配置信息加载专用）
 * @param c *config.AppConfig
 * @param confStruct interface{}
 * @return error
 */
func ApolloGetServerConfig(c *config.AppConfig, confStruct interface{}) error {
	propertiesStr, err := loadServerConfigFromApollo(c)
	if err != nil {
		return fmt.Errorf("LoadAllConfigFromApollo error: %s", err.Error())
	}

	// 创建properties格式配置信息处理对象
	properties := newProperties(propertiesStr)
	properties.setString(propertiesStr)

	// FIXME: 跑稳定之后删掉
	//fmt.Println("解析到的properties格式的配置信息:")
	//fmt.Println(string(properties.properties))

	if string(properties.properties) == "" {
		return fmt.Errorf("Agollo loadServerConfigFromApollo error, properties content is nil!")
	}

	// 将配置反序列化至服务配置
	err = properties.unmarshal(confStruct)
	if err != nil {
		return fmt.Errorf("UnmarshalProperties error: %s", err.Error())
	}

	return nil
}

// @Description 根据指定的apollo应用配置访问配置信息
// @Author Oracle
// @Version 1.0
// @Update Oracle 2021-11-24 init
func loadServerConfigFromApollo(c *config.AppConfig) (string, error) {
	// 初始化apollo客户端日志器
	// agollo.SetLogger(&agollolog.DefaultLogger{})

	// 根据apollo服务信息创建相对应的客户端
	client, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return c, nil
	})

	if err != nil {
		log.Fatal("Agollo startWithConfig error: ", err)
	}
	fmt.Printf("初始化Apollo配置成功, appid: [%s], namespace: [%s]\n", c.AppID, c.NamespaceName)

	var (
		propertiesStr string
	)

	// 根据apollo的namespace配置信息获取对应的配置
	cache := client.GetConfigCache(c.NamespaceName)
	cache.Range(func(key, value interface{}) bool {
		// 判断数值的实际类型, 如果是数组类的, 需要进行特殊处理, 用以viper库做反序列化
		// NOTE: 如果没有问题请不要修改
		if value == nil {
			fmt.Println("**************Nil value, key: ", key)
			return true
		}
		t := reflect.TypeOf(value)
		switch t.Kind() {
		case reflect.Slice:
			// 保存slice的value字符串
			valueSlice := fmt.Sprintf("%v\n", value)

			// 干掉字符串两遍的[], 拿到原串
			replacer := strings.NewReplacer("[", "", "]", "")
			sliceStr := replacer.Replace(valueSlice)

			// 对原串进行拆分解析, 按照逗号,分隔符重新组装新串
			elements := strings.Split(sliceStr, " ")

			// 生成最终由逗号(,)分割的字符串
			finalValue := strings.Join(elements, ",")

			// 将逗号分隔符组装好的字符串保存, 方便viper做反序列化
			propertiesStr += fmt.Sprintf(propertiesFormat, key, finalValue)

		default:
			// 非数组类的按照统一规则进行处理
			propertiesStr += fmt.Sprintf(propertiesFormat, key, value)
		}
		// fmt.Printf("变量类型为: %T, 其值为: %+v\n", value, value)

		return true
	})

	return propertiesStr, nil
}

// @Description 创建一个properties格式配置信息处理器对象
// @Author Oracle
// @Version 1.0
// @Update Oracle 2021-11-24 init
func newProperties(propertiesStr string) *CommonProperties {
	cp := &CommonProperties{}
	return cp
}

// @Description 根据string格式内容设置数据
// @Author Oracle
// @Version 1.0
// @Update Oracle 2021-11-24 init
func (cp *CommonProperties) setString(propertiesStr string) {
	cp.properties = []byte(propertiesStr)

	return
}

// @Description 将apollo配置信息反序列化至服务配置管理结构体对象当中
// @Author Oracle
// @Version 1.0
// @Update Oracle 2021-11-24 init
func (cp *CommonProperties) unmarshal(confStruct interface{}) error {
	// 初始化viper对象
	myViper := viper.New()

	// 设置待解析配置信息的数据格式类型
	myViper.SetConfigType(propertiesType)

	// 加载配置信息至viper对象
	err := myViper.ReadConfig(bytes.NewBuffer(cp.properties))
	if err != nil {
		// fmt.Println("Viper read config error: ", err)
		return err
	}

	// 将配置信息反序列化至配置管理结构体对象当中
	err = myViper.Unmarshal(confStruct)
	if err != nil {
		// fmt.Println("Viper unmarshal properties config to config structure error: ", err)
		return err
	}

	return nil
}
