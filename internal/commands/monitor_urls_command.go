package commands

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/pedrowilliamss/pingero-cli/internal/aggregations"
)

type MonitorUrlsCommandMessage string

const (
	CREATING_PEDING_OBSERVER MonitorUrlsCommandMessage = "Creating observer for %q"
	EXECUTING_OBSERVERS      MonitorUrlsCommandMessage = "Executing observers..."
)

type MonitorUrlsCommand struct {
	MessagePublisher io.Writer
	UrlAggregateMap  map[string]*aggregations.UrlAggregate
	mutex            sync.Mutex
	cancelFn         context.CancelFunc
	running          bool
}

func (muc *MonitorUrlsCommand) Execute() {
	muc.publishMessage(EXECUTING_OBSERVERS)
	_, muc.cancelFn = muc.executeObservers()
	muc.running = true
}

func (muc *MonitorUrlsCommand) Stop() {
	muc.mutex.Lock()
	defer muc.mutex.Unlock()

	muc.cancelFn()
	muc.running = false
}

func (muc *MonitorUrlsCommand) IsRunning() bool {
	return muc.running
}

func (c *MonitorUrlsCommand) executeObservers() (context.Context, context.CancelFunc) {
	generalCtx, generalCancelFn := context.WithCancel(context.Background())

	for _, aggregate := range c.UrlAggregateMap {
		go aggregate.Observer.Execute(generalCtx, 15*time.Second)
	}

	return generalCtx, generalCancelFn
}

func (c *MonitorUrlsCommand) AddUrl(url string, filePath string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ok := c.UrlAggregateMap[url]; !ok {
		c.UrlAggregateMap[url] = aggregations.CreateUrlAggregate(url, filePath)
	}
}

func CreateMonitorUrlsCommand(messageConsumer io.Writer, urlsAggregation map[string]*aggregations.UrlAggregate) *MonitorUrlsCommand {
	return &MonitorUrlsCommand{
		MessagePublisher: messageConsumer,
		UrlAggregateMap:  urlsAggregation,
		mutex:            sync.Mutex{},
	}
}

func (c *MonitorUrlsCommand) publishMessage(message MonitorUrlsCommandMessage) {
	c.MessagePublisher.Write([]byte(message + "\n"))
}
