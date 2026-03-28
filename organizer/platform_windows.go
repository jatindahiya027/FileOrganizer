//go:build windows

package organizer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

var (
	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	procSetFileAttrs = kernel32.NewProc("SetFileAttributesW")
	procGetFileAttrs = kernel32.NewProc("GetFileAttributesW")
)

const (
	fileAttrHidden = uint32(0x00000002)
	fileAttrSystem = uint32(0x00000004)
	fileAttrNormal = uint32(0x00000080)
	createNoWindow = uint32(0x08000000)
)

// MoveFile is the exported entry point used by organizer.go on Windows.
// It resolves name conflicts, then uses robocopy with no console window.
func MoveFile(src, dstDir, filename string) (string, error) {
	// Create destination folder ourselves — never let robocopy create it
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return "", fmt.Errorf("mkdir: %w", err)
	}
	stripHiddenSystem(dstDir)

	// Resolve a conflict-free filename before moving
	safeFilename := resolveConflict(dstDir, filename)
	dstPath := filepath.Join(dstDir, safeFilename)

	// If the name changed we must temporarily rename the source file so
	// robocopy moves it under the new name (robocopy has no rename flag)
	actualSrc := src
	if safeFilename != filename {
		tmpPath := filepath.Join(filepath.Dir(src), safeFilename)
		if err := os.Rename(src, tmpPath); err != nil {
			return "", fmt.Errorf("pre-rename for conflict resolution: %w", err)
		}
		actualSrc = tmpPath
		// Restore original name if robocopy fails
		defer func() {
			if _, err := os.Stat(tmpPath); err == nil {
				os.Rename(tmpPath, src)
			}
		}()
	}

	srcDir := filepath.Dir(actualSrc)
	cmd := exec.Command("robocopy",
		srcDir,
		dstDir,
		safeFilename,
		"/MOV",
		"/COPY:DAT",
		"/DCOPY:T",
		"/IS",
		"/IT",
		"/R:3",
		"/W:1",
		"/NJH",
		"/NJS",
		"/NFL",
		"/NDL",
		"/NP",
	)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: createNoWindow,
	}

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			code := exitErr.ExitCode()
			if code <= 1 {
				stripHiddenSystem(dstDir)
				return dstPath, nil
			}
			return "", fmt.Errorf("robocopy exit code %d", code)
		}
		return "", fmt.Errorf("robocopy: %w", err)
	}

	stripHiddenSystem(dstDir)
	return dstPath, nil
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

// stripHiddenSystem removes Hidden and System attribute bits from a path.
func stripHiddenSystem(path string) {
	p16, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return
	}
	cur, _, _ := procGetFileAttrs.Call(uintptr(unsafe.Pointer(p16)))
	if cur == 0xFFFFFFFF || cur == 0 {
		return
	}
	attrs := uint32(cur)
	if attrs&fileAttrHidden == 0 && attrs&fileAttrSystem == 0 {
		return
	}
	clean := attrs &^ (fileAttrHidden | fileAttrSystem)
	if clean == 0 {
		clean = fileAttrNormal
	}
	procSetFileAttrs.Call(uintptr(unsafe.Pointer(p16)), uintptr(clean))
}
