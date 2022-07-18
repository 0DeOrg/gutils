package dumputils

/**
 * @Author: lee
 * @Description:
 * @File: signal
 * @Date: 2021/10/13 6:05 下午
 */

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func RegisterSignal(cbExit func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT:
				if nil != cbExit {
					cbExit()
				}

			default:
				fmt.Println("other signal", s)
			}
		}
	}()
}
