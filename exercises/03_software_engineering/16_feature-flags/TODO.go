//go:build todo

package featureflags

import (
	"errors"
	"fmt"
	"sync"
)

// Variant represents a flag assignment variant.
// TODO: implement Variant type and constants.
type Variant string

const (
	Control   Variant = "control"
	Treatment Variant = "treatment"
)

// IsValid reports whether v is a recognized variant.
// TODO: implement Variant.IsValid.
func (v Variant) IsValid() bool {
	return false
}

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
// TODO: implement validateFlags — check Percentage range and override validity.
func validateFlags(flags map[string]Flag) error {
	return nil
}

// copyFlags creates a deep copy of the flags map.
// TODO: implement copyFlags — deep-copy both the map and each Flag's Overrides.
func copyFlags(flags map[string]Flag) map[string]Flag {
	return nil
}

// hashBucket returns a deterministic bucket 0-99 from the flag key and user ID.
// TODO: implement hashBucket using crypto/fnv.New32a().
// The hash must include both the flag key and user ID so that
// assignments are stable and orthogonal across flags.
func hashBucket(flagKey, userID string) int {
	return 0
}

// NewEvaluator creates a new evaluator with the given flags.
// It defensively copies the flags map so caller mutations have no effect.
// TODO: implement NewEvaluator — validate and deep-copy flags.
func NewEvaluator(flags map[string]Flag) (*Evaluator, error) {
	return nil, errors.New("not implemented")
}

// Replace atomically swaps the evaluator's configuration.
// If validation fails, the previous configuration is preserved.
// TODO: implement Replace — validate first, then swap under lock.
func (e *Evaluator) Replace(flags map[string]Flag) error {
	errors.New("not implemented")
	return nil // placeholder
}

// Evaluate returns the evaluation result for the given flag and user.
// Precedence: validate inputs → unknown flag → disabled → override → rollout.
// TODO: implement Evaluate following the evaluation precedence.
func (e *Evaluator) Evaluate(flagKey, userID string) (Evaluation, error) {
	fmt.Println("TODO: implement Evaluate")
	return Evaluation{}, errors.New("not implemented")
}
