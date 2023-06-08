package nacosutils

import (
	"github.com/0DeOrg/gutils"
	"log"
)

/**
 * @Author: lee
 * @Description:
 * @File: load_conf
 * @Date: 2023-05-22 8:49 下午
 */
func LoadNacosConfigWhenStart(filePath string) (*NacosConfig, error) {
	ret := &NacosConfig{}
	fileName := NACOS_PATH
	if "" != filePath {
		fileName = filePath
	}
	_, err := gutils.NewViper(fileName, ret)
	if nil != err {
		log.Fatal("read nacos config fatal", "err", err.Error())
	}

	return ret, nil
}
