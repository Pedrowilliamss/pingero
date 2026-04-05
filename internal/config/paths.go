package config

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Paths holds the filesystem locations used by pingero.
// Base is the root directory for all application data; ConfigFilePath
// is the location of the JSON configuration file.
type Paths struct {
	Base           string
	ConfigFilePath string
}

// PathWithDefaults returns a Paths rooted at [DefaultPathBase].
func PathWithDefaults() *Paths {
	base := DefaultPathBase()
	return &Paths{
		Base:           base,
		ConfigFilePath: filepath.Join(base, "config.json"),
	}
}

// PathWithBase returns a Paths rooted at base.
// If base is empty, it falls back to [PathWithDefaults].
func PathWithBase(base string) *Paths {
	if base == "" {
		return PathWithDefaults()
	}
	return &Paths{
		Base:           base,
		ConfigFilePath: filepath.Join(base, "config.json"),
	}
}

// DefaultPathBase returns the default base directory for application data,
// which is $HOME/.pingero. It panics if the home directory cannot be determined.
func DefaultPathBase() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(home, ".pingero")
}

// DefaultConfigFilePath returns the default path for the configuration file,
// which is $HOME/.pingero/config.json.
func DefaultConfigFilePath() string {
	return filepath.Join(DefaultPathBase(), "config.json")
}

// LogPath returns the directory where log files are stored.
func (p *Paths) LogPath() string {
	return filepath.Join(p.Base, "logs")
}

// CreateLogfileFor creates and returns the log file for the given URL,
// creating the logs directory if it does not exist.
func (p *Paths) CreateLogfileFor(rawURL string) (*os.File, error) {
	logPath := p.LogPath()
	if err := os.MkdirAll(logPath, 0o755); err != nil {
		return nil, err
	}

	return os.Create(filepath.Join(logPath, p.GenerateLogfileNameFor(rawURL)))
}

// GenerateLogfileNameFor returns the log filename for the given URL.
func (p *Paths) GenerateLogfileNameFor(rawURL string) string {
	return fmt.Sprintf("%s.log", sanitizeUrl(rawURL))
}

// sanitizeUrl converts a URL into a safe filename by replacing special
// characters with underscores. Falls back to a generic replacer if the URL
// cannot be parsed or has no host.
func sanitizeUrl(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil || u.Host == "" {
		return strings.NewReplacer(":", "_", "/", "_", ".", "_").Replace(rawURL)
	}

	host := strings.ReplaceAll(u.Host, ".", "_")
	trimmed := strings.Trim(u.Path, "/")
	if trimmed == "" {
		return host
	}
	return host + "_" + strings.ReplaceAll(trimmed, "/", "_")
}
