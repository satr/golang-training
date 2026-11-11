//go:build todo

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
// TODO: Store probes, initialize startTime to time.Now().
func NewChecker(probes map[string]Probe, minReadyTime time.Duration) *Checker {
	panic("TODO")
}

// Health runs all probes and returns the worst status with the first failing detail.
// TODO: Execute all probes, aggregate to worst status.
func (c *Checker) Health() (Status, string) {
	panic("TODO")
}

// Ready checks whether the service is ready to accept traffic.
// TODO: Return unhealthy if not marked ready OR min ready time not elapsed.
func (c *Checker) Ready() (Status, string) {
	panic("TODO")
}

// MarkReady signals that the service has completed startup.
// TODO: Set c.ready = true.
func (c *Checker) MarkReady() {
	panic("TODO")
}

// ServeHTTP handles /health and /ready endpoints.
// TODO: Route GET /health to Health(), GET /ready to Ready(), 404 otherwise.
func (c *Checker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	panic("TODO")
}

// healthResponse is the JSON shape returned by health/readiness endpoints.
type healthResponse struct {
	Status Status `json:"status"`
	Detail string `json:"detail,omitempty"`
}
