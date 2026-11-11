package safesharedstate

import (
	"errors"
	"fmt"
	"sync"
)

var (
	// ErrEmptyAccount reports an empty account name.
	ErrEmptyAccount = errors.New("account name must not be empty")
	// ErrInvalidAmount reports a non-positive deposit or transfer amount.
	ErrInvalidAmount = errors.New("amount must be positive")
	// ErrNegativeBalance reports an invalid negative initial balance.
	ErrNegativeBalance = errors.New("initial balance must not be negative")
	// ErrAccountNotFound reports that an account does not exist.
	ErrAccountNotFound = errors.New("account not found")
	// ErrInsufficientFunds reports that a transfer cannot be completed.
	ErrInsufficientFunds = errors.New("insufficient funds")
)

// Ledger stores account balances safely for concurrent use.
// Its zero value is readable but contains no accounts; use NewLedger to
// create a populated ledger.
type Ledger struct {
	mu       sync.RWMutex
	balances map[string]int64
}

// NewLedger validates and copies initial. Later changes to initial do not
// affect the ledger.
func NewLedger(initial map[string]int64) (*Ledger, error) {
	balances := make(map[string]int64, len(initial))
	for account, balance := range initial {
		if account == "" {
			return nil, ErrEmptyAccount
		}
		if balance < 0 {
			return nil, fmt.Errorf("account %q: %w: %d", account, ErrNegativeBalance, balance)
		}
		balances[account] = balance
	}

	return &Ledger{balances: balances}, nil
}

// Balance returns the current balance of account.
func (l *Ledger) Balance(account string) (int64, error) {
	if account == "" {
		return 0, ErrEmptyAccount
	}

	l.mu.RLock()
	defer l.mu.RUnlock()

	balance, ok := l.balances[account]
	if !ok {
		return 0, fmt.Errorf("%w: %q", ErrAccountNotFound, account)
	}
	return balance, nil
}

// Deposit atomically adds amount to account.
func (l *Ledger) Deposit(account string, amount int64) error {
	if account == "" {
		return ErrEmptyAccount
	}
	if amount <= 0 {
		return fmt.Errorf("%w: %d", ErrInvalidAmount, amount)
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if _, ok := l.balances[account]; !ok {
		return fmt.Errorf("%w: %q", ErrAccountNotFound, account)
	}
	l.balances[account] += amount
	return nil
}

// Transfer atomically moves amount from one account to another. A transfer
// to the same existing account is a successful no-op.
func (l *Ledger) Transfer(from, to string, amount int64) error {
	if from == "" || to == "" {
		return ErrEmptyAccount
	}
	if amount <= 0 {
		return fmt.Errorf("%w: %d", ErrInvalidAmount, amount)
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	fromBalance, ok := l.balances[from]
	if !ok {
		return fmt.Errorf("source %w: %q", ErrAccountNotFound, from)
	}
	if _, ok := l.balances[to]; !ok {
		return fmt.Errorf("destination %w: %q", ErrAccountNotFound, to)
	}
	if fromBalance < amount {
		return fmt.Errorf("account %q: %w: balance %d, amount %d", from, ErrInsufficientFunds, fromBalance, amount)
	}
	if from == to {
		return nil
	}

	l.balances[from] -= amount
	l.balances[to] += amount
	return nil
}

// Snapshot returns a copy of all balances. The caller may modify the returned
// map without racing with or changing the ledger.
func (l *Ledger) Snapshot() map[string]int64 {
	l.mu.RLock()
	defer l.mu.RUnlock()

	snapshot := make(map[string]int64, len(l.balances))
	for account, balance := range l.balances {
		snapshot[account] = balance
	}
	return snapshot
}

// Total returns the sum of all balances from one consistent read-locked view.
func (l *Ledger) Total() int64 {
	l.mu.RLock()
	defer l.mu.RUnlock()

	var total int64
	for _, balance := range l.balances {
		total += balance
	}
	return total
}
