# Files and `io.Reader` / `io.Writer`

## Theory

In Go, the standard library provides the `io.Reader` and `io.Writer` interfaces to work with streams of data. These interfaces are fundamental in Go's I/O system:

- `io.Reader`: Represents a source of data that can be read from
- `io.Writer`: Represents a destination for data that can be written to

These interfaces promote composition and abstraction. Functions that take `io.Reader` or `io.Writer` arguments can work with any type implementing those interfaces, not just files.

File I/O in Go typically uses the `os` package. You can create a file with `os.Create()` or open an existing one with `os.Open()`. Both return `*os.File`, which implements both `io.Reader` and `io.Writer`.

Go also provides `io.Copy`, which is a convenient function for copying data from any `io.Reader` to any `io.Writer`.

## Task

Implement two functions:
1. `CopyFile` - copies all data from a source file to a destination file
2. `CountLines` - counts the number of lines in a given input, where each line is terminated by a newline character (`\n`)

## Run

```bash
make 01_core_11_files_and_io
```

## Further practice

- Add error handling for cases when files don't exist or can't be opened
- Write a function that counts lines in memory using `strings.Reader`
- Experiment with `io.CopyBuffer` instead of `io.Copy`