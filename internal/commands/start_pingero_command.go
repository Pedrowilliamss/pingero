package commands

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pedrowilliamss/pingero-cli/internal/aggregations"
	"github.com/pedrowilliamss/pingero-cli/internal/logger"
)

type StartPingeroCommandMessage string

const (
	VERIFY_DIRECTORIES    StartPingeroCommandMessage = "Checking the file directory..."
	CREATING_DIRECTORIES  StartPingeroCommandMessage = "Creating required directories..."
	CHECKING_URL_FILE_LOG StartPingeroCommandMessage = "Checking urls file log"
	CREATING_URL_FILE_LOG StartPingeroCommandMessage = "Creating file log for the %s url"
	CREATING_OBSERVER     StartPingeroCommandMessage = "Creating observer for the %s url"
	SETTINGS_COMPLETED    StartPingeroCommandMessage = "Settings completed!! 🚀"
)

type StartPingeroCommand struct {
	MessagePublisher io.Writer
}

func (c *StartPingeroCommand) Execute(urls []string) map[string]*aggregations.UrlAggregate {
	urlAggregationMap := make(map[string]*aggregations.UrlAggregate, len(urls))
	directoriesToCreat, ok := c.verifyRequiredDirectories()
	if !ok {
		c.createRequiredDirectories(directoriesToCreat)
	}

	c.publishMessage(CHECKING_URL_FILE_LOG)
	for _, url := range urls {
		c.publishMessage(creatingUrlFileMessage(url))
		fileName := c.createFileLogger(url)

		c.publishMessage(creatingObserverMessage(url))
		urlAggregationMap[url] = aggregations.CreateUrlAggregate(url, fileName)
	}

	c.publishMessage(SETTINGS_COMPLETED)
	return urlAggregationMap
}

func (c *StartPingeroCommand) verifyRequiredDirectories() ([]string, bool) {
	c.publishMessage(VERIFY_DIRECTORIES)

	directoriesToCreat := make([]string, 0, 1)

	loggerDirExists := dirExists(logger.LOGGER_DIR)
	if !loggerDirExists {
		directoriesToCreat = append(directoriesToCreat, logger.LOGGER_DIR)
	}

	return directoriesToCreat, len(directoriesToCreat) == 0
}

func (c *StartPingeroCommand) createRequiredDirectories(directories []string) {
	c.publishMessage(CREATING_DIRECTORIES)

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

func (c *StartPingeroCommand) createFileLogger(url string) string {
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

		f, err = os.Create(fileName)

		if err != nil {
			panic(err)
		}
	}

	return fileName
}

func (c *StartPingeroCommand) publishMessage(message StartPingeroCommandMessage) {
	c.MessagePublisher.Write([]byte(message + "\n"))
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

func creatingObserverMessage(url string) StartPingeroCommandMessage {
	message := fmt.Sprintf(string(CREATING_OBSERVER), url)
	return StartPingeroCommandMessage(message)
}
