package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: go-envdir <envdir> <command> [args...]")
		os.Exit(1)
	}

	envDir := os.Args[1]
	cmdArgs := os.Args[2:]

	env, err := ReadDir(envDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read env dir: %v\n", err)
		os.Exit(1)
	}

	code := RunCmd(cmdArgs, env)
	os.Exit(code)
}
