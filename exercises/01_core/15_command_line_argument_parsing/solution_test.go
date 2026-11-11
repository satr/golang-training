package commandlineargumentparsing

import (
	"testing"
	"time"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    Runner
		wantErr bool
	}{
		{
			name: "no args",
			args: []string{},
			want: Runner{
				Host:    "localhost",
				Port:    8080,
				Verbose: false,
				Timeout: 30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "only host",
			args: []string{"-host", "example.com"},
			want: Runner{
				Host:    "example.com",
				Port:    8080,
				Verbose: false,
				Timeout: 30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "only port",
			args: []string{"-port", "9090"},
			want: Runner{
				Host:    "localhost",
				Port:    9090,
				Verbose: false,
				Timeout: 30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "verbose flag",
			args: []string{"-verbose"},
			want: Runner{
				Host:    "localhost",
				Port:    8080,
				Verbose: true,
				Timeout: 30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "all flags set",
			args: []string{"-host", "server.com", "-port", "8000", "-verbose", "-timeout", "1m"},
			want: Runner{
				Host:    "server.com",
				Port:    8000,
				Verbose: true,
				Timeout: 1 * time.Minute,
			},
			wantErr: false,
		},
		{
			name: "custom timeout",
			args: []string{"-timeout", "1m30s"},
			want: Runner{
				Host:    "localhost",
				Port:    8080,
				Verbose: false,
				Timeout: 90 * time.Second,
			},
			wantErr: false,
		},
		{
			name:    "unknown flag",
			args:    []string{"-unknown", "flag"},
			wantErr: true,
		},
		{
			name:    "invalid port value",
			args:    []string{"-port", "notanumber"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseArgs(tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseArgs() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if got.Host != tt.want.Host {
				t.Errorf("ParseArgs() Host = %v, want %v", got.Host, tt.want.Host)
			}
			if got.Port != tt.want.Port {
				t.Errorf("ParseArgs() Port = %v, want %v", got.Port, tt.want.Port)
			}
			if got.Verbose != tt.want.Verbose {
				t.Errorf("ParseArgs() Verbose = %v, want %v", got.Verbose, tt.want.Verbose)
			}
			if got.Timeout != tt.want.Timeout {
				t.Errorf("ParseArgs() Timeout = %v, want %v", got.Timeout, tt.want.Timeout)
			}
		})
	}
}

func TestRunnerRun(t *testing.T) {
	tests := []struct {
		name string
		r    Runner
		want string
	}{
		{
			name: "defaults",
			r:    Runner{},
			want: "host=localhost port=8080 verbose=off timeout=30s",
		},
		{
			name: "custom values",
			r: Runner{
				Host:    "server.com",
				Port:    9000,
				Verbose: true,
				Timeout: 1 * time.Minute,
			},
			want: "host=server.com port=9000 verbose=on timeout=1m0s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.r.Run()
			if got != tt.want {
				t.Errorf("Runner.Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
