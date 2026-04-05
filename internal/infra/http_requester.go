package infra

import (
	"net/http"
	"time"

	"github.com/pedrowilliamss/pingero/internal/observer"
)

type HttpRequesterFn func(url string) (resp *http.Response, err error)

func (fn HttpRequesterFn) Ping(url string) *observer.PingResult {
	pingResult := observer.PingResult{}
	start := time.Now()

	resp, err := fn(url)
	if err != nil {
		return &pingResult
	}

	pingResult.ResponseTime = time.Since(start)
	defer resp.Body.Close()

	pingResult.Success = resp.StatusCode >= 200 && resp.StatusCode < 300
	pingResult.StatusCode = resp.StatusCode

	return &pingResult
}
