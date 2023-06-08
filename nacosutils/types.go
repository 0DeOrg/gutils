package nacosutils

/**
 * @Author: lee
 * @Description:
 * @File: types
 * @Date: 2023-05-22 8:20 下午
 */

const NACOS_PATH = "apollo.yaml"

type NacosConfig struct {
	Servers []*ServerConfig
}

type ServerConfig struct {
	IpAddr   string `mapstructure:"ip-addr"         json:"ip-addr"          yaml:"ip-addr"`     // the nacos server address
	Port     uint64 `mapstructure:"port"            json:"port"          yaml:"port"`           // nacos server port
	GrpcPort uint64 `mapstructure:"grpc-port"       json:"grpc-port"          yaml:"grpc-port"` // nacos server grpc port, default=server port + 1000, this is not required
}

type ClientConfig struct {
}
