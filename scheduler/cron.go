package scheduler

// NewTickerElapseEveryMinute
/* @Description: 创建一个每分钟第几秒触发一次的ticker
 * @param second int
 * @return time.Ticker
 */
import (
	"time"
)

type Ticker struct {
	C    chan time.Time
	stop chan struct{}
}

func NewTickerElapseEveryMinute(elapse time.Duration) *Ticker {
	ret := &Ticker{
		C:    make(chan time.Time, 1),
		stop: make(chan struct{}, 1),
	}
	now := time.Now()

	s := EveryMinute(elapse)

	go func() {
		after := s.Next(now).Sub(now)
		timer := time.NewTimer(after)
		for {
			select {
			case now = <-timer.C:
				{
					ret.C <- now
					after = s.Next(now.Add(time.Second)).Sub(now)
					//log.Println("after", after)
					timer.Reset(after)
				}
			case <-ret.stop:
				{
					timer.Stop()
					return
				}
			}
		}
	}()

	return ret
}
