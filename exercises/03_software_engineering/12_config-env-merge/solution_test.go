package configenvmerge

import (
	"os"
	"strings"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig("APP")
	v, ok := cfg.Get("host")
	if ok {
		t.Error("expected missing key")
	}
	if v != "" {
		t.Errorf("got %q, want empty", v)
	}
}

func TestWithDefaults(t *testing.T) {
	cfg := DefaultConfig("APP")
	WithDefaults(cfg, map[string]string{"host": "localhost", "port": "8080"})
	host, ok := cfg.Get("host")
	if !ok || host != "localhost" {
		t.Errorf("host = %q, %v; want localhost, true", host, ok)
	}
}

func TestWithEnv_Override(t *testing.T) {
	os.Setenv("APP_PORT", "9090")
	defer os.Unsetenv("APP_PORT")

	cfg := DefaultConfig("APP")
	WithDefaults(cfg, map[string]string{"port": "8080"})
	WithEnv(cfg)

	port, ok := cfg.Get("port")
	if !ok || port != "9090" {
		t.Errorf("port = %q, %v; want 9090, true (env override)", port, ok)
	}
}

func TestWithEnv_IgnoreWrongPrefix(t *testing.T) {
	os.Setenv("OTHER_VAR", "x")
	defer os.Unsetenv("OTHER_VAR")

	cfg := DefaultConfig("APP")
	WithEnv(cfg)
	_, ok := cfg.Get("VAR")
	if ok {
		t.Error("should not pick up non-matching env var")
	}
}

func TestString_Output(t *testing.T) {
	cfg := DefaultConfig("APP")
	WithDefaults(cfg, map[string]string{"host": "local"})
	out := cfg.String()
	if !strings.Contains(out, "host=local") {
		t.Errorf("expected host=local in %q", out)
	}
}
