package scheduler

/**
 * @Author: lee
 * @Description:
 * @File: scheduler
 * @Date: 2022/3/17 4:52 下午
 */
import "time"

type Scheduler interface {
	Next(t time.Time) time.Time
}

type MinuterScheduler struct {
	delay time.Duration
}

type constantScheduler struct {
	delay time.Duration
	d     time.Duration
}

func EveryMinute(delay time.Duration) *constantScheduler {
	return &constantScheduler{
		delay: delay,
		d:     1 * time.Minute,
	}
}

func (s *constantScheduler) Next(t time.Time) time.Time {
	now := (time.Duration(t.Second())*time.Second + time.Duration(t.Nanosecond())) % s.d

	offset := (s.d + s.delay - now) % s.d

	return t.Add(offset)

}
