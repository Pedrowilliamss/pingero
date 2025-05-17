package logger

import (
	"fmt"
	"os"
	"time"
)

const LOG_SEPARATOR_SYMBOL = " -- "

type FileLogger struct {
	file *os.File
}

type LogMessage interface {
	Content() string
}

func (fl *FileLogger) AppendMessage(logMessage LogMessage) {
	fl.file.Write([]byte(fl.formatMessage(logMessage.Content())))
}

func (fl *FileLogger) formatMessage(message string) string {
	currentTime := time.Now().Format(time.RFC1123)
	return fmt.Sprintf("[%s]%s%s\n", currentTime, LOG_SEPARATOR_SYMBOL, message)
}

func CreateFileLogger(path string) (*FileLogger, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}

	return &FileLogger{
		file: f,
	}, nil
}
