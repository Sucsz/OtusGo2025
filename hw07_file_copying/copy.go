package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	// Progress-Bar
	//nolint:depguard
	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Open
	file, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnsupportedFile, fromPath)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("error closing file: %s\n", err)
		}
	}(file)

	// Validate && Do offset && Calc limit
	stat, err := os.Stat(fromPath)
	if err != nil {
		return fmt.Errorf("stat failed: %w", err)
	}
	if !stat.Mode().IsRegular() {
		return ErrUnsupportedFile
	}
	if offset > stat.Size() {
		return ErrOffsetExceedsFileSize
	}

	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		return fmt.Errorf("seek failed in %s: %w", fromPath, err)
	}

	if limit == 0 {
		limit = stat.Size() - offset
	} else {
		limit = min(limit, stat.Size()-offset)
	}

	// Create new file for copy
	copyFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("create file %s failed: %w", toPath, err)
	}
	defer func(copyFile *os.File) {
		err := copyFile.Close()
		if err != nil {
			fmt.Printf("error closing copy file: %s\n", err)
		}
	}(copyFile)

	// Create progress bar
	bar := pb.Full.Start64(limit)
	// End progress bar
	defer bar.Finish()
	// Show progress with Bytes
	bar.Set(pb.Bytes, true)
	barReader := bar.NewProxyReader(file)

	// Copy
	amount, err := io.CopyN(copyFile, barReader, limit)
	if err != nil && !(errors.Is(err, io.EOF)) {
		return fmt.Errorf("copy file failed: %w", err)
	}

	fmt.Printf("Copied %d bytes\n", amount)

	return nil
}
