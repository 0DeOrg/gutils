package dumputils

import (
	"bytes"
	"errors"
	"gitlab.qihangxingchen.com/qt/gutils/logutils"
	"go.uber.org/zap"
	"runtime"
)

/**
 * @Author: lee
 * @Description:
 * @File: routine
 * @Date: 2021/10/14 2:53 下午
 */

func HandlePanic(v ...interface{}) {
	var err error
	var stack string
	if r := recover(); nil != r {
		stack = string(PanicTrace(4))
		switch r.(type) {
		case error:
			err = r.(error)
			break
		case string:
			err = errors.New(r.(string))
			break
		default:
			err = errors.New("Unknown panic")
		}

		pc := make([]uintptr, 1)
		numFrames := runtime.Callers(4, pc)
		if numFrames < 1 {
			return
		}

		frame, _ := runtime.CallersFrames(pc).Next()
		//log.Println("rame function, file, line", frame.Function, frame.File, frame.Line)

		//log.Println("panic stack:\n "+stack+"\n", err.Error())

		logutils.Error("frame function, file, line", zap.String("func", frame.Function), zap.String("file", frame.File), zap.Int("line", frame.Line))
		logutils.Fatal("panic stack:\n "+stack+"\n", zap.Error(err))

	}
}

func PanicTrace(kb int) []byte {
	s := []byte("/src/runtime/panic.go")
	e := []byte("\ngoroutine ")
	line := []byte("\n")
	stack := make([]byte, kb<<10) //4KB
	length := runtime.Stack(stack, true)
	start := bytes.Index(stack, s)
	stack = stack[start:length]
	start = bytes.Index(stack, line) + 1
	stack = stack[start:]
	end := bytes.LastIndex(stack, line)
	if end != -1 {
		stack = stack[:end]
	}
	end = bytes.Index(stack, e)
	if end != -1 {
		stack = stack[:end]
	}
	stack = bytes.TrimRight(stack, "\n")
	return stack
}
