// Package config loads and validates service configuration.
package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	defaultSessionPath  = "test/data/sample-session.jsonl"
	defaultPollInterval = 5 * time.Second
)

// Config contains the session reader's runtime settings.
type Config struct {
	SessionPath  string
	PollInterval time.Duration
}

// Load reads optional values from .env, then applies process environment values.
// Existing process environment variables take precedence over .env.
func Load(dotenvPath string) (Config, error) {
	values, err := readDotenv(dotenvPath)
	if err != nil {
		return Config{}, err
	}

	lookup := func(key string) string {
		if value, ok := os.LookupEnv(key); ok {
			return value
		}
		return values[key]
	}

	sessionPath := strings.TrimSpace(lookup("CODEX_SESSION_PATH"))
	if sessionPath == "" {
		// Keep the older directory-oriented name as a convenient alias.
		sessionPath = strings.TrimSpace(lookup("CODEX_SESSION_DIR"))
	}
	if sessionPath == "" {
		sessionPath = defaultSessionPath
	}

	pollInterval := defaultPollInterval
	if raw := strings.TrimSpace(lookup("POLL_INTERVAL")); raw != "" {
		pollInterval, err = time.ParseDuration(raw)
		if err != nil {
			return Config{}, fmt.Errorf("parse POLL_INTERVAL %q: %w", raw, err)
		}
	}
	if pollInterval <= 0 {
		return Config{}, errors.New("POLL_INTERVAL must be greater than zero")
	}

	return Config{SessionPath: sessionPath, PollInterval: pollInterval}, nil
}

func readDotenv(path string) (map[string]string, error) {
	values := make(map[string]string)

	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return values, nil
	}
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for lineNumber := 1; scanner.Scan(); lineNumber++ {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		key, value, found := strings.Cut(line, "=")
		if !found || strings.TrimSpace(key) == "" {
			return nil, fmt.Errorf("parse %s line %d: expected KEY=VALUE", path, lineNumber)
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if len(value) >= 2 && ((value[0] == '"' && value[len(value)-1] == '"') ||
			(value[0] == '\'' && value[len(value)-1] == '\'')) {
			value = value[1 : len(value)-1]
		}
		values[key] = value
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}

	return values, nil
}
