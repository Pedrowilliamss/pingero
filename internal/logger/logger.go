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

type logger struct {
	file *os.File
}

func (fl *logger) AppendMessage(logMessage shared.LogMessage) {
	fl.file.Write([]byte(fl.formatMessage(logMessage.Content())))
}

func (fl *logger) formatMessage(message string) string {
	currentTime := time.Now().Format(time.RFC1123)
	return fmt.Sprintf("[%s]%s%s\n", currentTime, MESSAGE_SEPARATOR_SYMBOL, message)
}

func CreateFileLogger(url string) (*logger, error) {
	f, err := os.OpenFile(url, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}

		fmt.Printf("Arquivo de log para a url %s não encontrado, criando um novo...\n", url)
		f, err = os.Create(url)

		if err != nil {
			return nil, err
		}
	}

	return &logger{
		file: f,
	}, nil
}
