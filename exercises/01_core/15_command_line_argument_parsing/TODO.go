//go:build todo

package commandlineargumentparsing

import (
	"fmt"
	"time"
)

// Runner holds CLI configuration values.
type Runner struct {
	Host    string
	Port    int
	Verbose bool
	Timeout time.Duration
}

// ParseArgs parses flag arguments into a Runner configuration.
// Arguments should be like: []string{"-host", "example.com", "-port", "9090"}
// Returns an error for unknown flags or malformed values.
func ParseArgs(args []string) (Runner, error) {
	// TODO: Use flag flag.NewFlagSet or manual parsing to interpret args.
	// TODO: Return a Runner with parsed values or defaults, and any error.
	_ = fmt.Errorf
	return Runner{}, nil
}

// Run returns a formatted description of the Runner's current configuration.
func (r Runner) Run() string {
	// TODO: Format the runner's settings as a single-line string.
	// Format: "host=HOST port=PORT verbose=VERB timeout=TO" where VERB is
	// "on" or "off" and TO is either "30s" or whatever the user set.
	return ""
}
