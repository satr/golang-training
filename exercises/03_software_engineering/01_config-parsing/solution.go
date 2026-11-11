package configparsing

import (
	"fmt"
	"strconv"
	"strings"
)

// Config holds application configuration.
type Config struct {
	Name     string
	Port     int
	LogLevel string
}

// Defaults returns a Config with sensible default values.
func Defaults() Config {
	return Config{Name: "myapp", Port: 8080, LogLevel: "info"}
}

// LoadConfig builds a Config from defaults, file contents, and env vars.
func LoadConfig(fileContents string, env map[string]string) (Config, error) {
	cfg := Defaults()
	// Parse file contents
	cfg = parseFile(cfg, fileContents)
	// Overlay env vars
	cfg = parseEnv(cfg, env)
	// Validate
	if err := validate(cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func parseFile(cfg Config, raw string) Config {
	lines := splitLines(raw)
	for _, line := range lines {
		line = trimSpace(line)
		if line == "" {
			continue
		}
		parts := split(line, "=")
		if len(parts) != 2 {
			continue
		}
		key := trimSpace(parts[0])
		val := trimSpace(parts[1])
		switch key {
		case "Name":
			cfg.Name = val
		case "Port":
			p, err := strconv.Atoi(val)
			if err == nil {
				cfg.Port = p
			}
		case "LogLevel":
			cfg.LogLevel = val
		}
	}
	return cfg
}

func parseEnv(cfg Config, env map[string]string) Config {
	if name, ok := env["APP_NAME"]; ok && name != "" {
		cfg.Name = name
	}
	if port, ok := env["APP_PORT"]; ok && port != "" {
		p, err := strconv.Atoi(port)
		if err == nil {
			cfg.Port = p
		}
	}
	if level, ok := env["APP_LOG_LEVEL"]; ok && level != "" {
		cfg.LogLevel = level
	}
	return cfg
}

func validate(cfg Config) error {
	if cfg.Port < 1 || cfg.Port > 65535 {
		return fmt.Errorf("invalid port %d", cfg.Port)
	}
	return nil
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func split(s string, sep rune) []string {
	var parts []string
	start := 0
	for i := 0; i < len(s); i++ {
		if rune(s[i]) == sep {
			parts = append(parts, s[start:i])
			start = i + 1
		}
	}
	parts = append(parts, s[start:])
	return parts
}

func trimSpace(s string) string {
	return strings.TrimSpace(s)
}
