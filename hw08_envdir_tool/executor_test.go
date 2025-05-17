package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// cognitive complexity 35 of func `TestRunCmd` is high (> 30) (gocognit)
//
//nolint:gocognit
func TestRunCmd(t *testing.T) {
	tests := []struct {
		name        string
		initialEnv  map[string]string
		overrideEnv Environment
		cmdArgs     []string
		wantOut     string
		wantCode    int
		wantAbsent  string
	}{
		{
			name:        "simple command",
			initialEnv:  map[string]string{"A": "1"},
			overrideEnv: Environment{},
			cmdArgs:     []string{"/bin/bash", "-c", "echo hello"},
			wantOut:     "hello\n",
			wantCode:    0,
		},
		{
			name:        "override variable",
			initialEnv:  map[string]string{"FOO": "old"},
			overrideEnv: Environment{"FOO": {Value: "new", NeedRemove: false}},
			cmdArgs:     []string{"/usr/bin/env"},
			wantOut:     "FOO=new\n",
			wantCode:    0,
		},
		{
			name:        "remove variable",
			initialEnv:  map[string]string{"BAR": "baz"},
			overrideEnv: Environment{"BAR": {Value: "", NeedRemove: true}},
			cmdArgs:     []string{"/usr/bin/env"},
			wantAbsent:  "BAR=",
			wantCode:    0,
		},
		{
			name:     "command not found",
			cmdArgs:  []string{"/no/such/command"},
			wantCode: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Backup and set environment variables
			backup := make(map[string]string)
			for k, v := range tc.initialEnv {
				backup[k] = os.Getenv(k)
				if err := os.Setenv(k, v); err != nil {
					t.Fatalf("failed to set env %q: %v", k, err)
				}
			}
			defer func() {
				// Restore environment variables
				for k := range tc.initialEnv {
					if old, ok := backup[k]; ok {
						_ = os.Setenv(k, old)
					} else {
						_ = os.Unsetenv(k)
					}
				}
			}()

			// Capture stdout
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("failed to create pipe: %v", err)
			}
			defer r.Close()

			oldStdout := os.Stdout
			os.Stdout = w

			code := RunCmd(tc.cmdArgs, tc.overrideEnv)

			_ = w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			if _, err := buf.ReadFrom(r); err != nil {
				t.Fatalf("failed to read from pipe: %v", err)
			}
			out := buf.String()

			if tc.wantOut != "" && !strings.Contains(out, tc.wantOut) {
				t.Errorf("%q: stdout = %q, want to contain %q", tc.name, out, tc.wantOut)
			}
			if tc.wantAbsent != "" && strings.Contains(out, tc.wantAbsent) {
				t.Errorf("%q: stdout = %q, want to NOT contain %q", tc.name, out, tc.wantAbsent)
			}
			if code != tc.wantCode {
				t.Errorf("%q: exit code = %d, want %d", tc.name, code, tc.wantCode)
			}
		})
	}
}
