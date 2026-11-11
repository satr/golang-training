package featureflags

import (
	"errors"
	"fmt"
	"hash/fnv"
	"sync"
)

// Variant represents a flag assignment variant.
type Variant string

const (
	Control   Variant = "control"
	Treatment Variant = "treatment"
)

// IsValid reports whether v is a recognized variant.
func (v Variant) IsValid() bool { return v == Control || v == Treatment }

// Reason explains how an evaluation was decided.
type Reason string

const (
	ReasonDisabled Reason = "disabled"
	ReasonOverride Reason = "override"
	ReasonRollout  Reason = "rollout"
)

// Flag describes one feature flag's configuration.
type Flag struct {
	Enabled    bool
	Percentage int
	Overrides  map[string]Variant
}

// Evaluation is the result of evaluating a flag for a user.
type Evaluation struct {
	Variant Variant
	Reason  Reason
}

// Sentinel errors.
var (
	ErrEmptyFlagKey  = errors.New("flag key is empty")
	ErrEmptyUserID   = errors.New("user ID is empty")
	ErrUnknownFlag   = errors.New("flag not found")
	ErrBadPercentage = errors.New("percentage must be between 0 and 100")
	ErrBadOverride   = errors.New("override variant must be control or treatment")
)

// Evaluator holds the flags configuration and provides thread-safe evaluation.
type Evaluator struct {
	mu    sync.RWMutex
	flags map[string]Flag
}

// validateFlags checks every flag in the map for configuration errors.
func validateFlags(flags map[string]Flag) error {
	for key, f := range flags {
		if f.Percentage < 0 || f.Percentage > 100 {
			return fmt.Errorf("flag %q: %w", key, ErrBadPercentage)
		}
		if f.Overrides != nil {
			for uid, v := range f.Overrides {
				if !v.IsValid() {
					return fmt.Errorf("flag %q override %q: %w", key, uid, ErrBadOverride)
				}
			}
		}
	}
	return nil
}

// copyFlags creates a deep copy of the flags map.
func copyFlags(flags map[string]Flag) map[string]Flag {
	cp := make(map[string]Flag, len(flags))
	for k, f := range flags {
		flagCopy := f
		if f.Overrides != nil {
			flagCopy.Overrides = make(map[string]Variant, len(f.Overrides))
			for uid, v := range f.Overrides {
				flagCopy.Overrides[uid] = v
			}
		}
		cp[k] = flagCopy
	}
	return cp
}

// hashBucket returns a deterministic bucket 0-99 from the flag key and user ID.
func hashBucket(flagKey, userID string) int {
	h := fnv.New32a()
	h.Write([]byte(flagKey))
	h.Write([]byte("|"))
	h.Write([]byte(userID))
	return int(h.Sum32() % 100)
}

// NewEvaluator creates a new evaluator with the given flags.
// It defensively copies the flags map so caller mutations have no effect.
func NewEvaluator(flags map[string]Flag) (*Evaluator, error) {
	if err := validateFlags(flags); err != nil {
		return nil, err
	}
	return &Evaluator{flags: copyFlags(flags)}, nil
}

// Replace atomically swaps the evaluator's configuration.
// If validation fails, the previous configuration is preserved.
func (e *Evaluator) Replace(flags map[string]Flag) error {
	if err := validateFlags(flags); err != nil {
		return err
	}
	e.mu.Lock()
	e.flags = copyFlags(flags)
	e.mu.Unlock()
	return nil
}

// Evaluate returns the evaluation result for the given flag and user.
func (e *Evaluator) Evaluate(flagKey, userID string) (Evaluation, error) {
	if flagKey == "" {
		return Evaluation{}, ErrEmptyFlagKey
	}
	if userID == "" {
		return Evaluation{}, ErrEmptyUserID
	}

	e.mu.RLock()
	flag, ok := e.flags[flagKey]
	e.mu.RUnlock()

	if !ok {
		return Evaluation{}, fmt.Errorf("%q: %w", flagKey, ErrUnknownFlag)
	}

	if !flag.Enabled {
		return Evaluation{Variant: Control, Reason: ReasonDisabled}, nil
	}

	if v, ok := flag.Overrides[userID]; ok {
		return Evaluation{Variant: v, Reason: ReasonOverride}, nil
	}

	bucket := hashBucket(flagKey, userID)
	if bucket < flag.Percentage {
		return Evaluation{Variant: Treatment, Reason: ReasonRollout}, nil
	}
	return Evaluation{Variant: Control, Reason: ReasonRollout}, nil
}
