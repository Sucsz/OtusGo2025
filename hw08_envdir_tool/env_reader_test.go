package main

import (
	"os"
	"path/filepath"
	"testing"
)

// cognitive complexity 39 of func `TestReadDir` is high (> 30) (gocognit)
//
//nolint:gocognit
func TestReadDir(t *testing.T) {
	tests := []struct {
		name      string
		files     map[string][]byte
		subDirs   []string
		want      Environment
		wantError bool
	}{
		{
			name: "valid env files",
			files: map[string][]byte{
				"FOO":       []byte("bar"),
				"EMPTY":     {},
				"MULTILINE": []byte("line1\nline2"),
				"TRAILWS":   []byte("value with spaces \t\t\nmore"),
			},
			want: Environment{
				"FOO":       {Value: "bar", NeedRemove: false},
				"EMPTY":     {Value: "", NeedRemove: true},
				"MULTILINE": {Value: "line1", NeedRemove: false},
				"TRAILWS":   {Value: "value with spaces", NeedRemove: false},
			},
		},
		{
			name: "ignore files with equals sign",
			files: map[string][]byte{
				"FOO=BAD": []byte("should be ignored"),
				"BAR":     []byte("ok"),
			},
			want: Environment{
				"BAR": {Value: "ok", NeedRemove: false},
			},
		},
		{
			name: "terminal null replaced by newline, but line ends at first newline",
			files: map[string][]byte{
				"NULL": []byte("abc\x00def"),
			},
			want: Environment{
				"NULL": {Value: "abc\ndef", NeedRemove: false},
			},
		},
		{
			name: "ignore subdirectories",
			files: map[string][]byte{
				"FOO": []byte("bar"),
			},
			subDirs: []string{"SOME_DIR"},
			want: Environment{
				"FOO": {Value: "bar", NeedRemove: false},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()

			// Create test files
			for fName, content := range tc.files {
				fullPath := filepath.Join(dir, fName)
				err := os.WriteFile(fullPath, content, 0o644)
				if err != nil {
					t.Fatalf("failed to write file %s: %v", fName, err)
				}
			}

			// Create sub dirs
			for _, d := range tc.subDirs {
				err := os.Mkdir(filepath.Join(dir, d), 0o755)
				if err != nil {
					t.Fatalf("failed to create subdir %s: %v", d, err)
				}
			}

			got, err := ReadDir(dir)
			if tc.wantError {
				if err == nil {
					t.Fatalf("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(got) != len(tc.want) {
				t.Fatalf("got %d variables, want %d", len(got), len(tc.want))
			}

			for key, wantVal := range tc.want {
				gotVal, ok := got[key]
				if !ok {
					t.Errorf("expected key %q not found", key)
					continue
				}
				if gotVal != wantVal {
					t.Errorf("key %q: got %+v, want %+v", key, gotVal, wantVal)
				}
			}
		})
	}
}
