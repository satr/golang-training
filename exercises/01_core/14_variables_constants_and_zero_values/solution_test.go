package variablesconstantsandzerovalues

import (
	"errors"
	"testing"
)

func TestLogLevelLabel(t *testing.T) {
	tests := []struct {
		name  string
		level LogLevel
		want  string
	}{
		{name: "debug level", level: LevelDebug, want: "DEBUG"},
		{name: "info level", level: LevelInfo, want: "INFO"},
		{name: "warn level", level: LevelWarn, want: "WARN"},
		{name: "error level", level: LevelError, want: "ERROR"},
		{name: "negative level", level: -1, want: "UNKNOWN"},
		{name: "too large level", level: 42, want: "UNKNOWN"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.level.Label()
			if got != test.want {
				t.Fatalf("Label() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestGetDefaultConfig(t *testing.T) {
	cfg := GetDefaultConfig()

	if cfg.Timeout != 0 {
		t.Fatalf("Timeout = %d, want 0", cfg.Timeout)
	}
	if cfg.RetryLimit != 0 {
		t.Fatalf("RetryLimit = %d, want 0", cfg.RetryLimit)
	}
	if cfg.Verbose != false {
		t.Fatalf("Verbose = %t, want false", cfg.Verbose)
	}
	if cfg.Labels != nil {
		t.Fatalf("Labels = %v, want nil", cfg.Labels)
	}
	if cfg.LogLevel != LevelDebug {
		t.Fatalf("LogLevel = %d, want %d", cfg.LogLevel, LevelDebug)
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
		want error
	}{
		{name: "valid config", cfg: Config{Timeout: 10, RetryLimit: 50}, want: nil},
		{name: "invalid timeout", cfg: Config{Timeout: 0, RetryLimit: 50}, want: errors.New("Timeout must be greater than zero")},
		{name: "negative timeout", cfg: Config{Timeout: -5, RetryLimit: 50}, want: errors.New("Timeout must be greater than zero")},
		{name: "negative retry limit", cfg: Config{Timeout: 10, RetryLimit: -1}, want: errors.New("RetryLimit must be between 0 and 100")},
		{name: "too large retry limit", cfg: Config{Timeout: 10, RetryLimit: 101}, want: errors.New("RetryLimit must be between 0 and 100")},
		{name: "both invalid", cfg: Config{Timeout: 0, RetryLimit: -1}, want: errors.New("Timeout must be greater than zero; RetryLimit must be between 0 and 100")},
		{name: "boundary valid timeout", cfg: Config{Timeout: 1, RetryLimit: 50}, want: nil},
		{name: "boundary valid retry limit", cfg: Config{Timeout: 10, RetryLimit: 0}, want: nil},
		{name: "boundary valid retry limit max", cfg: Config{Timeout: 10, RetryLimit: 100}, want: nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := ValidateConfig(test.cfg)
			if (got == nil) != (test.want == nil) {
				t.Fatalf("ValidateConfig() = %v, want %v", got, test.want)
			}
			if got != nil && got.Error() != test.want.Error() {
				t.Fatalf("ValidateConfig() = %q, want %q", got.Error(), test.want.Error())
			}
		})
	}
}
