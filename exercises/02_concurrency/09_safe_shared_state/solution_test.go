package safesharedstate

import (
	"errors"
	"sync"
	"testing"
)

func mustLedger(t *testing.T, initial map[string]int64) *Ledger {
	t.Helper()

	ledger, err := NewLedger(initial)
	if err != nil {
		t.Fatalf("NewLedger() error = %v; want nil", err)
	}
	return ledger
}

func TestNewLedgerCopiesInitialState(t *testing.T) {
	initial := map[string]int64{"alice": 100}
	ledger := mustLedger(t, initial)

	initial["alice"] = 999
	initial["bob"] = 50

	got := ledger.Snapshot()
	want := map[string]int64{"alice": 100}
	if !mapsEqual(got, want) {
		t.Fatalf("Snapshot() = %v; want %v", got, want)
	}
}

func TestNewLedgerRejectsInvalidInitialState(t *testing.T) {
	tests := []struct {
		name    string
		initial map[string]int64
		wantErr error
	}{
		{name: "empty account", initial: map[string]int64{"": 10}, wantErr: ErrEmptyAccount},
		{name: "negative balance", initial: map[string]int64{"alice": -1}, wantErr: ErrNegativeBalance},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ledger, err := NewLedger(test.initial)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("NewLedger() error = %v; want %v", err, test.wantErr)
			}
			if ledger != nil {
				t.Fatalf("NewLedger() = %#v; want nil", ledger)
			}
		})
	}
}

func TestBalanceAndDeposit(t *testing.T) {
	ledger := mustLedger(t, map[string]int64{"alice": 100})

	if err := ledger.Deposit("alice", 25); err != nil {
		t.Fatalf("Deposit() error = %v; want nil", err)
	}
	got, err := ledger.Balance("alice")
	if err != nil {
		t.Fatalf("Balance() error = %v; want nil", err)
	}
	if got != 125 {
		t.Fatalf("Balance() = %d; want 125", got)
	}
}

func TestBalanceValidation(t *testing.T) {
	ledger := mustLedger(t, map[string]int64{"alice": 100})

	if _, err := ledger.Balance(""); !errors.Is(err, ErrEmptyAccount) {
		t.Fatalf("Balance(empty) error = %v; want ErrEmptyAccount", err)
	}
	if _, err := ledger.Balance("bob"); !errors.Is(err, ErrAccountNotFound) {
		t.Fatalf("Balance(missing) error = %v; want ErrAccountNotFound", err)
	}
}

func TestDepositValidationDoesNotMutateState(t *testing.T) {
	ledger := mustLedger(t, map[string]int64{"alice": 100})
	tests := []struct {
		name    string
		account string
		amount  int64
		wantErr error
	}{
		{name: "empty account", account: "", amount: 1, wantErr: ErrEmptyAccount},
		{name: "zero amount", account: "alice", amount: 0, wantErr: ErrInvalidAmount},
		{name: "negative amount", account: "alice", amount: -1, wantErr: ErrInvalidAmount},
		{name: "missing account", account: "bob", amount: 1, wantErr: ErrAccountNotFound},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := ledger.Deposit(test.account, test.amount); !errors.Is(err, test.wantErr) {
				t.Fatalf("Deposit() error = %v; want %v", err, test.wantErr)
			}
			if got := ledger.Total(); got != 100 {
				t.Fatalf("Total() = %d after rejected deposit; want 100", got)
			}
		})
	}
}

func TestTransfer(t *testing.T) {
	ledger := mustLedger(t, map[string]int64{"alice": 100, "bob": 40})

	if err := ledger.Transfer("alice", "bob", 30); err != nil {
		t.Fatalf("Transfer() error = %v; want nil", err)
	}
	want := map[string]int64{"alice": 70, "bob": 70}
	if got := ledger.Snapshot(); !mapsEqual(got, want) {
		t.Fatalf("Snapshot() = %v; want %v", got, want)
	}
	if got := ledger.Total(); got != 140 {
		t.Fatalf("Total() = %d; want 140", got)
	}
}

