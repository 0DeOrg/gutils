package nacosutils

/**
 * @Author: lee
 * @Description:
 * @File: nacos
 * @Date: 2023-05-22 8:20 下午
 */
import (
	"fmt"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
)

func NacosGetConfig(cfg *NacosConfig) error {
	scs := make([]*constant.ServerConfig, 0)
	for _, server := range cfg.Servers {
		sc := &constant.ServerConfig{
			IpAddr:   server.IpAddr,
			Port:     server.Port,
			GrpcPort: server.GrpcPort,
		}

		scs = append(scs, sc)
	}
	cc := constant.ClientConfig{
		NamespaceId:         Bootstrap.Nacos.Namespace.Id, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "logs",
		CacheDir:            "cache",
		//RotateTime:          "1h",
		//MaxAge:              3,
		LogLevel: "debug",
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": scs,
		"clientConfig":  cc,
	})

	if nil != err {
		return fmt.Errorf("CreateConfigClient err: %s", err.Error())
	}
}
