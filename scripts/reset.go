//go:build ignore

package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var exercisePathPattern = regexp.MustCompile(
	`^exercises/\d{2}_[a-z][a-z0-9_]*/\d{2}_[a-z][a-z0-9_]*$`,
)

func main() {
	exercise := flag.String("exercise", "", "exercise directory to reset")
	all := flag.Bool("all", false, "reset every exercise")
	flag.Parse()

	if (*exercise == "") == !*all {
		exit(errors.New("choose exactly one of -exercise or -all"))
	}

	if *all {
		if err := confirmAll(os.Stdin, os.Stdout); err != nil {
			exit(err)
		}
		paths, err := allExercises("exercises")
		if err != nil {
			exit(err)
		}
		for _, path := range paths {
			if err := reset(path); err != nil {
				exit(err)
			}
		}
		fmt.Printf("reset %d exercises\n", len(paths))
		return
	}

	if err := reset(filepath.Clean(*exercise)); err != nil {
		exit(err)
	}
	fmt.Println("reset", filepath.Clean(*exercise))
}

func confirmAll(in *os.File, out *os.File) error {
	fmt.Fprint(out, "This overwrites every solution.go. Type RESET ALL to continue: ")
	line, err := bufio.NewReader(in).ReadString('\n')
	if err != nil {
		return fmt.Errorf("read confirmation: %w", err)
	}
	if strings.TrimSpace(line) != "RESET ALL" {
		return errors.New("reset cancelled")
	}
	return nil
}

func allExercises(root string) ([]string, error) {
	var paths []string
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() && exercisePathPattern.MatchString(filepath.ToSlash(path)) {
			paths = append(paths, path)
			return filepath.SkipDir
		}
		return nil
	})
	return paths, err
}

func reset(path string) error {
	slashPath := filepath.ToSlash(path)
	if filepath.IsAbs(path) || !exercisePathPattern.MatchString(slashPath) {
		return fmt.Errorf("invalid exercise directory %q", path)
	}

	templatePath := filepath.Join(path, "TODO.go")
	template, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("read %s: %w", templatePath, err)
	}
	source, err := activeSource(template)
	if err != nil {
		return fmt.Errorf("%s: %w", templatePath, err)
	}

	temp, err := os.CreateTemp(path, ".solution-*.go")
	if err != nil {
		return fmt.Errorf("create temporary solution: %w", err)
	}
	tempName := temp.Name()
	defer os.Remove(tempName)

	if _, err := temp.Write(source); err != nil {
		temp.Close()
		return fmt.Errorf("write temporary solution: %w", err)
	}
	if err := temp.Chmod(0o644); err != nil {
		temp.Close()
		return fmt.Errorf("set solution permissions: %w", err)
	}
	if err := temp.Close(); err != nil {
		return fmt.Errorf("close temporary solution: %w", err)
	}
	if err := os.Rename(tempName, filepath.Join(path, "solution.go")); err != nil {
		return fmt.Errorf("replace solution: %w", err)
	}
	return nil
}

func activeSource(template []byte) ([]byte, error) {
	const prefix = "//go:build todo\n"
	if !bytes.HasPrefix(template, []byte(prefix)) {
		return nil, errors.New("missing //go:build todo constraint")
	}
	source := bytes.TrimPrefix(template, []byte(prefix))
	source = bytes.TrimPrefix(source, []byte("\n"))
	return source, nil
}

func exit(err error) {
	fmt.Fprintln(os.Stderr, "reset failed:", err)
	os.Exit(1)
}
