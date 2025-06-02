package aggregations

import (
	"io"
	"net/http"

	"github.com/pedrowilliamss/pingero-cli/internal/infra"
	"github.com/pedrowilliamss/pingero-cli/internal/logger"
	"github.com/pedrowilliamss/pingero-cli/internal/observer"
)

type UrlAggregate struct {
	Url      string
	FilePath string
	Logger   *logger.Logger
	Observer *observer.Observer
}

func CreateUrlAggregate(url, filePath string) *UrlAggregate {
	fileLogger, err := logger.CreateLogger(filePath)
	if err != nil {
		panic(err)
	}

	observer := observer.CreateObserver(url, fileLogger, infra.HttpRequesterFn(http.Get))

	return &UrlAggregate{
		Url:      url,
		FilePath: filePath,
		Logger:   fileLogger,
		Observer: observer,
	}
}

func (a *UrlAggregate) AddProxy(w io.Writer) {
	a.Logger.AddProxy(w)
}

func (a *UrlAggregate) RemoveProxy() {
	a.Logger.RemoveProxy()
}

func (a *UrlAggregate) GetMessages() [][]byte {
	return a.Logger.GetRecentMessages()
}
