package filesandio

import (
	"io"
	"os"
)

// CopyFile copies all data from the source file to the destination file.
// It should return an error if either file operation fails.
func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// CountLines counts the number of lines in the given reader.
// Each line is terminated by a newline character ('\n').
func CountLines(r io.Reader) (int, error) {
	buffer := make([]byte, 1024)
	lineCount := 0
	eof := false

	for !eof {
		n, err := r.Read(buffer)
		if err != nil {
			if err == io.EOF {
				eof = true
			} else {
				return 0, err
			}
		}

		for i := 0; i < n; i++ {
			if buffer[i] == '\n' {
				lineCount++
			}
		}
	}

	return lineCount, nil
}
