package logger

import (
	"reflect"
	"testing"

	"math/rand"
)

func TestLoggerBuffer(t *testing.T) {
	t.Run("Should be able to add a new message", func(t *testing.T) {
		loggerBuffer := createLoggerBufferStack(createMockLoggerRegister())

		want := createByteSlice(2)
		loggerBuffer.AddMessage(want)

		expectDeepEqual(t, loggerBuffer.messages[0], want)
	})

	t.Run("Should be able to store multiple messages", func(t *testing.T) {
		loggerBufferStack := createLoggerBufferStack(createMockLoggerRegister())

		messages := [][]byte{createByteSlice(2), createByteSlice(2)}
		for _, message := range messages {
			loggerBufferStack.AddMessage(message)
		}

		expectDeepEqual(t, loggerBufferStack.messages, messages)
	})

	t.Run("Should be able to get messages LIFO order", func(t *testing.T) {
		loggerBufferStack := createLoggerBufferStack(createMockLoggerRegister())

		messages := [][]byte{createByteSlice(2), createByteSlice(2)}
		for _, message := range messages {
			loggerBufferStack.AddMessage(message)
		}

		got := loggerBufferStack.GetMessages()
		expectDeepEqual(t, got, messages)
	})

	t.Run("Should be able to reduce size if the max size is catch", func(t *testing.T) {
		maxSize := 8
		quarter := maxSize / 4

		loggerBufferStack := loggerBuffer{
			messages: make([][]byte, 0),
			size:     0,
			maxSize:  8,
			register: createMockLoggerRegister(),
		}

		messages := [][]byte{
			createByteSlice(quarter),
			createByteSlice(quarter),
			createByteSlice(quarter),
			createByteSlice(quarter),
		}
		for _, message := range messages {
			loggerBufferStack.AddMessage(message)
		}

		gotedSize := loggerBufferStack.size
		wantedSize := quarter

		expectDeepEqual(t, gotedSize, wantedSize)
		expectDeepEqual(t, gotedSize, calcByteSliceSize(loggerBufferStack.messages)) // expect real size

		gotedMessages := loggerBufferStack.GetMessages()
		wantedMessges := messages[3:]
		expectDeepEqual(t, gotedMessages, wantedMessges)
	})

	t.Run("Should be able to save messages if maximum size is reached", func(t *testing.T) {
		maxSize := 2

		mockLoggerRegister := createMockLoggerRegister()
		loggerBufferStack := loggerBuffer{
			messages: make([][]byte, 0),
			size:     0,
			maxSize:  maxSize,
			register: mockLoggerRegister,
		}

		messages := [][]byte{createByteSlice(maxSize)}
		for _, message := range messages {
			loggerBufferStack.AddMessage(message)
		}

		if len(mockLoggerRegister.calls) == 0 {
			t.Fatal("Logger must be called if maximum size is reached")
		}

		got := mockLoggerRegister.calls
		want := messages

		expectDeepEqual(t, got, want)
	})

	t.Run("Should be able to stop the logger buffer", func(t *testing.T) {
		mockLoggerRegister := createMockLoggerRegister()
		loggerBufferStack := loggerBuffer{
			messages: make([][]byte, 0),
			size:     0,
			maxSize:  1024,
			register: mockLoggerRegister,
		}

		messages := [][]byte{createByteSlice(2), createByteSlice(3)}
		for _, message := range messages {
			loggerBufferStack.AddMessage(message)
		}

		loggerBufferStack.Stop()

		expectDeepEqual(t, loggerBufferStack.size, 0)
		expectDeepEqual(t, len(loggerBufferStack.messages), 0)

		if len(mockLoggerRegister.calls) == 0 {
			t.Fatal("Logger must be called once the stop function is called")
		}

		got := mockLoggerRegister.calls
		want := messages

		expectDeepEqual(t, got, want)
	})
}

func expectDeepEqual[T interface{}](t *testing.T, got, want T) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got: %v want: %v", got, want)
	}
}

func createByteSlice(size int) []byte {
	result := make([]byte, size)
	for i := range size {
		result[i] = randByte()
	}

	return result
}

func randByte() byte {
	return byte(rand.Intn(256))
}

func reverseSlice[T interface{}](slice []T) []T {
	n := len(slice)
	reversed := make([]T, n)

	for i := 0; i < n; i++ {
		reversed[i] = slice[n-1-i]
	}

	return reversed
}

func calcByteSliceSize(slice [][]byte) int {
	var totalBytes int
	for _, b := range slice {
		totalBytes += len(b)
	}

	return totalBytes
}

type mockLoggerRegister struct {
	calls [][]byte
}

func (mlr *mockLoggerRegister) Write(chunk []byte) (int, error) {
	mlr.calls = append(mlr.calls, chunk)
	return 0, nil
}

func createMockLoggerRegister() *mockLoggerRegister {
	return &mockLoggerRegister{
		calls: make([][]byte, 0),
	}
}
