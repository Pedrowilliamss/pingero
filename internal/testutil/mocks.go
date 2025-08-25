package testutil

import (
	"sync"
)

type MockBufferRecord struct {
	content []byte
}

func (m MockBufferRecord) ToBinary() []byte {
	return m.content
}

func CreateMockBufferRecord() MockBufferRecord {
	return MockBufferRecord{
		content: []byte{},
	}
}

func CreateMockBufferRecords(amount int) []MockBufferRecord {
	result := make([]MockBufferRecord, amount)
	for i := range amount {
		result[i] = MockBufferRecord{
			content: []byte{},
		}
	}
	return result
}

type MockBuffer struct {
	Calls                 []MockBufferRecord
	FlushAndClearIsCalled bool
}

func (m *MockBuffer) Push(record MockBufferRecord) {
	m.Calls = append(m.Calls, record)
}

func (m *MockBuffer) GetContent() []MockBufferRecord {
	return m.Calls
}

func (m *MockBuffer) FlushAndClear() {
	m.FlushAndClearIsCalled = true
}

func CreateMockBuffer() *MockBuffer {
	return &MockBuffer{
		Calls:                 make([]MockBufferRecord, 0),
		FlushAndClearIsCalled: false,
	}
}

type MockSink struct {
	Calls [][]byte
}

func (mlr *MockSink) Write(message []byte) (int, error) {
	mlr.Calls = append(mlr.Calls, message)
	return 0, nil
}

func CreateMockWriter() *MockSink {
	return &MockSink{
		Calls: make([][]byte, 0),
	}
}

type BufferRecord interface {
	ToBinary() []byte
}

func CreateMockRegister[T BufferRecord]() *MockRegister[T] {
	return &MockRegister[T]{
		mu:    sync.Mutex{},
		Calls: make([]T, 0),
		Chan:  make(chan struct{}),
	}
}

type MockRegister[T BufferRecord] struct {
	mu    sync.Mutex
	Calls []T
	Chan  chan struct{}
}

func (m *MockRegister[T]) Push(record T) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Calls = append(m.Calls, record)
	m.Chan <- struct{}{}
}

func (m *MockRegister[T]) GetContent() []T {
	return m.Calls
}

func (m *MockRegister[T]) FlushAndClear() {
}

type MockHttpRequester[T any] struct {
	Calls    []string
	Mu       sync.Mutex
	NotifyCh chan int
	result   T
}

func (mhr *MockHttpRequester[T]) MockResult(result T) {
	mhr.result = result
}

func (mhr *MockHttpRequester[T]) Ping(url string) T {
	mhr.Mu.Lock()
	defer mhr.Mu.Unlock()

	mhr.Calls = append(mhr.Calls, url)
	mhr.NotifyCh <- 0

	return mhr.result
}

func (mhr *MockHttpRequester[T]) CallCount() int {
	mhr.Mu.Lock()
	defer mhr.Mu.Unlock()
	return len(mhr.Calls)
}

func CreateMockHttpRequester[T any]() *MockHttpRequester[T] {
	return &MockHttpRequester[T]{
		Calls:    make([]string, 0, 10),
		NotifyCh: make(chan int, 0),
	}
}
