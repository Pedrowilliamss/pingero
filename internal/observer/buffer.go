package observer

import (
	"io"
	"math"
	"sync"
)

var (
	MAX_BUFFER_SIZE  = 400
	INIT_BUFFER_SIZE = MAX_BUFFER_SIZE / 4
)

type BufferRecord interface {
	ToBinary() []byte
}

type Buffer[T BufferRecord] interface {
	Push(record T)
	GetContent() []T
	FlushAndClear()
}

type buffer[T BufferRecord] struct {
	data    []T
	maxSize int
	sink    io.Writer
	mutex   sync.Mutex
}

func (b *buffer[T]) Push(record T) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.data = append(b.data, record)

	if b.isBufferFull() {
		b.cull()
	}
}

func (b *buffer[T]) GetContent() []T {
	return b.data
}

func (b *buffer[T]) FlushAndClear() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.saveChunk(b.data)
	b.data = make([]T, 0, INIT_BUFFER_SIZE)
}

func (b *buffer[T]) cull() {
	n := len(b.data)
	if n == 0 {
		return
	}

	quarter := int(math.Ceil(float64(n) * 0.75))
	b.saveChunk(b.data[:quarter])

	b.data = b.data[quarter:]
}

func (b *buffer[T]) isBufferFull() bool {
	return len(b.data) >= b.maxSize
}

func (b *buffer[T]) saveChunk(data []T) {
	binaryChunk := make([]byte, 0, len(data)*64)
	for _, record := range data {
		binaryChunk = append(binaryChunk, record.ToBinary()...)
		binaryChunk = append(binaryChunk, byte(10))
	}

	b.sink.Write(binaryChunk)
}

func BufferFrom[T BufferRecord](sink io.Writer) *buffer[T] {
	return &buffer[T]{
		data:    make([]T, 0, MAX_BUFFER_SIZE/4),
		maxSize: MAX_BUFFER_SIZE,
		sink:    sink,
	}
}
