package gutils

/**
 * @Author: lee
 * @Description:
 * @File: viper
 * @Date: 2021/9/14 11:11 上午
 */

import (
	"fmt"
	"github.com/0DeOrg/gutils/judge"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

const CONFIG_PATH = "config.yaml"

//外部命令行解析的时候赋值
var CfgPathFlag = ""

func NewViper(path string, pObj interface{}, callback ...func()) (*viper.Viper, error) {
	var config string
	if len(path) == 0 {
		//flag.StringVar(&config, "c", "", "choose config file.")
		//flag.Parse()
		if CfgPathFlag != "" {
			fmt.Printf("您正在使用命令行的-c参数传递的值,config的路径为%v\n", config)
		} else {
			config = CONFIG_PATH
		}
	} else {
		config = path
		fmt.Printf("您正在使用func NewViper()传递的值,config的路径为%v\n", config)
	}
	v := viper.New()
	v.SetConfigFile(config)
	err := v.ReadInConfig()
	if err != nil {
		log.Fatal("fatal error can't find config file: ", err.Error())
	}
	v.WatchConfig()

	if !judge.IsStructPtr(pObj) {
		log.Fatal("fatal error viper pObj must be struct pointer")
	}

	v.OnConfigChange(func(e fsnotify.Event) {
		log.Println("config file changed:", e.Name)
		if err := v.Unmarshal(pObj); err != nil {
			log.Println(err.Error())
		}

		for _, callFunc := range callback {
			go callFunc()
		}

	})

	if err := v.Unmarshal(pObj); err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return v, nil
}
