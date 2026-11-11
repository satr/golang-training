//go:build todo

package configparsing

// Config holds application configuration.
type Config struct {
	Name     string
	Port     int
	LogLevel string
}

// Defaults returns a Config with sensible default values.
func Defaults() Config {
	// TODO: Return a config with Name="myapp", Port=8080, LogLevel="info".
	return Config{}
}

// LoadConfig builds a Config from defaults, file contents, and env vars.
func LoadConfig(fileContents string, env map[string]string) (Config, error) {
	// TODO: Start with Defaults(), overlay Key=Value lines from fileContents,
	// then override with non-empty env values. Validate port is 1–65535.
	return Config{}, nil
}
