package observer

import (
	"context"
	"time"
)

type ObserverLogger interface {
	SaveUrlStatus(urlState UrlStatus)
}

type HttpRequester interface {
	Ping(url string) PingResult
}

type UrlStatus struct {
	Ok           bool
	ResponseTime time.Duration
	CheckedAt    time.Time
}

type PingResult struct {
	Success      bool
	ResponseTime time.Duration
}

type Observer struct {
	Url           string
	logger        ObserverLogger
	httpRequester HttpRequester
	running       bool
}

func (o *Observer) Execute(ctx context.Context) {
	o.running = true

	ticker := time.NewTicker(1 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			o.running = false
			return
		case <-ticker.C:
			o.logger.SaveUrlStatus(o.checkUrlStatus())
		}
	}
}

func (o *Observer) IsRunning() bool {
	return o.running
}

func (o *Observer) checkUrlStatus() UrlStatus {
	time := time.Now()
	pingResult := o.httpRequester.Ping(o.Url)

	return UrlStatus{
		Ok:           pingResult.Success,
		ResponseTime: pingResult.ResponseTime,
		CheckedAt:    time,
	}
}
