package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

var main_dir, configPath, configFile, logPath string

type UrlConfig struct {
	Url        string `json:"url"`
	LoggerFile string `json:"logger_file"`
}

type PingeroConfig struct {
	Urls map[string]UrlConfig `json:"urls"`
	mux  sync.Mutex           `json:"-"`
}

func (p *PingeroConfig) AddUrl(url string) *UrlConfig {
	p.mux.Lock()
	defer p.mux.Unlock()

	if value, ok := p.Urls[url]; ok {
		return &value
	}

	newUrlConfig := UrlConfig{
		Url:        url,
		LoggerFile: createLoggerFileName(url),
	}
	p.Urls[url] = newUrlConfig

	return &newUrlConfig
}

func (p *PingeroConfig) RemoveUrl(url string) {
	p.mux.Lock()
	defer p.mux.Unlock()

	if _, ok := p.Urls[url]; !ok {
		return
	}
	delete(p.Urls, url)
}

func CreateLoggerFiler(url string) *os.File {
	file, err := os.Create(filepath.Join(logPath, createLoggerFileName(url)))
	if err != nil {
		panic(err)
	}

	return file
}

func createLoggerFileName(url string) string {
	err := os.MkdirAll(logPath, 0o755)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s.log", sanitizeUrl(url))
}

func sanitizeUrl(url string) string {
	tokens := make([]string, 0, 2)
	for _, token := range strings.Split(strings.Split(url, "//")[1], ".") {
		tokens = append(tokens, strings.Split(token, "/")...)
	}

	return strings.Join(tokens, "_")
}

func CreatePingeroConfigFromFileSystem() (*PingeroConfig, error) {
	resolvePath()

	data, err := os.ReadFile(configFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			pingeroConfig := NewConfig()
			PingeroConfigToFile(pingeroConfig)

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

func PingeroConfigToFile(pingeroConfig *PingeroConfig) {
	fmt.Println(pingeroConfig)

	data, err := json.Marshal(pingeroConfig)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(configPath, 0o755)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(configFile, data, 0o644)
	if err != nil {
		panic(err)
	}
}

func NewConfig() *PingeroConfig {
	return &PingeroConfig{
		Urls: make(map[string]UrlConfig, 10),
		mux:  sync.Mutex{},
	}
}

func resolvePath() {
	main_dir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	configPath = path.Join(main_dir, "pingero")
	configFile = path.Join(configPath, "pingero_config.json")
	logPath = path.Join(main_dir, "logs")
}
