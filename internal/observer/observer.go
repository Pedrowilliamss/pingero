package observer

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pedrowilliamss/pingero-cli/internal/shared"
)

type ObserverLogger interface {
	AppendMessage(statusLogger shared.LogMessage)
}

type HttpRequester interface {
	Ping(url string) PingResult
}

type UrlStatus struct {
	Ok           bool
	ResponseTime time.Duration
}

func (us UrlStatus) Content() string {
	status := "Down"
	if us.Ok {
		status = "Up"
	}
	return fmt.Sprintf("Status: %s Response Time: %dms", status, us.ResponseTime.Milliseconds())
}

func UrlStatusFromString(s string) (*UrlStatus, error) {
	var us UrlStatus

	parts := strings.Fields(s)
	if len(parts) != 5 {
		return nil, fmt.Errorf("invalid format")
	}

	statusStr := parts[1]
	switch statusStr {
	case "Up":
		us.Ok = true
	case "Down":
		us.Ok = false
	default:
		return nil, fmt.Errorf("invalid status: %s", statusStr)
	}

	timeStr := parts[4]
	if !strings.HasSuffix(timeStr, "ms") {
		return nil, fmt.Errorf("invalid response time: %s", timeStr)
	}

	timeNumStr := strings.TrimSuffix(timeStr, "ms")
	ms, err := strconv.Atoi(timeNumStr)
	if err != nil {
		return nil, fmt.Errorf("invalid response time: %v", err)
	}

	us.ResponseTime = time.Duration(ms) * time.Millisecond

	return &us, nil
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
			urlStatus := o.checkUrlStatus()
			o.logger.AppendMessage(urlStatus)
		}
	}
}

func (o *Observer) IsRunning() bool {
	return o.running
}

func (o *Observer) checkUrlStatus() UrlStatus {
	pingResult := o.httpRequester.Ping(o.Url)

	return UrlStatus{
		Ok:           pingResult.Success,
		ResponseTime: pingResult.ResponseTime,
	}
}

func CreateObserver(Url string, logger ObserverLogger, httpRequester HttpRequester) *Observer {
	return &Observer{
		Url:           Url,
		logger:        logger,
		httpRequester: httpRequester,
		running:       false,
	}
}
