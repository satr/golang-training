//go:build todo

package variablesconstantsandzerovalues

import "fmt"

// LogLevel represents a log severity level.
type LogLevel int

const (
	// TODO: Use iota to assign sequential values starting at 0.
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
)

// Label returns the human-readable label for the level.
func (l LogLevel) Label() string {
	// TODO: Return "DEBUG", "INFO", "WARN", or "ERROR" based on level.
	// Return "UNKNOWN" for out-of-range values.
	_ = fmt.Errorf
	return ""
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
	// TODO: Return a zero-value Config.
	return Config{}
}

// ValidateConfig returns an error if Timeout <= 0 or RetryLimit is out of
// range [0, 100]. Returns nil when the config is valid.
func ValidateConfig(c Config) error {
	// TODO: Check Timeout > 0 and 0 <= RetryLimit <= 100.
	// Return a descriptive error or nil.
	return nil
}
