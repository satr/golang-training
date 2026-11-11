//go:build todo

package safesharedstate

import (
	"errors"
	"sync"
)

var (
	ErrEmptyAccount      = errors.New("account name must not be empty")
	ErrInvalidAmount     = errors.New("amount must be positive")
	ErrNegativeBalance   = errors.New("initial balance must not be negative")
	ErrAccountNotFound   = errors.New("account not found")
	ErrInsufficientFunds = errors.New("insufficient funds")
)

// Ledger stores account balances safely for concurrent use.
type Ledger struct {
	mu       sync.RWMutex
	balances map[string]int64
}

// NewLedger validates and copies initial.
func NewLedger(initial map[string]int64) (*Ledger, error) {
	panic("TODO: implement NewLedger")
}

// Balance returns the current balance of account.
func (l *Ledger) Balance(account string) (int64, error) {
	panic("TODO: implement Balance")
}

// Deposit atomically adds amount to account.
func (l *Ledger) Deposit(account string, amount int64) error {
	panic("TODO: implement Deposit")
}

// Transfer atomically moves amount from one account to another.
func (l *Ledger) Transfer(from, to string, amount int64) error {
	panic("TODO: implement Transfer")
}

// Snapshot returns an independent copy of all balances.
func (l *Ledger) Snapshot() map[string]int64 {
	panic("TODO: implement Snapshot")
}

// Total returns the sum of all balances from one consistent view.
func (l *Ledger) Total() int64 {
	panic("TODO: implement Total")
}