func TestTransferValidationIsAtomic(t *testing.T) {
	tests := []struct {
		name    string
		from    string
		to      string
		amount  int64
		wantErr error
	}{
		{name: "empty source", from: "", to: "bob", amount: 1, wantErr: ErrEmptyAccount},
		{name: "invalid amount", from: "alice", to: "bob", amount: 0, wantErr: ErrInvalidAmount},
		{name: "missing source", from: "carol", to: "bob", amount: 1, wantErr: ErrAccountNotFound},
		{name: "missing destination", from: "alice", to: "carol", amount: 1, wantErr: ErrAccountNotFound},
		{name: "insufficient funds", from: "alice", to: "bob", amount: 101, wantErr: ErrInsufficientFunds},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ledger := mustLedger(t, map[string]int64{"alice": 100, "bob": 40})
			before := ledger.Snapshot()

			err := ledger.Transfer(test.from, test.to, test.amount)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("Transfer() error = %v; want %v", err, test.wantErr)
			}
			if after := ledger.Snapshot(); !mapsEqual(after, before) {
				t.Fatalf("state changed after rejected transfer: before %v, after %v", before, after)
			}
		})
	}
}

func TestTransferToSameAccountIsNoOp(t *testing.T) {
	ledger := mustLedger(t, map[string]int64{"alice": 100})

	if err := ledger.Transfer("alice", "alice", 25); err != nil {
		t.Fatalf("Transfer() error = %v; want nil", err)
	}
	if got := ledger.Total(); got != 100 {
		t.Fatalf("Total() = %d; want 100", got)
	}
}

func TestSnapshotIsIndependent(t *testing.T) {
	ledger := mustLedger(t, map[string]int64{"alice": 100})

	snapshot := ledger.Snapshot()
	snapshot["alice"] = 0
	snapshot["bob"] = 50

	want := map[string]int64{"alice": 100}
	if got := ledger.Snapshot(); !mapsEqual(got, want) {
		t.Fatalf("Snapshot() = %v after caller mutation; want %v", got, want)
	}
}

func TestConcurrentDepositsDoNotLoseUpdates(t *testing.T) {
	ledger := mustLedger(t, map[string]int64{"alice": 0})
	const goroutines = 16
	const depositsPerGoroutine = 1_000

	start := make(chan struct{})
	errs := make(chan error, goroutines)
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for range goroutines {
		go func() {
			defer wg.Done()
			<-start
			for range depositsPerGoroutine {
				if err := ledger.Deposit("alice", 1); err != nil {
					errs <- err
					return
				}
			}
		}()
	}

	close(start)
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Fatalf("Deposit() error = %v; want nil", err)
	}

	want := int64(goroutines * depositsPerGoroutine)
	got, err := ledger.Balance("alice")
	if err != nil {
		t.Fatalf("Balance() error = %v; want nil", err)
	}
	if got != want {
		t.Fatalf("Balance() = %d; want %d", got, want)
	}
}

func TestConcurrentTransfersPreserveInvariant(t *testing.T) {
	const initialBalance int64 = 100_000
	ledger := mustLedger(t, map[string]int64{
		"alice": initialBalance,
		"bob":   initialBalance,
	})
	const goroutines = 12
	const rounds = 1_000

	start := make(chan struct{})
	errs := make(chan error, goroutines)
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for range goroutines {
		go func() {
			defer wg.Done()
			<-start
			for range rounds {
				if err := ledger.Transfer("alice", "bob", 1); err != nil {
					errs <- err
					return
				}
				if err := ledger.Transfer("bob", "alice", 1); err != nil {
					errs <- err
					return
				}
			}
		}()
	}

	close(start)
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Fatalf("Transfer() error = %v; want nil", err)
	}

	if got, want := ledger.Total(), 2*initialBalance; got != want {
		t.Fatalf("Total() = %d; want %d", got, want)
	}
	want := map[string]int64{"alice": initialBalance, "bob": initialBalance}
	if got := ledger.Snapshot(); !mapsEqual(got, want) {
		t.Fatalf("Snapshot() = %v; want %v", got, want)
	}
}

func mapsEqual(left, right map[string]int64) bool {
	if len(left) != len(right) {
		return false
	}
	for key, value := range left {
		rightValue, ok := right[key]
		if !ok || rightValue != value {
			return false
		}
	}
	return true
}
