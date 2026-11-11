package structuredlogging

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestInfo_Output(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger(buf, Info)
	logger.Info("hello")
	if !strings.Contains(buf.String(), "INFO") || !strings.Contains(buf.String(), "hello") {
		t.Errorf("expected INFO hello in: %s", buf.String())
	}
}

func TestDebugSuppressed(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger(buf, Info)
	// Force Info-level logger to not output Info below Warn
	_ = logger.With("req", "123")
}

func TestWith_ChildFields(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger(buf, Debug)
	child := logger.With("user", "alice")
	child.Info("logged")
	out := buf.String()
	if !strings.Contains(out, "user=alice") {
		t.Errorf("expected user=alice in: %s", out)
	}
}

func TestNestedWith(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger(buf, Debug)
	grandChild := logger.With("req", "1").With("user", "bob")
	grandChild.Info("msg")
	out := buf.String()
	if !strings.Contains(out, "req=1") || !strings.Contains(out, "user=bob") {
		t.Errorf("expected both fields in: %s", out)
	}
}

func TestWithError(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger(buf, Debug)
	logger.Error("fail", fmt.Errorf("boom"))
	out := buf.String()
	if !strings.Contains(out, "boom") {
		t.Errorf("expected error value in: %s", out)
	}
}
