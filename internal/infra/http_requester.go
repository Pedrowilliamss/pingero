package infra

import (
	"net/http"
	"time"

	"github.com/pedrowilliamss/pingero-cli/internal/observer"
)

type HttpRequesterFn func(url string) (resp *http.Response, err error)

func (fn HttpRequesterFn) Ping(url string) (pingResult observer.PingResult) {
	start := time.Now()

	resp, err := fn(url)

	pingResult.ResponseTime = time.Since(start)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	pingResult.Success = resp.StatusCode >= 200 && resp.StatusCode < 300
	return
}
