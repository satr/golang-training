package healthreadiness

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

// Status represents the health state of a service.
type Status string

const (
	StatusOK       Status = "ok"
	StatusWarning  Status = "warning"
	StatusCritical Status = "critical"
)

// severity ranks status values for comparison.
func severity(s Status) int {
	switch s {
	case StatusOK:
		return 0
	case StatusWarning:
		return 1
	case StatusCritical:
		return 2
	default:
		return 3
	}
}

// Probe checks a single dependency and returns its status.
type Probe func(context.Context) (Status, string)

// Checker aggregates multiple probes into health and readiness endpoints.
type Checker struct {
	probes       map[string]Probe
	ready        bool
	mu           sync.Mutex
	minReadyTime time.Duration
	startTime    time.Time
}

// NewChecker creates a checker with the given probes and minimum ready time.
func NewChecker(probes map[string]Probe, minReadyTime time.Duration) *Checker {
	if minReadyTime <= 0 {
		minReadyTime = time.Second
	}
	return &Checker{
		probes:       probes,
		minReadyTime: minReadyTime,
		startTime:    time.Now(),
	}
}

// Health runs all probes and returns the worst status.
func (c *Checker) Health() (Status, string) {
	var worst Status
	var detail string
	for _, probe := range c.probes {
		st, d := probe(context.Background())
		if severity(st) > severity(worst) {
			worst = st
			detail = d
		}
	}
	return worst, detail
}

// Ready checks whether the service is ready to accept traffic.
func (c *Checker) Ready() (Status, string) {
	c.mu.Lock()
	ready := c.ready
	c.mu.Unlock()

	if !ready || time.Since(c.startTime) < c.minReadyTime {
		return StatusCritical, "not ready"
	}
	return c.Health()
}

// MarkReady signals that the service has completed startup.
func (c *Checker) MarkReady() {
	c.mu.Lock()
	c.ready = true
	c.mu.Unlock()
}

// ServeHTTP handles /health and /ready endpoints.
func (c *Checker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var status Status
	var detail string

	switch r.URL.Path {
	case "/health":
		status, detail = c.Health()
	case "/ready":
		status, detail = c.Ready()
	default:
		http.NotFound(w, r)
		return
	}

	resp := healthResponse{Status: status, Detail: detail}
	if status != StatusOK {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(resp)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// healthResponse is the JSON shape returned by health/readiness endpoints.
type healthResponse struct {
	Status Status `json:"status"`
	Detail string `json:"detail,omitempty"`
}
