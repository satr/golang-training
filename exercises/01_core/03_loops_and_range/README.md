# Loops and range

## Theory

Go has one loop keyword: `for`. The `range` form yields an index and value for
slices. When ranging over a string, it decodes UTF-8 and yields byte offsets
and Unicode code points (`rune`), so it is safer than indexing bytes when
counting characters.

```go
for _, value := range values {
    // use value
}
```

A nil slice and an empty slice both have length zero and are safe to range
over.

## Task

Implement `Sum` by ranging over all integers. Implement `RuneCounts` so it
counts Unicode code points, not bytes, and returns an empty non-nil map for an
empty string.

## Run

```bash
make 01_core_03_loops_and_range
```

## Further practice

Return rune counts in first-seen order without relying on randomized map
iteration.

