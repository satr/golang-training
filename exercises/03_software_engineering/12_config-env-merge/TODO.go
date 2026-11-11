//go:build todo

package configenvmerge

type Config struct{}

func DefaultConfig(prefix string) *Config                          { return nil }
func WithDefaults(cfg *Config, defaults map[string]string) *Config { return cfg }
func WithEnv(cfg *Config) *Config                                  { return cfg }
func (c *Config) Get(key string) (string, bool)                    { return "", false }
func (c *Config) String() string                                   { return "" }
