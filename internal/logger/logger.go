package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/pedrowilliamss/pingero-cli/internal/shared"
)

var (
	MESSAGE_SEPARATOR_SYMBOL = " -- "
	LOGGER_DIR               = filepath.Join(shared.MAIN_DIR, "logs", "urls")
)

type Logger struct {
	file         *os.File
	proxy        io.Writer
	loggerBuffer loggerBuffer
}

func (l *Logger) AppendMessage(logMessage shared.LogMessage) {
	message := l.formatMessage(logMessage.Content())
	if l.proxy != nil {
		l.proxy.Write(message)
	}
	go l.file.Write(message)
}

func (l *Logger) AddProxy(proxy io.Writer) {
	l.proxy = proxy
}

func (l *Logger) RemoveProxy() {
	l.proxy = nil
}

func (l *Logger) GetRecentMessages() [][]byte {
	return l.loggerBuffer.GetMessages()
}

func (l *Logger) formatMessage(message string) []byte {
	currentTime := time.Now().Format(time.RFC1123)
	return []byte(fmt.Sprintf("[%s]%s%s\n", currentTime, MESSAGE_SEPARATOR_SYMBOL, message))
}

func CreateLogger(filePath string) (*Logger, error) {
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}

		f, err = os.Create(filePath)

		if err != nil {
			return nil, err
		}
	}

	return &Logger{
		file:         f,
		loggerBuffer: createLoggerBufferStack(f),
	}, nil
}

func CreateFileLoggerWithProxy(filePath string, proxy io.Writer) (*Logger, error) {
	logger, err := CreateLogger(filePath)
	if err != nil {
		return nil, err
	}

	logger.proxy = proxy
	return logger, nil
}
