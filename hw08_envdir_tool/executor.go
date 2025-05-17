package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		fmt.Fprintln(os.Stderr, "no command specified")
		return 1
	}

	envMap := make(map[string]string)

	// Get env
	for _, v := range os.Environ() {
		// Split first '=', that e.g. foo=bar=baz == "foo: bar=baz"
		keyValue := strings.SplitN(v, "=", 2)
		if len(keyValue) == 2 {
			envMap[keyValue[0]] = keyValue[1]
		}
	}

	// Change ev with local env
	for k, v := range env {
		if v.NeedRemove {
			delete(envMap, k)
		} else {
			envMap[k] = v.Value
		}
	}

	// Transform for exec
	finalEnv := make([]string, 0, len(envMap))
	for k, v := range envMap {
		finalEnv = append(finalEnv, k+"="+v)
	}

	// Create cmd
	//
	//nolint:gosec
	execCmd := exec.Command(cmd[0], cmd[1:]...)
	execCmd.Env = finalEnv

	// Bind streams
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	// Run
	if err := execCmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		fmt.Fprintf(os.Stderr, "run %q failed: %v\n", cmd[0], err)
		return 1
	}
	return 0
}
