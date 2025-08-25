package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"

	pingero "github.com/pedrowilliamss/pingero-cli/internal"
	"github.com/pedrowilliamss/pingero-cli/internal/observer"
	"github.com/pedrowilliamss/pingero-cli/internal/presenter"
)

var viewerCH = make(chan observer.UrlStatus, 1000)

func logging(urls []string) {
	presenters := presenter.CreateUrlsStatusPresenter(urls)
	for urlStatus := range viewerCH {
		presenter := presenters[urlStatus.Url]
		presenter.FromStatus(&urlStatus)

		os.Stdout.Write(append(presenter.ToBinary(), byte(10)))
	}
}

func main() {
	resolveParams()

	pingero, err := pingero.CreatePingero(
		pingero.WithURLs(urlsArgs),
		pingero.WithChannel(viewerCH),
	)
	if err != nil {
		panic(err)
	}
	defer pingero.SaveAndStop()

	pingero.Start()

	go logging(pingero.GetUrls())
	shutdown()
}

type urlsFlag []string

func (u *urlsFlag) String() string {
	return strings.Join(*u, ", ")
}

func (i *urlsFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var urlsArgs urlsFlag

func resolveParams() {
	flag.Var(&urlsArgs, "url", "URLs that should be noted")
	flag.Parse()
}

func shutdown() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	<-stop
}
