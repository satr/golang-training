package featureflags

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

// evalOK wraps Evaluate and records any failure to the testing.TB.
func evalOK(t *testing.T, e *Evaluator, flagKey, userID string) Evaluation {
	t.Helper()
	got, err := e.Evaluate(flagKey, userID)
	if err != nil {
		t.Fatalf("Evaluate(%q, %q): %v", flagKey, userID, err)
	}
	return got
}

// makeFlags returns a simple flags map for most tests.
func makeFlags() map[string]Flag {
	return map[string]Flag{
		"dark-mode": {
			Enabled:    true,
			Percentage: 50,
			Overrides:  nil,
		},
		"disabled-flag": {
			Enabled:    false,
			Percentage: 0,
			Overrides:  nil,
		},
	}
}

// --- NewEvaluator ---

func TestNewValidatorErrors(t *testing.T) {
	tests := []struct {
		name    string
		flags   map[string]Flag
		wantErr error
	}{
		{
			name: "negative percentage",
			flags: map[string]Flag{
				"bad": {Enabled: true, Percentage: -1},
			},
			wantErr: ErrBadPercentage,
		},
		{
			name: "percentage over 100",
			flags: map[string]Flag{
				"bad": {Enabled: true, Percentage: 101},
			},
			wantErr: ErrBadPercentage,
		},
		{
			name: "invalid override variant",
			flags: map[string]Flag{
				"flag": {
					Enabled:    true,
					Percentage: 50,
					Overrides:  map[string]Variant{"u1": "invalid"},
				},
			},
			wantErr: ErrBadOverride,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewEvaluator(tt.flags)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("expected error wrapping %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestNewEvaluator_OK(t *testing.T) {
	flags := makeFlags()
	e, err := NewEvaluator(flags)
	if err != nil {
		t.Fatalf("NewEvaluator: %v", err)
	}
	if e == nil {
		t.Fatal("expected non-nil evaluator")
	}
}

// --- Evaluate: disabled flag, rollout, overrides ---

func TestEvaluate(t *testing.T) {
	e, err := NewEvaluator(makeFlags())
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		// setFlag optionally modifies the evaluator before evaluation.
		setup      func(*Evaluator)
		flagKey    string
		userID     string
		wantVar    Variant
		wantReason Reason
		wantErr    error
	}{
		{
			name:       "disabled flag returns control + disabled",
			flagKey:    "disabled-flag",
			userID:     "user-1",
			wantVar:    Control,
			wantReason: ReasonDisabled,
		},
		{
			name: "zero percent rollout returns control + rollout",
			setup: func(e *Evaluator) {
				e.Replace(map[string]Flag{
					"zero": {Enabled: true, Percentage: 0},
				})
			},
			flagKey:    "zero",
			userID:     "user-1",
			wantVar:    Control,
			wantReason: ReasonRollout,
		},
		{
			name: "hundred percent rollout returns treatment + rollout",
			setup: func(e *Evaluator) {
				e.Replace(map[string]Flag{
					"all": {Enabled: true, Percentage: 100},
				})
			},
			flagKey:    "all",
			userID:     "user-1",
			wantVar:    Treatment,
			wantReason: ReasonRollout,
		},
		{
			name:    "unknown flag returns error",
			flagKey: "nonexistent",
			userID:  "user-1",
			wantErr: ErrUnknownFlag,
		},
		{
			name:    "empty flag key returns error",
			flagKey: "",
			userID:  "user-1",
			wantErr: ErrEmptyFlagKey,
		},
		{
			name:    "empty user ID returns error",
			flagKey: "dark-mode",
			userID:  "",
			wantErr: ErrEmptyUserID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(e)
			}

			got, err := e.Evaluate(tt.flagKey, tt.userID)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error wrapping %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Variant != tt.wantVar {
				t.Errorf("variant = %v, want %v", got.Variant, tt.wantVar)
			}
			if got.Reason != tt.wantReason {
				t.Errorf("reason = %v, want %v", got.Reason, tt.wantReason)
			}
		})
	}
}

