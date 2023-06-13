package nacosutils

/**
 * @Author: lee
 * @Description:
 * @File: types
 * @Date: 2023-05-22 8:20 下午
 */

const NACOS_PATH = "nacos.yaml"

type NacosConfig struct {
	Servers []*ServerConfig `mapstructure:"servers"            json:"servers"          yaml:"servers"`
	Client  *ClientConfig   `mapstructure:"client"            json:"client"          yaml:"client"`
}

type ServerConfig struct {
	IpAddr   string `mapstructure:"ip-addr"         json:"ip-addr"          yaml:"ip-addr"`     // the nacos server address
	Port     uint64 `mapstructure:"port"            json:"port"          yaml:"port"`           // nacos server port
	GrpcPort uint64 `mapstructure:"grpc-port"       json:"grpc-port"          yaml:"grpc-port"` // nacos server grpc port, default=server port + 1000, this is not required
}

type ClientConfig struct {
	UserName  string `mapstructure:"user-name"            json:"user-name"          yaml:"user-name"`
	Password  string `mapstructure:"password"            json:"password"          yaml:"password"`
	Namespace string `mapstructure:"namespace"            json:"namespace"          yaml:"namespace"`
	DataId    string `mapstructure:"data-id"            json:"data-id"          yaml:"data-id"`
	Group     string `mapstructure:"group"            json:"group"          yaml:"group"`
	LogPath   string `mapstructure:"log-path"            json:"log-path"          yaml:"log-path"`
	CachePath string `mapstructure:"cache-path"            json:"cache-path"          yaml:"cache-path"`
	LogLevel  string `mapstructure:"log-level"            json:"log-level"          yaml:"log-level"`
}
