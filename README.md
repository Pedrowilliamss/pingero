# Pingero

Pingero is a lightweight CLI tool for real-time URL health monitoring. It continuously pings one or more URLs and reports live status metrics directly in your terminal.

## What it does

- Continuously pings each URL at a configurable interval
- Displays live metrics per URL: response time (max/min), uptime, downtime, and HTTP status code frequency
- Persists monitoring history to disk (`~/logs/`) and saves your URL list between sessions (`~/.pingero/pingero_config.json`)
- Handles multiple URLs simultaneously, each monitored in its own goroutine
- Flushes buffered logs to disk automatically and on graceful shutdown (Ctrl+C)

## Output example

```
url: https://example.com | status: up | max_time: 212ms | min_time: 38ms | up_time: 5m12s | down_time: 0s | status_checks: [200: 62 checks]
url: https://api.example.com | status: down | max_time: 541ms | min_time: 102ms | up_time: 3m44s | down_time: 1m28s | status_checks: [200: 45 checks] [503: 17 checks]
```

## Installation

```bash
git clone https://github.com/pedrowilliamss/pingero-cli.git
cd pingero-cli
go build -o pingero/pingero-cli ./cmd/pingero-cli/
```

## Usage

```bash
# Monitor a single URL
./pingero/pingero-cli -url https://example.com

# Monitor multiple URLs
./pingero/pingero-cli -url https://example.com -url https://api.example.com
```

URLs passed via `-url` are merged with any URLs already saved in your config, so previously monitored URLs resume automatically.

Stop monitoring with `Ctrl+C` — all data is flushed to disk before exit.

## Requirements

- Go 1.21+
- No external dependencies
