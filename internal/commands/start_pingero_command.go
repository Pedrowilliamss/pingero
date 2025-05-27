package commands

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pedrowilliamss/pingero-cli/internal/logger"
)

type StartPingeroCommandMessage string

const (
	VERIFY_DIRECTORIES    StartPingeroCommandMessage = "Checking the file directory..."
	CREATING_DIRECTORIES  StartPingeroCommandMessage = "Creating required directories..."
	CHECKING_URL_FILE_LOG StartPingeroCommandMessage = "Checking urls file log"
	CREATING_URL_FILE_LOG StartPingeroCommandMessage = "Creating file log for the %s url"
	SETTINGS_COMPLETED    StartPingeroCommandMessage = "Settings completed!! 🚀"
)

type StartPingeroCommand struct {
	MessageConsumer io.Writer
}

func (spc *StartPingeroCommand) Start(urls []string) map[string]string {
	urlFileMap := make(map[string]string, 10)
	directoriesToCreat, ok := spc.verifyRequiredDirectories()
	if !ok {
		spc.createRequiredDirectories(directoriesToCreat)
	}

	for _, url := range urls {
		urlFileMap[url] = spc.createFileLogger(url)
	}

	spc.publishMessage(SETTINGS_COMPLETED)
	return urlFileMap
}

func (spc *StartPingeroCommand) verifyRequiredDirectories() ([]string, bool) {
	spc.publishMessage(VERIFY_DIRECTORIES)

	directoriesToCreat := make([]string, 0, 1)

	loggerDirExists := dirExists(logger.LOGGER_DIR)
	if !loggerDirExists {
		directoriesToCreat = append(directoriesToCreat, logger.LOGGER_DIR)
	}

	return directoriesToCreat, len(directoriesToCreat) == 0
}

func (spc *StartPingeroCommand) createRequiredDirectories(directories []string) {
	spc.publishMessage(CREATING_DIRECTORIES)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	for _, dir := range directories {
		err := os.MkdirAll(filepath.Join(homeDir, dir), 0755)
		if err != nil {
			panic(err)
		}
	}
}

func (spc *StartPingeroCommand) createFileLogger(url string) string {
	spc.publishMessage(CHECKING_URL_FILE_LOG)

	hashFileName := sha256.Sum256([]byte(url))
	stringFileName := hex.EncodeToString(hashFileName[:]) + ".log"

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	fileName := filepath.Join(homeDir, logger.LOGGER_DIR, stringFileName)

	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0o644)
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	if err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}

		spc.publishMessage(creatingUrlFileMessage(url))
		f, err = os.Create(fileName)

		if err != nil {
			panic(err)
		}
	}

	return fileName
}

func (spc *StartPingeroCommand) publishMessage(message StartPingeroCommandMessage) {
	spc.MessageConsumer.Write([]byte(message + "\n"))
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func creatingUrlFileMessage(url string) StartPingeroCommandMessage {
	message := fmt.Sprintf(string(CREATING_URL_FILE_LOG), url)
	return StartPingeroCommandMessage(message)
}
