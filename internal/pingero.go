package pingero

import (
	"context"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/pedrowilliamss/pingero-cli/internal/config"
	"github.com/pedrowilliamss/pingero-cli/internal/config/paths"
	"github.com/pedrowilliamss/pingero-cli/internal/infra"
	"github.com/pedrowilliamss/pingero-cli/internal/observer"
)

var httpRequesterFn = infra.HttpRequesterFn(http.Get)

type Pingero struct {
	config           *config.PingeroConfig
	observersMap     map[string]*observer.Observer
	urlStatusChannel chan observer.UrlStatus
}

func (p *Pingero) Start() {
	for _, observer := range p.observersMap {
		go observer.Run(context.Background(), 5*time.Second)
	}
}

func (p *Pingero) StartObserver(url string) {
	observer, ok := p.observersMap[url]
	if !ok {
		return
	}

	go observer.Run(context.Background(), 5*time.Second)
}

func (p *Pingero) AddNewUrl(url string) {
	if _, ok := p.observersMap[url]; ok {
		return
	}

	p.config.AddUrl(url)
	p.observersMap[url] = createObserver(p.config.Paths, url, p.urlStatusChannel)
}

func (p *Pingero) RemoveUrl(url string) {
	p.config.RemoveUrl(url)
	p.observersMap[url].Stop()

	delete(p.observersMap, url)
}

func (p *Pingero) SaveAndStop() {
	config.PingeroConfigToFile(p.config)
	var wg sync.WaitGroup
	for _, observer := range p.observersMap {
		wg.Go(func() { observer.Stop() })
	}

	wg.Wait()
}

func (p *Pingero) StopObserver(url string) {
	observer, ok := p.observersMap[url]
	if !ok {
		return
	}

	observer.Stop()
}

func (p *Pingero) GetUrlsCount() int {
	return len(p.observersMap)
}

func (p *Pingero) GetUrls() []string {
	urls := make([]string, 0, p.GetUrlsCount())
	for url := range p.observersMap {
		urls = append(urls, url)
	}

	return urls
}

type PingeroOptionsFunc func(opts *PingeroOptions)

type PingeroOptions struct {
	urls             []string
	urlStatusChannel chan observer.UrlStatus
	paths            *paths.Paths
}

func defaultPingeroOptions() *PingeroOptions {
	return &PingeroOptions{
		urls:             make([]string, 0),
		urlStatusChannel: nil,
		paths:            paths.New(),
	}
}

func WithChannel(channel chan observer.UrlStatus) PingeroOptionsFunc {
	return func(p *PingeroOptions) {
		p.urlStatusChannel = channel
	}
}

func WithURLs(urls []string) PingeroOptionsFunc {
	return func(p *PingeroOptions) {
		p.urls = urls
	}
}

func WithPaths(paths *paths.Paths) PingeroOptionsFunc {
	return func(p *PingeroOptions) {
		p.paths = paths
	}
}

func CreatePingero(optionsFn ...PingeroOptionsFunc) (*Pingero, error) {
	opts := defaultPingeroOptions()
	for _, fn := range optionsFn {
		fn(opts)
	}

	pingeroConfig, err := config.CreatePingeroConfigFromFileSystem(opts.paths)
	if err != nil {
		return nil, err
	}

	var ch chan observer.UrlStatus
	if opts.urlStatusChannel != nil {
		ch = opts.urlStatusChannel
	} else {
		ch = make(chan observer.UrlStatus, max(len(pingeroConfig.Urls), 1)*1000)
	}

	if len(pingeroConfig.Urls) < len(opts.urls) {
		for _, url := range opts.urls {
			pingeroConfig.AddUrl(url)
		}
	}

	urls := make([]string, 0, len(pingeroConfig.Urls)+len(opts.urls))
	for url := range pingeroConfig.Urls {
		urls = append(urls, url)
	}
	observersMap := createObserversMap(opts.paths, append(urls, opts.urls...), ch)

	return &Pingero{
		config:           pingeroConfig,
		observersMap:     observersMap,
		urlStatusChannel: ch,
	}, nil
}

func createObserversMap(paths *paths.Paths, urls []string, ch chan<- observer.UrlStatus) map[string]*observer.Observer {
	observersMap := make(map[string]*observer.Observer, len(urls))
	for _, url := range urls {
		logfile := paths.CreateLogfileFor(url)
		observersMap[url] = createObserver(url, logfile, ch)
	}

	return observersMap
}

func createObserver(url string, logfile *os.File, viewer chan<- observer.UrlStatus) *observer.Observer {
	return observer.CreateObserver(url, logfile, httpRequesterFn, viewer)
}
