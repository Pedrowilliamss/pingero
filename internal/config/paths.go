package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Paths struct {
	Base           string
	ConfigFilePath string
}

func PathWithDefaults() *Paths {
	return PathWithBase(DefaultPathBase())
}

// Gets the default base directory for path
func DefaultPathBase() string {
	base, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return path.Join(base, ".pingero")
}

// Returns the default config file location
func DefaultConfigFilePath() string {
	return path.Join(DefaultPathBase(), "config.json")
}

// Gets the main paths using the provided directory as a base
func PathWithBase(base string) *Paths {
	return &Paths{
		Base:           base,
		ConfigFilePath: DefaultConfigFilePath(),
	}
}

// Logs returns the path to the logs directory
func (p *Paths) LogPath() string {
	return path.Join(p.Base, "logs")
}

// Creates a log file for the provided URL
func (p *Paths) CreateLogfileFor(url string) (*os.File, error) {
	logPath := p.LogPath()
	err := os.MkdirAll(logPath, 0o755)
	if err != nil {
		return nil, err
	}

	file, err := os.Create(filepath.Join(logPath, p.GenerateLogfileNameFor(url)))
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (p *Paths) GenerateLogfileNameFor(url string) string {
	return fmt.Sprintf("%s.log", sanitizeUrl(url))
}

func sanitizeUrl(url string) string {
	tokens := make([]string, 0, 2)
	for token := range strings.SplitSeq(strings.Split(url, "//")[1], ".") {
		tokens = append(tokens, strings.Split(token, "/")...)
	}

	return strings.Join(tokens, "_")
}
