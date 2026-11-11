# Safe Shared State and Mutexes

This exercise uses an in-memory account ledger to practise protecting shared
state with `sync.RWMutex`. The important part is not the map itself: it is
defining which operations must be indivisible and preserving ownership at the
API boundary.

## Why a map needs protection

A Go map is not safe for concurrent reads and writes. More subtly, protecting
individual map accesses is insufficient for a transfer:

```go
if balances[from] >= amount {
	balances[from] -= amount
	balances[to] += amount
}
```

The check and both updates form one state transition. Other goroutines must not
observe or modify the ledger halfway through it. Hold one exclusive lock over
the entire transition.

## `Mutex` and `RWMutex`

`sync.Mutex` allows one holder. `sync.RWMutex` distinguishes readers from a
writer:

- `RLock` permits concurrent readers.
- `Lock` excludes every reader and writer.
- Every successful lock must have a matching unlock.

Use `RWMutex` here because `Balance`, `Snapshot`, and `Total` only read state,
while `Deposit` and `Transfer` modify it. Do not assume `RWMutex` is always
faster than `Mutex`; measure in real systems.

The mutex is part of `Ledger`, and methods use pointer receivers. Copying a
ledger after first use would copy synchronization state and is invalid.

## API to implement

```go
type Ledger struct {
	mu       sync.RWMutex
	balances map[string]int64
}

func NewLedger(initial map[string]int64) (*Ledger, error)
func (l *Ledger) Balance(account string) (int64, error)
func (l *Ledger) Deposit(account string, amount int64) error
func (l *Ledger) Transfer(from, to string, amount int64) error
func (l *Ledger) Snapshot() map[string]int64
func (l *Ledger) Total() int64
```

### Construction

`NewLedger` must:

- reject an empty account name with `ErrEmptyAccount`;
- reject a negative initial balance with an error wrapping
  `ErrNegativeBalance`;
- copy the input map so the caller does not share mutable state with the
  ledger;
- accept a nil map as an empty ledger.

### Reads

`Balance` validates the name, uses a read lock, and returns an error wrapping
`ErrAccountNotFound` for an unknown account.

`Total` calculates the sum while holding one read lock. Locking separately for
each account would not represent one consistent point in time.

`Snapshot` returns a new map populated under a read lock. Never expose the
internal map: a caller could otherwise mutate it without holding the mutex.

### Writes

`Deposit` requires a non-empty existing account and a positive amount. The
existence check and update happen under one write lock.

`Transfer` requires two non-empty existing accounts, a positive amount, and
sufficient source funds. Validation that depends on current balances and both
updates happen under the same write lock. A valid transfer to the same account
is a successful no-op.

Rejected operations must leave the ledger unchanged.

## Keep lock scope explicit

A common pattern is:

```go
l.mu.Lock()
defer l.mu.Unlock()

// inspect and modify protected fields
```

Validate arguments that do not depend on shared state before locking. Validate
current balances while locked. Avoid network calls, channel operations, disk
I/O, or caller-provided callbacks while holding the lock; their duration is
unbounded and they can introduce deadlocks.

Do not call another method that acquires the same mutex while already holding
it. Go mutexes are not reentrant. Put shared lock-assuming logic in an
unexported helper when necessary.

## Mutexes establish synchronization

Unlocking a mutex synchronizes with a later successful lock of the same mutex.
Writes performed before `Unlock` are therefore visible to the later lock
holder. A mutex protects both exclusion and memory visibility.

A race-free program still can be logically wrong. For example, locking debit
and credit separately removes the data race but allows another goroutine to
observe a temporarily reduced total. Tests therefore check both race safety
and the ledger invariant.

## Testing concurrent code

The tests use a closed start channel as a barrier so many goroutines begin
together. A `sync.WaitGroup` waits for completion. No timing sleeps are needed.

Run the package normally and with the race detector:

```bash
go test ./exercises/02_concurrency/09_safe_shared_state
go test -race ./exercises/02_concurrency/09_safe_shared_state
```

The race detector identifies conflicting unsynchronized memory accesses. It
only detects executions that actually occur, so deterministic contention and
meaningful invariant assertions are both important.

## Exercise workflow

Reset `solution.go` from the template:

```bash
make reset EXERCISE=exercises/02_concurrency/09_safe_shared_state
```

Then implement it and run the focused target:

```bash
make 02_concurrency_09_safe_shared_state
```

`TODO.go` has the `//go:build todo` constraint because it is a reset template.
Do not build with the `todo` tag while the completed `solution.go` is present;
both files declare the same API.

## Common mistakes

- Reading or writing the map without holding the appropriate lock.
- Locking only the two assignments rather than the full check-and-update.
- Returning the internal map from `Snapshot`.
- Retaining the caller's initial map instead of copying it.
- Using `RLock` in a method that mutates state.
- Forgetting an unlock on an error return.
- Calling a lock-taking method while already holding the same lock.
- Copying a `Ledger` value after its mutex has been used.
- Assuming that absence of a race report proves the state invariant.

## Further practice

- Add account creation without exposing the internal map.
- Detect integer overflow in deposits and transfers.
- Compare one ledger-wide lock with sharded account locks.
- Implement the same counter-like workload with `sync/atomic` and benchmark it.
- Add contention benchmarks for `Mutex` and `RWMutex`.
