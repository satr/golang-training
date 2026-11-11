package configenvmerge

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	values map[string]string
	prefix string
}

func DefaultConfig(prefix string) *Config {
	return &Config{
		prefix: prefix,
		values: map[string]string{},
	}
}

func WithDefaults(cfg *Config, defaults map[string]string) *Config {
	for k, v := range defaults {
		cfg.values[k] = v
	}
	return cfg
}

func WithEnv(cfg *Config) *Config {
	for _, e := range os.Environ() {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := parts[0]
		if !strings.HasPrefix(key, cfg.prefix+"_") {
			continue
		}
		configKey := strings.TrimPrefix(key, cfg.prefix+"_")
		cfg.values[configKey] = parts[1]
	}
	return cfg
}

func (c *Config) Get(key string) (string, bool) {
	v, ok := c.values[key]
	return v, ok
}

func (c *Config) String() string {
	parts := make([]string, 0, len(c.values))
	for k, v := range c.values {
		parts = append(parts, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(parts, ", ")
}
