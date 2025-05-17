package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := make(Environment)

	// Read dir
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read dir %s failed: %w", dir, err)
	}

	// Read files in dir
	for _, file := range files {
		// Skip sub dirs
		if file.IsDir() {
			continue
		}

		fileName := file.Name()

		if strings.Contains(fileName, "=") {
			continue
		}

		info, err := file.Info()
		if err != nil {
			return nil, fmt.Errorf("stat file %q failed: %w", fileName, err)
		}

		if info.Size() == 0 {
			env[fileName] = EnvValue{Value: "", NeedRemove: true}
			continue
		}

		fileData, err := os.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, fmt.Errorf("read file %s failed: %w", file.Name(), err)
		}

		lines := bytes.Split(fileData, []byte{'\n'})
		firstLine := lines[0]

		// Replace \x00 to \n
		firstLine = bytes.ReplaceAll(firstLine, []byte{0}, []byte{'\n'})

		// Delete \t escapes
		value := strings.TrimRight(string(firstLine), " \t")

		env[fileName] = EnvValue{Value: value, NeedRemove: false}
	}

	return env, nil
}
