package confutils

/**
 * @Author: lee
 * @Description:
 * @File: appolo
 * @Date: 2022/2/18 5:37 下午
 */

import (
	"fmt"
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/apolloconfig/agollo/v4/storage"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"gutils/convert"
	"log"
)

type changeListener struct {
	confSt interface{}
	vp     *viper.Viper
}

// ApolloGetServerConfig
/* @Description: Apollo配置动态加载（服务整体配置信息加载专用）
 * @param c *config.AppConfig
 * @param confStruct interface{}
 * @return error
 */
func ApolloGetServerConfig(c *config.AppConfig, confStruct interface{}) error {

	if err := convert.MustBeStructPtr(confStruct); nil != err {
		return err
	}
	// 根据apollo服务信息创建相对应的客户端
	client, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return c, nil
	})

	if nil != err {
		return err
	}

	if nil == client {
		return fmt.Errorf("StartWithConfig client is nil")
	}

	cache := client.GetConfigCache(c.NamespaceName)
	vp := viper.New()
	listener := &changeListener{
		vp:     vp,
		confSt: confStruct,
	}

	client.AddChangeListener(listener)

	cache.Range(func(key, value interface{}) bool {
		vp.Set(cast.ToString(key), value)
		return true
	})

	err = vp.Unmarshal(confStruct)

	if nil != err {
		return err
	}

	return nil
}

func (l *changeListener) OnChange(event *storage.ChangeEvent) {
	for key, value := range event.Changes {
		l.vp.Set(key, value.NewValue)
		log.Printf("changeListener OnChange key: %s and new value is: %v", key, value)
	}

	err := l.vp.Unmarshal(l.confSt)
	if nil != err {
		log.Print("changeListener OnChange err", err.Error())
	}
}

//OnNewestChange 监控最新变更
func (l *changeListener) OnNewestChange(event *storage.FullChangeEvent) {

}
