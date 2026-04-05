package config

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type UrlConfig struct {
	Url        string `json:"url"`
	LoggerFile string `json:"logger_file"`
}

type PingeroConfig struct {
	Paths *Paths
	Urls  map[string]UrlConfig `json:"urls"`
	mux   sync.Mutex           `json:"-"`
}

func (p *PingeroConfig) AddUrl(url string) *UrlConfig {
	p.mux.Lock()
	defer p.mux.Unlock()

	if value, ok := p.Urls[url]; ok {
		return &value
	}

	newUrlConfig := UrlConfig{
		Url:        url,
		LoggerFile: p.Paths.GenerateLogfileNameFor(url),
	}
	p.Urls[url] = newUrlConfig

	return &newUrlConfig
}

func (p *PingeroConfig) RemoveUrl(url string) {
	p.mux.Lock()
	defer p.mux.Unlock()

	delete(p.Urls, url)
}

func CreatePingeroConfigFromFileSystem(paths *Paths) (*PingeroConfig, error) {
	data, err := os.ReadFile(paths.ConfigFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			pingeroConfig := NewConfig(paths)
			err := PingeroConfigToFile(pingeroConfig)
			if err != nil {
				panic(err)
			}

			return pingeroConfig, nil
		}
		return nil, err
	}

	var config PingeroConfig
	if err = json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func PingeroConfigToFile(pingeroConfig *PingeroConfig) error {
	data, err := json.Marshal(pingeroConfig)
	if err != nil {
		return err
	}

	err = os.MkdirAll(pingeroConfig.Paths.Base, 0o755)
	if err != nil {
		return err
	}
	err = os.MkdirAll(pingeroConfig.Paths.LogPath(), 0o755)
	if err != nil {
		return err
	}

	err = os.WriteFile(pingeroConfig.Paths.ConfigFilePath, data, 0o644)
	if err != nil {
		return err
	}

	return nil
}

func NewConfig(paths *Paths) *PingeroConfig {
	return &PingeroConfig{
		Urls:  make(map[string]UrlConfig, 10),
		Paths: paths,
		mux:   sync.Mutex{},
	}
}
