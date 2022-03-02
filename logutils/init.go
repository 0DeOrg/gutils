package logutils

/**
 * @Author: lee
 * @Description:
 * @File: init
 * @Date: 2021/9/14 3:41 下午
 */
import (
	"fmt"
	"go.uber.org/zap"
	"log"
)

var (
	loggerModule ILogger
	logInit      = false

	errorNotInit = fmt.Errorf("log module not inited")
)

type ILogger interface {
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	DPanic(msg string, fields ...zap.Field)
	Infof(format string, values ...interface{})
	Errorf(format string, values ...interface{})
	Warnf(format string, values ...interface{})
	Debugf(format string, values ...interface{})
	Fatalf(format string, values ...interface{})
}

func InitLogger(config interface{}) {
	var err error
	if v, ok := config.(ZapConfig); ok {
		loggerModule, err = newZapLogModule(v)
		if nil != err {
			log.Fatal("zap log init fault, ", err.Error(), v)
		}

		logInit = true
	}
}

func Logger() ILogger {
	return loggerModule
}

func Info(msg string, fields ...zap.Field) {
	if !logInit {
		panic(errorNotInit)
	}
	loggerModule.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	if !logInit {
		panic(errorNotInit)
	}
	loggerModule.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	if !logInit {
		panic(errorNotInit)
	}
	loggerModule.Error(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	if !logInit {
		panic(errorNotInit)
	}
	loggerModule.Debug(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	if !logInit {
		panic(errorNotInit)
	}
	loggerModule.Fatal(msg, fields...)
}

func DPanic(msg string, fields ...zap.Field) {
	if !logInit {
		panic(errorNotInit)
	}
	loggerModule.DPanic(msg, fields...)
}
