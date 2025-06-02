package commands

import (
	"io"

	"github.com/pedrowilliamss/pingero-cli/internal/aggregations"
	"github.com/pedrowilliamss/pingero-cli/internal/logger"
)

type ViewLogsCommand struct {
	w               io.Writer
	urlsAggregation map[string]*aggregations.UrlAggregate
}

func (c *ViewLogsCommand) Execute() {
	for _, aggregate := range c.urlsAggregation {
		proxyLogger := logger.CreateProxyLogger(aggregate.Url, c.w)
		aggregate.Logger.AddProxy(proxyLogger)
	}
}

func (c *ViewLogsCommand) Stop() {
	for _, aggregate := range c.urlsAggregation {
		aggregate.Logger.RemoveProxy()
	}
}

func CreateViewLogsCommand(messagePublisher io.Writer, urlsAggregation map[string]*aggregations.UrlAggregate) *ViewLogsCommand {
	return &ViewLogsCommand{
		urlsAggregation: urlsAggregation,
		w:               messagePublisher,
	}
}
