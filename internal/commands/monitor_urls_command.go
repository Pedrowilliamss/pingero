package commands

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/pedrowilliamss/pingero-cli/internal/infra"
	"github.com/pedrowilliamss/pingero-cli/internal/logger"
	"github.com/pedrowilliamss/pingero-cli/internal/observer"
)

type MonitorUrlsCommandMessage string

const (
	CREATING_PENDING_OBSERVERS MonitorUrlsCommandMessage = "Creating pending observers..."
	CREATING_PEDING_OBSERVER   MonitorUrlsCommandMessage = "Creating observer for %q"
	EXECUTING_OBSERVERS        MonitorUrlsCommandMessage = "Executing observers..."
)

type MonitorUrlsCommand struct {
	MessageConsumer io.Writer
	Urls            map[string]string
	observersMap    map[string]*observer.Observer
	mutex           sync.Mutex
}

func (muc *MonitorUrlsCommand) MonitorUrls() {
	muc.publishMessage(CREATING_PENDING_OBSERVERS)
	muc.createPendingObservers()

	muc.publishMessage(EXECUTING_OBSERVERS)
	muc.executeObservers()
}

func (muc *MonitorUrlsCommand) createPendingObservers() {
	muc.mutex.Lock()
	defer muc.mutex.Unlock()

	for url, fileName := range muc.Urls {
		if muc.observersMap[url] != nil {
			continue
		}

		muc.publishMessage(creatingPendingObserverMessage(url))
		fileLogger, err := logger.CreateFileLogger(fileName)
		if err != nil {
			panic(err)
		}

		muc.observersMap[url] = observer.CreateObserver(url, fileLogger, infra.HttpRequesterFn(http.Get))
	}
}

func (muc *MonitorUrlsCommand) executeObservers() (context.Context, context.CancelFunc) {
	generalCtx, generalCancelFn := context.WithCancel(context.Background())
	for _, observer := range muc.observersMap {
		if !observer.IsRunning() {
			go observer.Execute(generalCtx)
		}
	}

	return generalCtx, generalCancelFn
}

func (muc *MonitorUrlsCommand) AddUrl(url string, fileName string) {
	muc.mutex.Lock()
	defer muc.mutex.Unlock()

	if _, ok := muc.observersMap[url]; !ok {
		fileLogger, err := logger.CreateFileLogger(fileName)
		if err != nil {
			panic(err)
		}

		muc.observersMap[url] = observer.CreateObserver(url, fileLogger, infra.HttpRequesterFn(http.Get))
		muc.Urls[url] = fileName
	}
}

func CreateMonitorUrlsCommand(messageConsumer io.Writer, urlsMap map[string]string) *MonitorUrlsCommand {
	return &MonitorUrlsCommand{
		MessageConsumer: messageConsumer,
		Urls:            urlsMap,
		observersMap:    make(map[string]*observer.Observer),
		mutex:           sync.Mutex{},
	}
}

func (spc *MonitorUrlsCommand) publishMessage(message MonitorUrlsCommandMessage) {
	spc.MessageConsumer.Write([]byte(message + "\n"))
}

func creatingPendingObserverMessage(url string) MonitorUrlsCommandMessage {
	message := fmt.Sprintf(string(CREATING_PEDING_OBSERVER), url)
	return MonitorUrlsCommandMessage(message)
}
