//go:build todo

package filesandio

import (
	"io"
	"os"
)

// CopyFile copies all data from the source file to the destination file.
// It should return an error if either file operation fails.
func CopyFile(src, dst string) error {
	// TODO: Open source and destination files
	// TODO: Copy data from src to dst using io.Copy or similar
	// TODO: Close both files
	return nil
}

// CountLines counts the number of lines in the given reader.
// Each line is terminated by a newline character ('\n').
func CountLines(r io.Reader) (int, error) {
	// TODO: Read data from r and count lines
	// TODO: Handle potential errors during reading
	return 0, nil
}
