# Race Detection and Deterministic Concurrency Tests

This exercise implements a concurrent memoization cache. Concurrent callers
requesting the same key share one in-flight load, while loaders for different
keys may run at the same time.

The implementation practises race prevention. The tests practise making
concurrency observable with channels and barriers instead of hoping that a
short `time.Sleep` creates the desired interleaving.

## A race is about unsynchronized memory access

A data race occurs when goroutines access the same memory concurrently, at
least one access is a write, and the accesses are not synchronized. This
apparently simple cache is unsafe:

```go
if value, ok := values[key]; ok {
	return value
}
value := load(key)
values[key] = value
```

Concurrent map reads and writes can race or panic. Adding a mutex around each
map access removes that race, but does not guarantee that two callers will not
both observe a miss and execute the same expensive load. That is a logical
concurrency bug.

## Result publication with a closed channel

Each cache entry has a completion channel:

```go
type entry struct {
	ready chan struct{}
	value string
	err   error
}
```

The first caller installs an in-flight entry while holding the mutex, releases
the mutex, and runs the loader. Same-key callers find that entry and wait for
`ready` to close.

The first caller stores `value` and `err` before closing `ready`. Closing a
channel synchronizes with receivers unblocked by that close, so waiters can
safely read the published result afterward.

Do not hold the cache mutex while running the loader. Loader duration is
unbounded, it may call other code, and holding the global lock would serialize
loads for unrelated keys.

## API to implement

```go
type Loader func(key string) (string, error)

type Cache struct {
	mu      sync.Mutex
	entries map[string]*entry
}

func (c *Cache) Get(key string, load Loader) (string, error)
func (c *Cache) Delete(key string) bool
func (c *Cache) Len() int
```

The zero value of `Cache` must be usable.

### `Get`

- Reject an empty key with `ErrEmptyKey`.
- Reject a nil loader with `ErrNilLoader`.
- Return a completed cached value without running the loader again.
- Let exactly one caller load a missing key at a time.
- Make concurrent same-key callers wait for the in-flight result.
- Do not serialize loaders for different keys.
- Wrap and return loader errors to current callers.
- Remove failed entries so a later call can retry.

Panics from a loader are outside this exercise's contract. In production, a
panic policy would be required so waiters are not left blocked forever.

### `Delete`

Delete a completed successful entry and report `true`. Return `false` for an
empty key, a missing key, or an in-flight load. Do not interrupt or replace an
in-flight load.

### `Len`

Return the number of completed successful values. Do not count in-flight
loads. It is a point-in-time observation; another goroutine may change the
cache immediately after it returns.

## Deterministic concurrency tests

Use synchronization to create the state the test needs:

```go
started := make(chan struct{})
release := make(chan struct{})

load := func(string) (string, error) {
	close(started) // observable: loader entered
	<-release      // controlled: loader remains in flight
	return "value", nil
}
```

The test waits on `started`, inspects in-flight behavior, and then closes
`release`. Channel close broadcasts to all receivers and makes a useful start
barrier. `sync.WaitGroup` proves completion.

Short sleeps do not prove that a goroutine reached a particular line. They
make tests dependent on scheduler speed and machine load. The tests use
timeouts only as generous guards against a hang, never as the mechanism that
creates correctness.

## Using the race detector

Run:

```bash
go test -race ./exercises/02_concurrency/10_race_detection
```

The detector instruments memory accesses and reports conflicting accesses
that occur during the run. A report includes the access stacks and goroutine
creation stacks. Start from the first conflicting accesses, not merely the
last line printed.

Useful options include:

```bash
go test -race -count=10 ./exercises/02_concurrency/10_race_detection
GORACE="halt_on_error=1" go test -race ./exercises/02_concurrency/10_race_detection
```

The detector is dynamic: it cannot report an execution path that tests never
exercise. Passing `-race` proves neither functional correctness nor absence of
all possible races. Combine it with deterministic overlap and assertions about
semantic properties such as one load per key.

## Exercise workflow

Reset `solution.go` from the template:

```bash
make reset EXERCISE=exercises/02_concurrency/10_race_detection
```

Implement the TODOs and run:

```bash
make 02_concurrency_10_race_detection
go test -race ./exercises/02_concurrency/10_race_detection
```

`TODO.go` uses `//go:build todo` because it is the reset template. Do not
enable that build tag while the completed `solution.go` is also present.

## Common mistakes

- Accessing the entries map outside the mutex.
- Holding the mutex while executing a loader.
- Unlocking between checking for an entry and installing the in-flight entry.
- Closing `ready` before storing the result.
- Closing an entry's channel more than once.
- Caching failures permanently when retries are required.
- Deleting an in-flight entry and allowing a duplicate loader to start.
- Treating a clean race-detector run as proof of logical correctness.
- Using sleeps to force an ordering in tests.
- Copying a `Cache` after its mutex has been used.

## Further practice

- Add context-aware waiting without cancelling a load needed by other callers.
- Add expiry and decide how stale entries interact with in-flight refreshes.
- Recover loader panics and release all waiters with a defined error.
- Benchmark a single lock against sharded locks.
- Compare the design with `golang.org/x/sync/singleflight`.
