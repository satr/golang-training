# Dependency injection

## Theory

Dependency injection lets code receive dependencies from outside instead of
creating them internally. In Go, constructor functions accept interfaces and
return configured receivers, making code testable.

The pattern is simple: the receiving type holds an interface, the constructor
accepts that interface, and callers supply any implementation they like. Tests
pass lightweight fakes; production code passes the real implementation.

```go
type Service struct {
	storage Storage
}

func NewService(s Storage) *Service {
	return &Service{storage: s}
}
```

## Task

Define a `Storage` interface with `Get`, `Set`, and `Delete` methods.
Implement `MapStore` as an in-memory store backed by `map[string]string` and
protected by a `sync.RWMutex`. Implement `Service` that wraps the storage:
`Process` reads the value, uppercases it, and writes it back; `Get` and
`Delete` delegate through to storage. `NewService` must reject `nil` storage.

## Run

```bash
make 03_software_engineering_02_dependency_injection
```

## Further practice

Add an additional interface method `Keys() []string` and a `Clear` method to
Service that iterates storage and removes every entry.
