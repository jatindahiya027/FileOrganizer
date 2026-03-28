//go:build !windows

package organizer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// MoveFile is the exported entry point used by organizer.go on macOS/Linux.
// It resolves name conflicts then moves via os.Rename or copy+delete.
func MoveFile(src, dstDir, filename string) (string, error) {
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return "", fmt.Errorf("mkdir: %w", err)
	}

	safeFilename := resolveConflict(dstDir, filename)
	dst := filepath.Join(dstDir, safeFilename)

	return dst, moveUnix(src, dst)
}

// resolveConflict returns a filename that does not already exist in dstDir.
// "photo.jpg" → "photo_(1).jpg" → "photo_(2).jpg" …
func resolveConflict(dstDir, filename string) string {
	if _, err := os.Stat(filepath.Join(dstDir, filename)); os.IsNotExist(err) {
		return filename
	}
	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)
	for i := 1; i < 100000; i++ {
		newName := fmt.Sprintf("%s_(%d)%s", base, i, ext)
		if _, err := os.Stat(filepath.Join(dstDir, newName)); os.IsNotExist(err) {
			return newName
		}
	}
	return fmt.Sprintf("%s_(%d)%s", base, os.Getpid(), ext)
}

func moveUnix(src, dst string) error {
	if err := os.Rename(src, dst); err == nil {
		return nil
	}
	return copyThenDelete(src, dst)
}

func copyThenDelete(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open src: %w", err)
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("stat src: %w", err)
	}

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("create dst: %w", err)
	}

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		dstFile.Close()
		os.Remove(dst)
		return fmt.Errorf("copy: %w", err)
	}

	dstFile.Close()
	return os.Remove(src)
}