// TestStableAssignment verifies that the same user+flag always yields the
// same variant, regardless of how many times Evaluate is called.
func TestStableAssignment(t *testing.T) {
	e, err := NewEvaluator(makeFlags())
	if err != nil {
		t.Fatal(err)
	}

	const flagKey = "dark-mode"
	const userID = "stable-user"

	var first Evaluation
	for i := 0; i < 10; i++ {
		got, err := e.Evaluate(flagKey, userID)
		if err != nil {
			t.Fatalf("iteration %d: %v", i, err)
		}
		if i == 0 {
			first = got
		} else if got != first {
			t.Fatalf("iteration %d: got %v, want stable %v", i, got, first)
		}
	}
}

// --- Override tests ---

func TestOverrides(t *testing.T) {
	tests := []struct {
		name       string
		flags      map[string]Flag
		flagKey    string
		userID     string
		wantVar    Variant
		wantReason Reason
	}{
		{
			name: "control override",
			flags: map[string]Flag{
				"flag": {
					Enabled:    true,
					Percentage: 100,
					Overrides:  map[string]Variant{"user-ctrl": Control},
				},
			},
			flagKey:    "flag",
			userID:     "user-ctrl",
			wantVar:    Control,
			wantReason: ReasonOverride,
		},
		{
			name: "treatment override",
			flags: map[string]Flag{
				"flag": {
					Enabled:    true,
					Percentage: 0,
					Overrides:  map[string]Variant{"user-trt": Treatment},
				},
			},
			flagKey:    "flag",
			userID:     "user-trt",
			wantVar:    Treatment,
			wantReason: ReasonOverride,
		},
		{
			name: "override beats 100 percent rollout",
			flags: map[string]Flag{
				"flag": {
					Enabled:    true,
					Percentage: 100,
					Overrides:  map[string]Variant{"sticky": Control},
				},
			},
			flagKey:    "flag",
			userID:     "sticky",
			wantVar:    Control,
			wantReason: ReasonOverride,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := NewEvaluator(tt.flags)
			if err != nil {
				t.Fatal(err)
			}

			got := evalOK(t, e, tt.flagKey, tt.userID)
			if got.Variant != tt.wantVar {
				t.Errorf("variant = %v, want %v", got.Variant, tt.wantVar)
			}
			if got.Reason != tt.wantReason {
				t.Errorf("reason = %v, want %v", got.Reason, tt.wantReason)
			}
		})
	}
}

// --- Defensive copying ---

func TestConstructorDefensiveCopy(t *testing.T) {
	flags := map[string]Flag{
		"dark-mode": {
			Enabled:    true,
			Percentage: 50,
			Overrides:  map[string]Variant{"u1": Treatment},
		},
	}
	e, err := NewEvaluator(flags)
	if err != nil {
		t.Fatal(err)
	}

	// Mutate the original map and overrides.
	flags["dark-mode"].Overrides["u1"] = Control
	flags["new-flag"] = Flag{Enabled: true, Percentage: 0}

	// Evaluator must still see the original override.
	got := evalOK(t, e, "dark-mode", "u1")
	if got.Variant != Treatment {
		t.Errorf("variant = %v, want Treatment (defensive copy violated)", got.Variant)
	}

	// The new key must not exist.
	_, err = e.Evaluate("new-flag", "user")
	if !errors.Is(err, ErrUnknownFlag) {
		t.Errorf("expected ErrUnknownFlag for mutated key, got %v", err)
	}
}

