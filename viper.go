package gutils
/**
 * @Author: lee
 * @Description:
 * @File: viper
 * @Date: 2021/9/14 11:11 上午
 */

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"gutils/judge"
)

const CONFIG_PATH = "config.yaml"

func NewViper(path string, pObj interface{}) *viper.Viper {
	var config string
	if len(path) == 0 {
		flag.StringVar(&config, "c", "", "choose config file.")
		flag.Parse()
		if config != "" {
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
		panic(fmt.Errorf("fatal error can't find config file: %s \n", err))
	}
	v.WatchConfig()

	if !judge.IsStructPtr(pObj) {
		panic(fmt.Errorf("fatal error viper pObj must be struct pointer"))
	}

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := v.Unmarshal(pObj); err != nil {
			fmt.Println(err)
		}
	})

	if err := v.Unmarshal(pObj); err != nil {
		fmt.Println(err)
	}
	return v
}