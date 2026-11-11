# Interfaces and test doubles

## Theory

Go interfaces are structural—any type that implements the required methods
satisfies the interface automatically. This implicit satisfaction makes it
easy to write lightweight test doubles that stand in for real dependencies
without framework overhead.

A test double implements the same interface as production code but provides
predictable behavior. It might forward calls to known data, record interactions,
or simulate errors. Because interfaces are small and explicit, tests stay
focused on the unit under test.

```go
type FakeRepo struct {
	Data map[string]User
}

func (f *FakeRepo) Find(id string) (User, error) {
	// Return known test data.
}
```

## Task

Define a `User` struct and a `UserRepo` interface with `Find` and `Save`
methods. Implement `InMemoryRepo` as a backing store and `UserHandler` as a
service layer that uses the repository. Write tests with a custom test double
so that `GetUser` and `UpdateUser` can be verified in isolation. Handle
context cancellation gracefully.

## Run

```bash
make 03_software_engineering_03_interfaces_and_test_doubles
```

## Further practice

Add a `Delete` method to `UserRepo` and an `IsDeleted` field to the test double
that tracks which IDs were removed.
