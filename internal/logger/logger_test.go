package logger

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/pedrowilliamss/pingero-cli/internal/observer"
)

func TestLogger_SaveUrlStatus(t *testing.T) {
	t.Run("should be able to save success UrlStatus", func(t *testing.T) {
		mockWriter := createMockWriter()
		logger := Logger{
			writer: mockWriter,
		}

		urlStatus := observer.UrlStatus{
			Ok:           true,
			ResponseTime: 10,
		}

		logger.SaveUrlStatus(urlStatus)

		got := mockWriter.Calls[0]
		want := []byte(fmt.Sprintf("Status: %s | Time Response: %dms", "Success", urlStatus.ResponseTime))

		expecEqual(t, got, want)
	})

	t.Run("should be able to save failure UrlStatus", func(t *testing.T) {
		mockWriter := createMockWriter()
		logger := Logger{
			writer: mockWriter,
		}

		urlStatus := observer.UrlStatus{
			Ok:           false,
			ResponseTime: 20,
		}

		logger.SaveUrlStatus(urlStatus)

		got := mockWriter.Calls[0]
		want := []byte(fmt.Sprintf("Status: %s | Time Response: %dms", "Failure", urlStatus.ResponseTime))

		expecEqual(t, got, want)
	})
}

func expecEqual[T interface{}](t *testing.T, got, want T) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got: %v want: %v", got, want)
	}
}

type mockWriter struct {
	Calls [][]byte
}

func (mr *mockWriter) Write(b []byte) (n int, err error) {
	mr.Calls = append(mr.Calls, b)
	return 0, nil
}

func createMockWriter() *mockWriter {
	return &mockWriter{
		Calls: make([][]byte, 0),
	}
}
