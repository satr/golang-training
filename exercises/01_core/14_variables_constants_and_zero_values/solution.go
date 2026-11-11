package variablesconstantsandzerovalues

import (
	"fmt"
	"strings"
)

// LogLevel represents a log severity level.
type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
)

// Label returns the human-readable label for the level.
func (l LogLevel) Label() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Config holds application configuration.
type Config struct {
	Timeout    int
	RetryLimit int
	Verbose    bool
	Labels     []string
	LogLevel   LogLevel
}

// GetDefaultConfig returns a Config at its zero value.
func GetDefaultConfig() Config {
	return Config{}
}

// ValidateConfig returns an error if Timeout <= 0 or RetryLimit is out of
// range [0, 100]. Returns nil when the config is valid.
func ValidateConfig(c Config) error {
	var problems []string
	if c.Timeout <= 0 {
		problems = append(problems, "Timeout must be greater than zero")
	}
	if c.RetryLimit < 0 || c.RetryLimit > 100 {
		problems = append(problems, "RetryLimit must be between 0 and 100")
	}
	if len(problems) > 0 {
		return fmt.Errorf("invalid config: %s", strings.Join(problems, "; "))
	}
	return nil
}
