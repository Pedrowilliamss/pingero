package logger

import (
	"io"
	"math"
	"sync"
)

var (
	DEFAULT_CHUNK_SIZE = 50
	MAX_BUFFER_SIZE    = 1 * 1024 * 1024 // 1MB
)

type loggerBuffer struct {
	messages [][]byte
	size     int
	maxSize  int
	register io.Writer
	mutex    sync.Mutex
}

func (lrb *loggerBuffer) AddMessage(message []byte) {
	lrb.mutex.Lock()
	defer lrb.mutex.Unlock()

	lrb.messages = append(lrb.messages, message)
	lrb.size += len(message)

	if lrb.isBufferFull() {
		lrb.cull()
	}
}

func (lrb *loggerBuffer) GetMessages() [][]byte {
	return lrb.messages
}

func (lrb *loggerBuffer) Stop() {
	lrb.mutex.Lock()
	defer lrb.mutex.Unlock()

	lrb.size = 0
	lrb.saveChunk(lrb.messages)
	lrb.messages = make([][]byte, 0)
}

func (lrb *loggerBuffer) cull() {
	n := len(lrb.messages)
	if n == 0 {
		return
	}

	quarter := int(math.Ceil(float64(n) * 0.75))
	lrb.saveChunk(lrb.messages[:quarter])

	newMessages := make([][]byte, len(lrb.messages[quarter:]))
	copy(newMessages, lrb.messages[quarter:])

	lrb.messages = newMessages
	lrb.calcSize()
}

func (lrb *loggerBuffer) isBufferFull() bool {
	return lrb.size >= lrb.maxSize
}

func (lrb *loggerBuffer) calcSize() {
	var totalBytes int
	for _, b := range lrb.messages {
		totalBytes += len(b)
	}

	lrb.size = totalBytes
}

func (lrb *loggerBuffer) saveChunk(messages [][]byte) {
	for _, message := range messages {
		lrb.register.Write(message)
	}
}

func createLoggerBufferStack(register io.Writer) loggerBuffer {
	return loggerBuffer{
		messages: make([][]byte, 0, 50),
		size:     0,
		maxSize:  MAX_BUFFER_SIZE,
		register: register,
	}
}
