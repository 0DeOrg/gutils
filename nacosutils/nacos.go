package nacosutils

/**
 * @Author: lee
 * @Description:
 * @File: nacos
 * @Date: 2023-05-22 8:20 下午
 */
import (
	"bytes"
	"fmt"
	"github.com/0DeOrg/gutils"
	"github.com/0DeOrg/gutils/convert"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"log"
)

// NacosGetConfig
/* @Description: 通过nacos 获取配置
 * @param path string nacos 系统配置文件路径
 * @param confPtr interface{} 导出配置结构体，必须传指针
 * @param confType string 配置文件类型，yaml， 等
 * @return error
 */

func ReadNacosConfig(path string) (*NacosConfig, error) {
	nacosConf := NacosConfig{}
	fileName := NACOS_PATH
	if "" != path {
		fileName = path
	}

	_, err := gutils.NewViper(fileName, &nacosConf)
	if nil != err {
		return nil, fmt.Errorf("nacos.yaml load err: %s", err.Error())
	}

	return &nacosConf, nil
}
func NacosGetConfig(nacosConf *NacosConfig, confPtr interface{}, confType string) error {
	if err := convert.MustBeStructPtr(confPtr); nil != err {
		return err
	}

	scs := make([]constant.ServerConfig, 0)
	for _, server := range nacosConf.Servers {
		sc := constant.ServerConfig{
			IpAddr:   server.IpAddr,
			Port:     server.Port,
			GrpcPort: server.GrpcPort,
		}

		scs = append(scs, sc)
	}
	cc := constant.ClientConfig{
		Username:            nacosConf.Client.UserName,
		Password:            nacosConf.Client.Password,
		NamespaceId:         nacosConf.Client.Namespace, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		NotLoadCacheAtStart: true,
		LogDir:              nacosConf.Client.LogPath,
		CacheDir:            nacosConf.Client.CachePath,
		//RotateTime:          "1h",
		//MaxAge:              3,
		LogLevel: nacosConf.Client.LogLevel,
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		constant.KEY_SERVER_CONFIGS: scs,
		constant.KEY_CLIENT_CONFIG:  cc,
	})

	if nil != err {
		return fmt.Errorf("CreateConfigClient err: %s", err.Error())
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: nacosConf.Client.DataId,
		Group:  nacosConf.Client.Group,
	})

	if err != nil {
		return fmt.Errorf("load nacos config err:%s", err.Error())
	}

	if "" == content {
		return fmt.Errorf("nacos content is empty")
	}

	vp := viper.New()
	vp.SetConfigType(confType)
	listener := changeListener{
		vp:       vp,
		confPtr:  confPtr,
		confType: confType,
	}

	// 填充到配置文件
	err = listener.OnChange(nacosConf.Client.Namespace, nacosConf.Client.Group, nacosConf.Client.DataId, content)
	if nil != err {
		return fmt.Errorf("viper read nacos config err: %s", err.Error())
	}

	// 注册监听
	err = configClient.ListenConfig(vo.ConfigParam{
		DataId: nacosConf.Client.DataId,
		Group:  nacosConf.Client.Group,
		OnChange: func(namespace, group, dataId, content string) {
			err = listener.OnChange(namespace, group, dataId, content)
			if nil != err {
				log.Println("nacos config OnChange", "err", err.Error())
			} else {
				log.Println("nacos config OnChange success", "content", content)
			}
		},
	})

	return nil
}

type changeListener struct {
	vp       *viper.Viper
	confPtr  interface{}
	confType string
}

func (l *changeListener) OnChange(namespace, group, dataId, content string) error {
	vp := viper.New()
	vp.SetConfigType(l.confType)

	err := vp.ReadConfig(bytes.NewReader([]byte(content)))
	if nil != err {
		return fmt.Errorf("viper read config err: %s", err.Error())
	}

	pairs := vp.AllSettings()

	for key, value := range pairs {
		l.vp.Set(key, value)
	}

	err = vp.Unmarshal(l.confPtr)
	if nil != err {
		return fmt.Errorf("unmarshal viper err: %s", err.Error())
	}

	return nil
}
