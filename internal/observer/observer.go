package observer

import (
	"context"
	"fmt"
	"io"
	"time"
)

type Register interface {
	Push(record *UrlStatus)
	GetContent() []*UrlStatus
	FlushAndClear()
}

type HttpRequester interface {
	Ping(url string) *PingResult
}

type UrlStatus struct {
	Url          string
	Up           bool
	StatusCode   int
	ResponseTime time.Duration
	CheckedAt    time.Time
}

func (u *UrlStatus) ToBinary() []byte {
	status := "down"
	if u.Up {
		status = "up"
	}

	return fmt.Appendf(nil,
		"[%s] %s: status: %s status_code: %d response_time: %s",
		u.CheckedAt.Format(time.RFC1123),
		u.Url,
		status,
		u.StatusCode,
		u.ResponseTime,
	)
}

type PingResult struct {
	Success      bool
	StatusCode   int
	ResponseTime time.Duration
}

type Observer struct {
	url           string
	register      Register
	httpRequester HttpRequester
	running       bool
	viewer        chan<- UrlStatus
}

func (o *Observer) Run(ctx context.Context, tickerTime time.Duration) {
	o.running = true
	o.ping()

	ticker := time.NewTicker(tickerTime)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			o.running = false
			return
		case <-ticker.C:
			o.ping()
		}
	}
}

func (o *Observer) Stop() {
	o.running = false
	o.register.FlushAndClear()
}

func (o *Observer) IsRunning() bool {
	return o.running
}

func (o *Observer) GetCurrentStatus() []*UrlStatus {
	return o.register.GetContent()
}

func (o *Observer) ping() {
	pingResult := o.httpRequester.Ping(o.url)
	status := UrlStatus{
		Url:          o.url,
		Up:           pingResult.Success,
		StatusCode:   pingResult.StatusCode,
		ResponseTime: pingResult.ResponseTime,
		CheckedAt:    time.Now(),
	}
	o.register.Push(&status)

	if o.viewer != nil {
		o.viewer <- status
	}
}

func CreateObserver(Url string, repository io.Writer, httpRequester HttpRequester, viewer chan<- UrlStatus) *Observer {
	return &Observer{
		url:           Url,
		register:      BufferFrom[*UrlStatus](repository),
		httpRequester: httpRequester,
		running:       false,
		viewer:        viewer,
	}
}
