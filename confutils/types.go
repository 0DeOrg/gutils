package confutils

/**
 * @Author: lee
 * @Description:
 * @File: types
 * @Date: 2022/2/18 5:40 下午
 */

const APOLLO_PATH = "apollo.yaml"

// Apollo 配置文件
type Apollo struct {
	Apollo ApolloConfig `mapstructure:"apollo"          json:"apollo"          yaml:"apollo"`
}

// Apollo 配置文件详细信息
type ApolloConfig struct {
	AppID          string `mapstructure:"appid"                   json:"appid"                   yaml:"appid"`
	Cluster        string `mapstructure:"cluster"                 json:"cluster"                 yaml:"cluster"`
	IP             string `mapstructure:"ip"                      json:"ip"                      yaml:"ip"`
	NamespaceName  string `mapstructure:"namespace-name"          json:"namespace-name"          yaml:"namespace-name"`
	IsBackupConfig bool   `mapstructure:"is-backup-config"        json:"is-backup-config"        yaml:"is-backup-config"`
	Secret         string `mapstructure:"secret"                  json:"secret"                  yaml:"secret"`
	// BackupConfigPath        string               `mapstructure:"backupConfigPath"        json:"backupConfigPath"        yaml:"backupConfigPath"`
	// SyncServerTimeout       int                  `mapstructure:"syncServerTimeout"       json:"syncServerTimeout"       yaml:"syncServerTimeout"`
	// notificationsMap        *notificationsMap    `mapstructure:"notificationsMap"        json:"notificationsMap"        yaml:"notificationsMap"`
	// currentConnApolloConfig *CurrentApolloConfig `mapstructure:"currentConnApolloConfig" json:"currentConnApolloConfig" yaml:"currentConnApolloConfig"`
}
