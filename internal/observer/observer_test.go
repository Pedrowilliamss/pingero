package observer

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/pedrowilliamss/pingero-cli/internal/shared"
)

func TestObserver(t *testing.T) {
	t.Run("should be possible to know if an observer is running", func(t *testing.T) {
		httpRequester := createMockHttpRequester()
		observer := Observer{
			Url:           "mock-url.com.br",
			logger:        createFakeLogger(),
			httpRequester: httpRequester,
		}

		expectObserverRunningState(t, observer, false)

		ctx, cancelFn := context.WithCancel(context.Background())

		go observer.Execute(ctx, 1*time.Second)

		time.Sleep(50 * time.Millisecond)
		cancelFn()

		expectObserverRunningState(t, observer, true)
	})

	t.Run("should be able to log the state of the observer's url", func(t *testing.T) {
		httpRequester := createMockHttpRequester()
		logger := createFakeLogger()
		observer := Observer{
			Url:           "mock-url.com.br",
			logger:        logger,
			httpRequester: httpRequester,
		}

		mockSuccess := false
		mockResponseTime := 50 * time.Millisecond

		httpRequester.MockResult(PingResult{
			Success:      mockSuccess,
			ResponseTime: mockResponseTime,
		})

		ctx, cancelFn := context.WithCancel(context.Background())
		go observer.Execute(ctx, 1*time.Second)
		time.Sleep(50 * time.Millisecond)
		cancelFn()

		if len(logger.Calls) == 0 {
			t.Errorf("logger should be called")
		}

		got := logger.Calls[0]

		expectEqual(t, got.Ok, mockSuccess)
		expectEqual(t, got.ResponseTime, mockResponseTime)
	})

	t.Run("observer should be able to execute in loop", func(t *testing.T) {
		const NUMBER_OF_CALLS = 3
		mockHttpRequester := createMockHttpRequester()
		observer := Observer{
			Url:           "mock-url.com.br",
			logger:        createFakeLogger(),
			httpRequester: mockHttpRequester,
		}

		ctx, cancelFn := context.WithCancel(context.Background())

		go func() {
			defer cancelFn()
			for {
				select {
				case calls := <-mockHttpRequester.notifyCh:
					if calls >= NUMBER_OF_CALLS {
						return
					}
				case <-time.After(1 * time.Second):
					return
				}
			}
		}()

		observer.Execute(ctx, 1*time.Second)

		mockHttpRequester.mu.Lock()
		gotCall := len(mockHttpRequester.Calls)

		if gotCall != NUMBER_OF_CALLS {
			t.Errorf("expected the httpRequester to receive %d calls, but %d were made", NUMBER_OF_CALLS, gotCall)
		}
	})
}

func expectObserverRunningState(t *testing.T, observer Observer, want bool) {
	t.Helper()
	got := observer.IsRunning()
	if got != want {
		t.Errorf("Expecte that observer running state is: %t, but the heal state is %t", want, got)
	}
}

func expectEqual[T interface{}](t *testing.T, got, want T) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

type fakeLogger struct {
	mu    sync.Mutex
	Calls []UrlStatus
}

func (fl *fakeLogger) AppendMessage(status shared.LogMessage) {
	fl.mu.Lock()
	defer fl.mu.Unlock()

	fmt.Println(status.Content())
	urlStatus, _ := UrlStatusFromString(status.Content())
	fl.Calls = append(fl.Calls, *urlStatus)
}

func createFakeLogger() *fakeLogger {
	return &fakeLogger{
		Calls: make([]UrlStatus, 0),
		mu:    sync.Mutex{},
	}
}

type mockHttpRequester struct {
	Calls    []string
	mu       sync.Mutex
	notifyCh chan int
	result   PingResult
}

func (mhr *mockHttpRequester) MockResult(result PingResult) {
	mhr.result = result
}

func (mhr *mockHttpRequester) Ping(url string) PingResult {
	mhr.mu.Lock()
	defer mhr.mu.Unlock()

	mhr.Calls = append(mhr.Calls, url)
	mhr.notifyCh <- len(mhr.Calls)

	return mhr.result
}

func createMockHttpRequester() *mockHttpRequester {
	return &mockHttpRequester{
		Calls:    make([]string, 0, 1),
		mu:       sync.Mutex{},
		notifyCh: make(chan int, 10),
	}
}
