package scheduler

import (
	"log"
	"testing"
	"time"
)

/**
 * @Author: lee
 * @Description:
 * @File: cron_test
 * @Date: 2022/3/17 11:02 上午
 */

func Test_Cron(t *testing.T) {
	ticker := NewTickerElapseEveryMinute(5 * time.Second)
	log.Println(time.Now().String())
	//ticker := time.NewTicker(time.Second)
	for {
		select {
		case tick := <-ticker.C:
			{
				log.Println(tick.String())
			}
		}
	}

}
