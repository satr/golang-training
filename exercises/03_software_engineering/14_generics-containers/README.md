# Generics: Containers

## Theory

Go 1.18 introduced type parameters with constraints. The `any` alias stands for `interface{}`. The `comparable` constraint limits types to values that support `==`. Generic containers let you write type-safe collections without code duplication.

## Task

Implement a generic Stack[T] with Push, Pop, Peek, Size, and IsEmpty methods. Add a generic Map[K, V] with Set and Get that panics on nil keys.

## Run

```bash
make 03_software_engineering_14_generics_containers
```

## Further practice

Add generic filter/map operations on slices, a generic cache with TTL.
