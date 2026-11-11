package healthreadiness

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHealthAllProbesOK(t *testing.T) {
	probes := map[string]Probe{
		"db":    func(context.Context) (Status, string) { return StatusOK, "" },
		"cache": func(context.Context) (Status, string) { return StatusOK, "" },
	}
	c := NewChecker(probes, 1*time.Millisecond)
	s, d := c.Health()
	if s != StatusOK {
		t.Errorf("Health() status = %q; want %q", s, StatusOK)
	}
	if d != "" {
		t.Errorf("Health() detail = %q; want empty", d)
	}
}

func TestHealthWorstStatusWins(t *testing.T) {
	probes := map[string]Probe{
		"db":    func(context.Context) (Status, string) { return StatusOK, "" },
		"cache": func(context.Context) (Status, string) { return StatusWarning, "slow" },
	}
	c := NewChecker(probes, 1*time.Millisecond)
	s, d := c.Health()
	if s != StatusWarning {
		t.Errorf("Health() status = %q; want %q", s, StatusWarning)
	}
	if d != "slow" {
		t.Errorf("Health() detail = %q; want %q", d, "slow")
	}
}

func TestHealthCriticalWins(t *testing.T) {
	probes := map[string]Probe{
		"db":    func(context.Context) (Status, string) { return StatusWarning, "slow" },
		"cache": func(context.Context) (Status, string) { return StatusCritical, "down" },
	}
	c := NewChecker(probes, 1*time.Millisecond)
	s, d := c.Health()
	if s != StatusCritical {
		t.Errorf("Health() status = %q; want %q", s, StatusCritical)
	}
	if d != "down" {
		t.Errorf("Health() detail = %q; want %q", d, "down")
	}
}

func TestReadyBeforeMark(t *testing.T) {
	c := NewChecker(nil, 100*time.Hour)
	s, _ := c.Ready()
	if s != StatusCritical {
		t.Errorf("Ready() status = %q; want %q (not marked ready)", s, StatusCritical)
	}
}

func TestReadyAfterMarkReady(t *testing.T) {
	c := NewChecker(nil, 1*time.Millisecond)
	c.MarkReady()
	// minReadyTime is 1ms, so we just verify MarkReady unblocks when time is sufficient
	c.startTime = time.Now().Add(-time.Second)
	s, _ := c.Ready()
	if s != StatusOK {
		t.Errorf("Ready() status = %q; want %q", s, StatusOK)
	}
}

func TestServeHTTPHealthEndpoint(t *testing.T) {
	c := NewChecker(nil, 1*time.Millisecond)
	srv := httptest.NewServer(c)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/health")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("GET /health status = %d; want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestServeHTTPReadyEndpoint(t *testing.T) {
	c := NewChecker(nil, 100*time.Hour)
	srv := httptest.NewServer(c)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/ready")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("GET /ready status = %d; want %d", resp.StatusCode, http.StatusServiceUnavailable)
	}
}

func TestServeHTTPUnknownRoute(t *testing.T) {
	c := NewChecker(nil, 1*time.Millisecond)
	srv := httptest.NewServer(c)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/unknown")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("GET /unknown status = %d; want %d", resp.StatusCode, http.StatusNotFound)
	}
}

func TestHealthEmptyProbes(t *testing.T) {
	c := NewChecker(map[string]Probe{}, 1*time.Millisecond)
	s, _ := c.Health()
	if s != StatusOK {
		t.Errorf("Health() with no probes = %q; want %q", s, StatusOK)
	}
}
