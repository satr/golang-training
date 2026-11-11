# Package Architecture

## Theory

Go apps should structure packages by feature, not by technical layer. Dependencies point inward: the core domain depends on interfaces, while external packages implement them. This is the dependency inversion principle — high-level policy depends on abstractions, low-level details implement them.

## Task

Build a `Store` interface in the domain layer and a concrete `InMemoryStore` that implements it. Wire them together demonstrating that the Store interface hides the concrete implementation from domain logic.

## Run

```bash
make 03_software_engineering_09_package_architecture
```

## Further practice

Add a file-based or database-backed store, apply the same interface.
