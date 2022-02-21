package confutils

/**
 * @Author: lee
 * @Description:
 * @File: load_server_conf_file
 * @Date: 2022/2/18 5:39 下午
 */

import (
	"github.com/apolloconfig/agollo/v4/env/config"
	"gutils"
)

// LoadApolloConfigWhenStartup
/* @Description:
 * @param confPath string
 * @param confFile string
 * @param runEnv string
 * @return *config.AppConfig
 */
func LoadApolloConfigWhenStartup(path string) *config.AppConfig {
	// 通过viper解析配置文件并反序列化至结构体
	apolloConf := Apollo{}
	fileName := APOLLO_PATH
	if "" == path {
		fileName = path
	}
	gutils.NewViper(fileName, &apolloConf)

	// 返回agollo库使用的配置结构
	apolloAppConf := &config.AppConfig{
		AppID:          apolloConf.Apollo.AppID,
		Cluster:        apolloConf.Apollo.Cluster,
		NamespaceName:  apolloConf.Apollo.NamespaceName,
		IP:             apolloConf.Apollo.IP,
		IsBackupConfig: apolloConf.Apollo.IsBackupConfig,
		Secret:         apolloConf.Apollo.Secret,
	}

	return apolloAppConf
}
