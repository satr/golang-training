//go:build ignore

package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	trackPattern     = regexp.MustCompile(`^(\d{2})_[a-z][a-z0-9_]*$`)
	exercisePattern  = regexp.MustCompile(`^(\d{2})_[a-z][a-z0-9_-]*$`)
	requiredFiles    = []string{"README.md", "TODO.go", "doc.go", "solution.go", "solution_test.go"}
	requiredHeadings = []string{
		"## Theory",
		"## Task",
		"## Run",
		"## Further practice",
	}
)

func main() {
	if err := verify("exercises"); err != nil {
		fmt.Fprintln(os.Stderr, "layout verification failed:", err)
		os.Exit(1)
	}
	fmt.Println("layout verification passed")
}

func verify(root string) error {
	tracks, err := childDirs(root)
	if err != nil {
		return err
	}
	if err := verifyNumberedNames(tracks, trackPattern, "track"); err != nil {
		return err
	}

	for _, track := range tracks {
		trackPath := filepath.Join(root, track)
		exercises, err := childDirs(trackPath)
		if err != nil {
			return err
		}
		if err := verifyNumberedNames(exercises, exercisePattern, track+" exercise"); err != nil {
			return err
		}
		for _, exercise := range exercises {
			if err := verifyExercise(filepath.Join(trackPath, exercise)); err != nil {
				return err
			}
		}
	}
	return rejectArtifacts(root)
}

func childDirs(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)
	return names, nil
}

func verifyNumberedNames(names []string, pattern *regexp.Regexp, kind string) error {
	for index, name := range names {
		match := pattern.FindStringSubmatch(name)
		if match == nil {
			return fmt.Errorf("%s name %q does not match %s", kind, name, pattern)
		}
		number, err := strconv.Atoi(match[1])
		if err != nil {
			return fmt.Errorf("parse number in %s: %w", name, err)
		}
		if want := index + 1; number != want {
			return fmt.Errorf("%s numbering: %q is %02d, want %02d", kind, name, number, want)
		}
	}
	return nil
}

func verifyExercise(path string) error {
	for _, name := range requiredFiles {
		info, err := os.Stat(filepath.Join(path, name))
		if err != nil {
			return fmt.Errorf("%s: required file %s: %w", path, name, err)
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("%s/%s is not a regular file", path, name)
		}
	}

	readme, err := os.ReadFile(filepath.Join(path, "README.md"))
	if err != nil {
		return fmt.Errorf("read exercise README: %w", err)
	}
	position := -1
	for _, heading := range requiredHeadings {
		next := strings.Index(string(readme), heading)
		if next <= position {
			return fmt.Errorf("%s: missing or out-of-order heading %q", path, heading)
		}
		position = next
	}

	todo, err := os.ReadFile(filepath.Join(path, "TODO.go"))
	if err != nil {
		return fmt.Errorf("read TODO template: %w", err)
	}
	if !strings.HasPrefix(string(todo), "//go:build todo\n") {
		return fmt.Errorf("%s/TODO.go must start with //go:build todo", path)
	}
	return nil
}

func rejectArtifacts(root string) error {
	generated := map[string]bool{
		".DS_Store": true, "coverage.out": true,
	}
	return filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		name := entry.Name()
		if generated[name] || strings.HasSuffix(name, ".test") ||
			strings.HasSuffix(name, ".db") || strings.HasSuffix(name, ".sqlite") {
			return errors.New("unexpected generated artifact: " + path)
		}
		return nil
	})
}
