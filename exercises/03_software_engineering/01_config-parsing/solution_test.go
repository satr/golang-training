package configparsing

import (
	"fmt"
	"strings"
	"testing"
)

func TestDefaults(t *testing.T) {
	cfg := Defaults()

	want := Config{Name: "myapp", Port: 8080, LogLevel: "info"}
	if cfg != want {
		t.Fatalf("Defaults() = %+v; want %+v", cfg, want)
	}
}

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name      string
		file      string
		env       map[string]string
		wantPort  int
		wantName  string
		wantLevel string
		wantErr   bool
	}{
		{
			name:      "defaults only",
			file:      "",
			env:       map[string]string{},
			wantPort:  8080,
			wantName:  "myapp",
			wantLevel: "info",
			wantErr:   false,
		},
		{
			name:      "file override",
			file:      "Name=custom\nPort=9090",
			env:       map[string]string{},
			wantPort:  9090,
			wantName:  "custom",
			wantLevel: "info",
			wantErr:   false,
		},
		{
			name:      "env override",
			file:      "Name=fromfile\nPort=9090",
			env:       map[string]string{"APP_NAME": "envapp"},
			wantPort:  9090,
			wantName:  "envapp",
			wantLevel: "info",
			wantErr:   false,
		},
		{
			name:      "full overlay",
			file:      "Name=fromfile\nPort=9090\nLogLevel=debug",
			env:       map[string]string{"APP_NAME": "envapp", "APP_LOG_LEVEL": "error"},
			wantPort:  9090,
			wantName:  "envapp",
			wantLevel: "error",
			wantErr:   false,
		},
		{
			name:    "invalid port",
			file:    "Port=abc",
			env:     map[string]string{},
			wantErr: true,
		},
		{
			name:    "invalid port from env",
			file:    "",
			env:     map[string]string{"APP_PORT": "0"},
			wantErr: true,
		},
		{
			name:      "env empty string no override",
			file:      "Name=fromfile",
			env:       map[string]string{"APP_NAME": ""},
			wantPort:  8080,
			wantName:  "fromfile",
			wantLevel: "info",
			wantErr:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg, err := LoadConfig(test.file, test.env)
			if test.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			var sb strings.Builder
			if cfg.Name != test.wantName {
				sb.WriteString(fmt.Sprintf("Name=%q want=%q ", cfg.Name, test.wantName))
			}
			if cfg.Port != test.wantPort {
				sb.WriteString(fmt.Sprintf("Port=%d want=%d ", cfg.Port, test.wantPort))
			}
			if cfg.LogLevel != test.wantLevel {
				sb.WriteString(fmt.Sprintf("LogLevel=%q want=%q ", cfg.LogLevel, test.wantLevel))
			}
			if sb.Len() > 0 {
				t.Fatalf("LoadConfig() mismatch: %s", sb.String())
			}
		})
	}
}
