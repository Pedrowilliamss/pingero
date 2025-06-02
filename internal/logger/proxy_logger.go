package logger

import (
	"fmt"
	"io"
)

type ProxyLogger struct {
	url string
	w   io.Writer
}

func (p ProxyLogger) Write(m []byte) (int, error) {
	message := []byte(fmt.Sprintf("%s -> %s", p.url, string(m)))
	p.w.Write(message)

	return len(message), nil
}

func CreateProxyLogger(url string, w io.Writer) *ProxyLogger {
	return &ProxyLogger{
		url: url,
		w:   w,
	}
}