func TestReplaceDefensiveCopy(t *testing.T) {
	e, err := NewEvaluator(map[string]Flag{
		"flag": {Enabled: true, Percentage: 50},
	})
	if err != nil {
		t.Fatal(err)
	}

	newFlags := map[string]Flag{
		"flag": {
			Enabled:    true,
			Percentage: 50,
			Overrides:  map[string]Variant{"u1": Treatment},
		},
	}
	if err := e.Replace(newFlags); err != nil {
		t.Fatal(err)
	}

	// Mutate the caller-side map.
	newFlags["flag"].Overrides["u1"] = Control
	newFlags["other"] = Flag{Enabled: true, Percentage: 0}

	got := evalOK(t, e, "flag", "u1")
	if got.Variant != Treatment {
		t.Errorf("variant = %v, want Treatment (Replace defensive copy violated)", got.Variant)
	}
	_, err = e.Evaluate("other", "user")
	if !errors.Is(err, ErrUnknownFlag) {
		t.Errorf("expected ErrUnknownFlag, got %v", err)
	}
}

// --- Replace: validation failure preserves old config ---

func TestReplaceRejectPreservesOld(t *testing.T) {
	validFlags := map[string]Flag{
		"dark-mode": {Enabled: true, Percentage: 50},
	}
	e, err := NewEvaluator(validFlags)
	if err != nil {
		t.Fatal(err)
	}

	// Verify old config works first.
	_ = evalOK(t, e, "dark-mode", "user-1")

	// Attempt replace with invalid flags.
	badFlags := map[string]Flag{
		"bad": {Enabled: true, Percentage: 200},
	}
	err = e.Replace(badFlags)
	if !errors.Is(err, ErrBadPercentage) {
		t.Fatalf("expected ErrBadPercentage, got %v", err)
	}

	// Old config must still work.
	_ = evalOK(t, e, "dark-mode", "user-1")
}

// TestReplace applies a valid replacement and verifies new behavior takes effect.
func TestReplace(t *testing.T) {
	e, err := NewEvaluator(map[string]Flag{
		"old": {Enabled: true, Percentage: 100},
	})
	if err != nil {
		t.Fatal(err)
	}

	err = e.Replace(map[string]Flag{
		"new": {Enabled: true, Percentage: 0},
	})
	if err != nil {
		t.Fatalf("Replace: %v", err)
	}

	// Old key must be gone.
	_, err = e.Evaluate("old", "user")
	if !errors.Is(err, ErrUnknownFlag) {
		t.Errorf("expected ErrUnknownFlag for old key, got %v", err)
	}

	// New key must exist and evaluate.
	got := evalOK(t, e, "new", "user")
	if got.Variant != Control {
		t.Errorf("variant = %v, want Control", got.Variant)
	}
	if got.Reason != ReasonRollout {
		t.Errorf("reason = %v, want ReasonRollout", got.Reason)
	}
}

// --- Concurrency ---

func TestConcurrentEvaluation(t *testing.T) {
	e, err := NewEvaluator(makeFlags())
	if err != nil {
		t.Fatal(err)
	}

	const goroutines = 50
	done := make(chan struct{})
	var wg sync.WaitGroup

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				t.Error(ctx.Err())
			default:
				evalOK(t, e, "dark-mode", "user-1")
			}
		}(i)
	}

	wg.Wait()
	close(done)
}

func TestConcurrentEvalAndReplace(t *testing.T) {
	e, err := NewEvaluator(map[string]Flag{
		"flag": {Enabled: true, Percentage: 50},
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	started := make(chan struct{})
	var wg sync.WaitGroup

	// Evaluators.
	for i := 0; i < 40; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-started
			select {
			case <-ctx.Done():
			default:
				evalOK(t, e, "flag", "user-1")
			}
		}()
	}

	// Replacers.
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			<-started
			select {
			case <-ctx.Done():
			default:
				// Build alternating configs; validation always passes.
				_ = e.Replace(map[string]Flag{
					"flag": {Enabled: true, Percentage: n % 100},
				})
			}
		}(i)
	}

	close(started) // release all goroutines simultaneously
	wg.Wait()
}
