package presenter

import (
	"fmt"
	"time"

	"github.com/pedrowilliamss/pingero-cli/internal/observer"
)

type MetricsTimerPresenter struct {
	MaxTime         time.Duration
	Mintime         time.Duration
	DownTime        time.Duration
	UpTime          time.Duration
	lastState       bool
	lastSwitchState time.Time
}

func (m *MetricsTimerPresenter) update(u *observer.UrlStatus) {
	if m.MaxTime < u.ResponseTime {
		m.MaxTime = u.ResponseTime
	}
	if m.Mintime == 0 || m.Mintime > u.ResponseTime {
		m.Mintime = u.ResponseTime
	}

	if u.Up != m.lastState {
		if u.Up == false {
			m.UpTime = time.Since(m.lastSwitchState)
			m.DownTime = 0
		} else {
			m.DownTime = time.Since(m.lastSwitchState)
			m.UpTime = 0
		}

		m.lastState = u.Up
		m.lastSwitchState = time.Now()
	} else {
		if u.Up == true {
			m.UpTime = time.Since(m.lastSwitchState)
		} else {
			m.DownTime = time.Since(m.lastSwitchState)
		}
	}
}

func createMetricsTimerPresenter() *MetricsTimerPresenter {
	return &MetricsTimerPresenter{
		lastSwitchState: time.Now(),
		lastState:       true,
	}
}

type UrlStatusPresenter struct {
	Up           bool
	Url          string
	TimerMetrics *MetricsTimerPresenter
	StatusCheck  map[int16]*struct {
		httpStatus int16
		checks     int
	}
}

func (u *UrlStatusPresenter) FromStatus(s *observer.UrlStatus) {
	status := int16(s.StatusCode)
	if value, ok := u.StatusCheck[status]; ok {
		value.checks++
	} else {
		u.StatusCheck[status] = &struct {
			httpStatus int16
			checks     int
		}{
			httpStatus: status,
			checks:     1,
		}
	}

	u.TimerMetrics.update(s)
}

func (u *UrlStatusPresenter) ToBinary() []byte {
	status := "down"
	if u.Up {
		status = "up"
	}

	b := fmt.Appendf(nil, "url: %s | status: %s | max_time: %s | min_time: %s | up_time: %s | down_time: %s | status_checks:",
		u.Url,
		status,
		u.TimerMetrics.MaxTime,
		u.TimerMetrics.Mintime,
		u.TimerMetrics.UpTime,
		u.TimerMetrics.DownTime,
	)

	for _, v := range u.StatusCheck {
		b = fmt.Appendf(b, " [%d: %d checks]", v.httpStatus, v.checks)
	}

	return b
}

func CreateUrlStatusPresenter(url string) *UrlStatusPresenter {
	return &UrlStatusPresenter{
		Up:           true,
		Url:          url,
		TimerMetrics: createMetricsTimerPresenter(),
		StatusCheck: map[int16]*struct {
			httpStatus int16
			checks     int
		}{},
	}
}

func CreateUrlsStatusPresenter(urls []string) map[string]*UrlStatusPresenter {
	presenters := make(map[string]*UrlStatusPresenter, len(urls))
	for _, url := range urls {
		presenters[url] = CreateUrlStatusPresenter(url)
	}

	return presenters
}
