package daoutils
/**
 * @Author: lee
 * @Description:
 * @File: mysql
 * @Date: 2021/9/15 3:22 下午
 */
import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
)



type MySQLCfg struct {
	Params       string `mapstructure:"params"         json:"params"        yaml:"params"`
	MaxIdleConns int    `mapstructure:"max-idle-conns" json:"maxIdleConns"  yaml:"max-idle-conns"`
	MaxOpenConns int    `mapstructure:"max-open-conns" json:"maxOpenConns"  yaml:"max-open-conns"`
	LogMode      string   `mapstructure:"log-mode"       json:"logMode"       yaml:"log-mode"`
	Prefix       string `mapstructure:"prefix"         json:"prefix"        yaml:"prefix"`
	URL          string `mapstructure:"url"           json:"url"             yaml:"url"`
	Dbname       string `mapstructure:"db-name"        json:"dbname"        yaml:"db-name"`
	Username     string `mapstructure:"username"       json:"username"      yaml:"username"`
	Password     string `mapstructure:"password"       json:"password"      yaml:"password"`
	DefaultStringSize uint `mapstructure:"default-str-size"       json:"defaultStrSize"      yaml:"default-str-size"`
}

type MySQLClient struct {
	cfg MySQLCfg
	DB *gorm.DB
}

func NewMysqlClient(cfg MySQLCfg) *MySQLClient {
	ret := MySQLClient{
		cfg: cfg,
	}

	return &ret
}

var _ IDaoClient = (*MySQLClient)(nil)
func (c *MySQLClient) Connect () error {
	dsn := c.DSN()
	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         c.cfg.DefaultStringSize,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}

	gormConfig := c.generateGormConfig()

	if db, err := gorm.Open(mysql.New(mysqlConfig), gormConfig); err != nil {
		return err
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(c.cfg.MaxIdleConns)
		sqlDB.SetMaxOpenConns(c.cfg.MaxOpenConns)
		c.DB = db
	}

	return nil
}

func (c MySQLClient) DSN() string {
	dsn := c.cfg.Username + ":" + c.cfg.Password + "@tcp(" + c.cfg.URL + ")/" + c.cfg.Dbname + "?" + c.cfg.Params
	return dsn
}

func (c MySQLClient) generateGormConfig() *gorm.Config{
	gormConfig := gorm.Config{
	}

	switch strings.ToLower( c.cfg.LogMode ) {
	case "silent":

		break
	case "info":

		break
	case "warn":

		break
	case "error":

		break
	default:

		break
	}

	return &gormConfig
}