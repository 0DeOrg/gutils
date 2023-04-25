package logutils

/**
 * @Author: lee
 * @Description:
 * @File: zap
 * @Date: 2021/9/13 6:04 下午
 */

import (
	"fmt"
	"github.com/0DeOrg/gutils/fileutils"
	"github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"time"
)

type ZapConfig struct {
	Directory    string `json:"directory"     yaml:"directory"   mapstructure:"directory"`
	ShowLine     bool   `json:"show-line"     yaml:"show-line"   mapstructure:"show-line"`
	ZapLevel     string `json:"zap-level"     yaml:"zap-level"   mapstructure:"zap-level"`
	Archive      string `json:"archive"     yaml:"archive"       mapstructure:"archive"`
	Format       string `json:"format"     yaml:"format"         mapstructure:"format"`
	LinkName     string `json:"link-name"     yaml:"link-name"   mapstructure:"link-name"`
	LogInConsole bool   `json:"log-in-console"     yaml:"log-in-console"    mapstructure:"log-in-console"`
	EncodeLevel  string `json:"encode-level"     yaml:"encode-level"        mapstructure:"encode-level"`
}

var DefaultZapConfig = ZapConfig{
	Directory:    "log",
	ZapLevel:     "info",
	Archive:      "log",
	Format:       "console",
	LinkName:     "log/latest_log",
	LogInConsole: true,
	EncodeLevel:  "LowercaseColorLevelEncoder",
	ShowLine:     true,
}

type ZapLogModule struct {
	logger *zap.Logger
	config ZapConfig
}

var _ ILogger = (*ZapLogModule)(nil)

func (m *ZapLogModule) Info(msg string, fields ...zap.Field) {
	m.logger.Info(msg, fields...)
}

func (m *ZapLogModule) Infof(format string, vals ...interface{}) {

}

func (m *ZapLogModule) Error(msg string, fields ...zap.Field) {
	m.logger.Error(msg, fields...)
}

func (m *ZapLogModule) Errorf(format string, vals ...interface{}) {

}

func (m *ZapLogModule) Warn(msg string, fields ...zap.Field) {
	m.logger.Warn(msg, fields...)
}

func (m *ZapLogModule) Warnf(format string, vals ...interface{}) {

}

func (m *ZapLogModule) Debug(msg string, fields ...zap.Field) {
	m.logger.Debug(msg, fields...)
}

func (m *ZapLogModule) Debugf(format string, vals ...interface{}) {

}

func (m *ZapLogModule) Fatal(msg string, fields ...zap.Field) {
	m.logger.Fatal(msg, fields...)
}

func (m *ZapLogModule) Fatalf(format string, vals ...interface{}) {

}
func (m *ZapLogModule) DPanic(msg string, fields ...zap.Field) {
	m.logger.DPanic(msg, fields...)
}

func (m *ZapLogModule) Panic(msg string, fields ...zap.Field) {
	m.logger.Panic(msg, fields...)
}

var zapConfig ZapConfig
var level zapcore.Level

func newZapLogModule(config ZapConfig) (*ZapLogModule, error) {
	logger, err := newZapLogger(config)
	if nil != err {
		return nil, err
	}
	ret := ZapLogModule{
		logger: logger,
		config: config,
	}

	return &ret, nil
}

func newZapLogger(config ZapConfig) (logger *zap.Logger, err error) {
	zapConfig = config
	if err = fileutils.CreateDirectoryIfNotExist(zapConfig.Directory, os.ModePerm); nil != err {
		return nil, err
	}

	// 初始化配置文件的Level
	switch zapConfig.ZapLevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "dPanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}

	if level == zap.DebugLevel || level == zap.ErrorLevel {
		logger = zap.New(getEncoderCore(), zap.AddStacktrace(level))
	} else {
		logger = zap.New(getEncoderCore())
	}
	if zapConfig.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
		logger = logger.WithOptions(zap.AddCallerSkip(2))
	}

	return logger, nil
}

// getEncoderConfig 获取zapcore.EncoderConfig
func getEncoderConfig() (config zapcore.EncoderConfig) {
	config = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	switch {
	case zapConfig.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	case zapConfig.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		config.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case zapConfig.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
		config.EncodeLevel = zapcore.CapitalLevelEncoder
	case zapConfig.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	return config
}

// getEncoder 获取zapcore.Encoder
func getEncoder() zapcore.Encoder {
	if "json" == zapConfig.Format {
		return zapcore.NewJSONEncoder(getEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(getEncoderConfig())
}
func getEncoderCore() (core zapcore.Core) {
	writer, err := GetWriteSyncer() // 使用file-rotatelogs进行日志分割
	if err != nil {
		fmt.Printf("Get Write Syncer Failed err:%v", err.Error())
		return
	}
	return zapcore.NewCore(getEncoder(), writer, level)
}

// 自定义日志输出时间格式
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// zap logger中加入file-rotatelogs
func GetWriteSyncer() (zapcore.WriteSyncer, error) {
	var filePath string
	filePath = path.Join(zapConfig.Directory, zapConfig.Archive+"-%Y-%m-%d.log")

	var linkName rotatelogs.Option
	if zapConfig.Archive == "" {
		linkName = rotatelogs.WithLinkName(zapConfig.LinkName)
	} else {
		linkName = rotatelogs.WithLinkName(zapConfig.LinkName + "_" + zapConfig.Archive)
	}

	fileWriter, err := rotatelogs.New(
		filePath,
		linkName,
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if zapConfig.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter)), err
	}
	return zapcore.AddSync(fileWriter), err
}
