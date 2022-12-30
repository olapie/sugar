package xtime

import "time"

type Clock interface {
	Now() time.Time
}

type LocalClock struct {
}

func (l LocalClock) Now() time.Time {
	return time.Now()
}

type SystemUpTimer interface {
	// SystemUpTime returns system uptime in seconds
	SystemUpTime() int64
}

type ServerClock struct {
	serverTime *time.Time
	upTime     int64
	timer      SystemUpTimer
}

func NewServerClock(serverTime time.Time, timer SystemUpTimer) *ServerClock {
	return &ServerClock{
		serverTime: &serverTime,
		upTime:     timer.SystemUpTime(),
		timer:      timer,
	}
}

func (l ServerClock) Now() time.Time {
	elapsedSeconds := l.timer.SystemUpTime() - l.upTime
	return l.serverTime.Add(time.Second * time.Duration(elapsedSeconds))
}
