package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pedrowilliamss/pingero-cli/internal/shared"
)

const (
	MESSAGE_SEPARATOR_SYMBOL = " -- "
)

var (
	LOGGER_DIR = filepath.Join(shared.MAIN_DIR, "logs", "urls")
)

type FileLogger struct {
	file *os.File
}

func (fl *FileLogger) AppendMessage(logMessage shared.LogMessage) {
	fl.file.Write([]byte(fl.formatMessage(logMessage.Content())))
}

func (fl *FileLogger) formatMessage(message string) string {
	currentTime := time.Now().Format(time.RFC1123)
	return fmt.Sprintf("[%s]%s%s\n", currentTime, MESSAGE_SEPARATOR_SYMBOL, message)
}

func CreateFileLogger(url string) (*FileLogger, error) {
	urlFilePath := createUrlFilePath(url)

	f, err := os.OpenFile(urlFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}

		fmt.Printf("Arquivo de log para a url %s não encontrado, criando um novo...\n", url)
		f, err = os.Create(urlFilePath)

		if err != nil {
			return nil, err
		}
	}

	return &FileLogger{
		file: f,
	}, nil
}

func createUrlFilePath(fileName string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	logsDir := filepath.Join(homeDir, LOGGER_DIR)

	err = os.MkdirAll(logsDir, 0755)

	if err != nil {
		panic(err)
	}
	return filepath.Join(logsDir, fileName)
}
