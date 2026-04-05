package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"

	pingero "github.com/pedrowilliamss/pingero/internal"
	"github.com/pedrowilliamss/pingero/internal/config"
	"github.com/pedrowilliamss/pingero/internal/observer"
	"github.com/pedrowilliamss/pingero/internal/presenter"
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
	flags := parseFlags()

	pingero, err := pingero.CreatePingero(
		pingero.WithURLs(flags.urls),
		pingero.WithChannel(viewerCH),
		pingero.WithPaths(flags.paths),
	)
	if err != nil {
		panic(err)
	}
	defer pingero.SaveAndStop()

	pingero.Start()

	go logging(pingero.GetUrls())
	shutdown()
}

func shutdown() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	<-stop
}

type urlsFlag []string

type Flags struct {
	paths *config.Paths
	urls  urlsFlag
}

func (u *urlsFlag) String() string {
	return strings.Join(*u, ", ")
}

func (i *urlsFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func parseFlags() *Flags {
	c := Flags{}
	c.paths = config.PathWithDefaults()

	flag.Var(&c.urls, "url", "URLs that should be noted")
	flag.StringVar(&c.paths.ConfigFilePath, "config", config.DefaultConfigFilePath(), "Path to the configuration file")
	flag.Parse()

	return &c
}
