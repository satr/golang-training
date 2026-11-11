package commandlineargumentparsing

import (
	"flag"
	"fmt"
	"strings"
	"time"
)

// Runner holds CLI configuration values.
type Runner struct {
	Host    string
	Port    int
	Verbose bool
	Timeout time.Duration
}

// defaults defines the fallback values for Runner fields.
var defaults = struct {
	Host    string
	Port    int
	Verbose bool
	Timeout time.Duration
}{
	Host:    "localhost",
	Port:    8080,
	Timeout: 30 * time.Second,
}

// ParseArgs parses flag arguments into a Runner configuration.
// Arguments should be like: []string{"-host", "example.com", "-port", "9090"}
// Returns an error for unknown flags or malformed values.
func ParseArgs(args []string) (Runner, error) {
	// Use a flag set to parse arguments without touching os.Args
	fs := flag.NewFlagSet("runner", flag.ContinueOnError)

	host := fs.String("host", defaults.Host, "server host")
	port := fs.Int("port", defaults.Port, "server port")
	verbose := fs.Bool("verbose", defaults.Verbose, "enable verbose logging")
	timeout := fs.Duration("timeout", defaults.Timeout, "request timeout")

	err := fs.Parse(args)
	if err != nil {
		return Runner{}, fmt.Errorf("parse args: %w", err)
	}

	return Runner{
		Host:    *host,
		Port:    *port,
		Verbose: *verbose,
		Timeout: *timeout,
	}, nil
}

// Run returns a formatted description of the Runner's current configuration.
func (r Runner) Run() string {
	verbose := "off"
	if r.Verbose {
		verbose = "on"
	}
	return fmt.Sprintf("host=%s port=%d verbose=%s timeout=%s", r.Host, r.Port, verbose, r.Timeout)
}

// Helper for tests: check if a string is a valid flag prefix.
func isValidFlag(s string) bool {
	return strings.HasPrefix(s, "-") && len(s) > 1
}
