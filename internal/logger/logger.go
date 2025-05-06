package logger

import (
	"fmt"
	"io"

	"github.com/pedrowilliamss/pingero-cli/internal/observer"
)

type Logger struct {
	writer io.Writer
}

func (l *Logger) SaveUrlStatus(status observer.UrlStatus) {
	l.writer.Write(l.urlStatusToBinary(status))
}

func (l *Logger) urlStatusToBinary(urlStatus observer.UrlStatus) []byte {
	status := "Failure"
	if urlStatus.Ok {
		status = "Success"
	}

	return []byte(fmt.Sprintf("Status: %s | Time Response: %dms", status, urlStatus.ResponseTime))
}
