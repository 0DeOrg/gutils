package pprofutils

/**
 * @Author: lee
 * @Description:
 * @File: types
 * @Date: 2023-01-16 1:59 下午
 */

type PProfConfig struct {
	Enable bool   `json:"enable"     yaml:"enable"   mapstructure:"enable"`
	Host   string `json:"host"     yaml:"host"   mapstructure:"host"`
	Port   int    `json:"port"     yaml:"port"   mapstructure:"port"`
}
