package infra

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pedrowilliamss/pingero-cli/internal/observer"
)

var httpRequesterFn = HttpRequesterFn(http.Get)

func TestHttpRequester_Ping(t *testing.T) {
	t.Run("should be able to execute an ping", func(t *testing.T) {
		responseTime := 20 * time.Millisecond
		ts := createFakeServer(http.StatusOK, responseTime)
		defer ts.Close()

		got := httpRequesterFn.Ping(ts.URL)

		expecPingResult(t, got, observer.PingResult{
			Success:      true,
			ResponseTime: responseTime,
		})
	})

	t.Run("Ping should return a failure when the url is unavailable", func(t *testing.T) {
		responseTime := 20 * time.Millisecond
		ts := createFakeServer(http.StatusInternalServerError, responseTime)
		defer ts.Close()

		got := httpRequesterFn.Ping(ts.URL)

		expecPingResult(t, got, observer.PingResult{
			Success:      false,
			ResponseTime: responseTime,
		})
	})

	t.Run("Ping shoudl return a failue when the url its wrong", func(t *testing.T) {
		got := httpRequesterFn.Ping("invalid-url")

		expecPingResult(t, got, observer.PingResult{
			Success:      false,
			ResponseTime: 0,
		})
	})
}

func createFakeServer(statusCode int, responseTime time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(responseTime)
		w.WriteHeader(statusCode)
	}))
}

func expecPingResult(t *testing.T, got, want observer.PingResult) {
	t.Helper()

	if got.Success != want.Success {
		t.Fatal("the ping was expected to be unsuccessful")
	}
	if got.ResponseTime < want.ResponseTime && got.ResponseTime > want.ResponseTime+5 {
		t.Errorf("the response time was expected to be between 50 and 60 ms. Goted response time: %d", got.ResponseTime)
	}
}
