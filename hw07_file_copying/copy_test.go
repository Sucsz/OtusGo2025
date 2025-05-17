package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

//nolint:gocognit
func TestCopy(t *testing.T) {
	tmpDir := t.TempDir()

	testCases := []struct {
		name     string
		src      string
		dst      string
		offset   int64
		limit    int64
		wantFile string
		wantErr  bool
	}{
		{
			name:     "offset=0 limit=0 (full file)",
			src:      "input.txt",
			dst:      "result_full.txt",
			offset:   0,
			limit:    0,
			wantFile: "out_offset0_limit0.txt",
		},
		{
			name:     "offset=0 limit=10",
			src:      "input.txt",
			dst:      "result_0_10.txt",
			offset:   0,
			limit:    10,
			wantFile: "out_offset0_limit10.txt",
		},
		{
			name:     "offset=0 limit=1000",
			src:      "input.txt",
			dst:      "result_0_1000.txt",
			offset:   0,
			limit:    1000,
			wantFile: "out_offset0_limit1000.txt",
		},
		{
			name:     "offset=0 limit=10000",
			src:      "input.txt",
			dst:      "result_0_10000.txt",
			offset:   0,
			limit:    10000,
			wantFile: "out_offset0_limit10000.txt",
		},
		{
			name:     "offset=100 limit=1000",
			src:      "input.txt",
			dst:      "result_100_1000.txt",
			offset:   100,
			limit:    1000,
			wantFile: "out_offset100_limit1000.txt",
		},
		{
			name:     "offset=6000 limit=1000",
			src:      "input.txt",
			dst:      "result_6000_1000.txt",
			offset:   6000,
			limit:    1000,
			wantFile: "out_offset6000_limit1000.txt",
		},
		{
			name:    "offset exceeds file size",
			src:     "input.txt",
			dst:     "result_exceed.txt",
			offset:  999999,
			limit:   100,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Copy source file
			originalSrc := filepath.Join("testdata", tc.src)
			tempSrc := filepath.Join(tmpDir, "copy_"+tc.src)
			data, err := os.ReadFile(originalSrc)
			if err != nil {
				t.Fatalf("failed to read the source file: %v", err)
			}
			if err := os.WriteFile(tempSrc, data, 0o644); err != nil {
				t.Fatalf("failed to create a temporary copy of the source file: %v", err)
			}
			srcPath := tempSrc
			dstPath := filepath.Join(tmpDir, tc.dst)

			err = Copy(srcPath, dstPath, tc.offset, tc.limit)
			if tc.wantErr {
				if err == nil {
					t.Errorf("an error was expected, but no error occurred")
				}
				return
			}
			if err != nil {
				t.Fatalf("an unexpected mistake: %v", err)
			}

			// Compare
			if tc.wantFile != "" {
				gotBytes, err := os.ReadFile(dstPath)
				if err != nil {
					t.Fatalf("failed to read the result file: %v", err)
				}

				wantBytes, err := os.ReadFile(filepath.Join("testdata", tc.wantFile))
				if err != nil {
					t.Fatalf("failed to read the reference file: %v", err)
				}

				if !bytes.Equal(gotBytes, wantBytes) {
					t.Errorf("files do not match:\nreceived:\n%s\nexpected:\n%s",
						string(gotBytes), string(wantBytes))
				}
			}
		})
	}
}
