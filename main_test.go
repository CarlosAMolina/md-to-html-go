package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestMainGo(t *testing.T) {
	tmpDir := "/tmp/dir-to-convert"
	testdataSrc := "testdata/dir-to-convert"
	testdataExpected := "testdata/dir-converted"
	if err := os.RemoveAll(tmpDir); err != nil {
		t.Fatalf("failed to delete %s: %v", tmpDir, err)
	}
	if err := copyDir(testdataSrc, tmpDir); err != nil {
		t.Fatalf("failed to copy %s to %s: %v", testdataSrc, tmpDir, err)
	}
	cmd := exec.Command("go", "run", ".", tmpDir)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to execute main.go: %v\nStderr: %s", err, stderr.String())
	}
	if err := compareDirs(testdataExpected, tmpDir); err != nil {
		t.Fatalf("directories do not match: %v", err)
	}
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}
		return copyFile(path, target)
	})
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, info.Mode())
}

func compareDirs(expected, actual string) error {
	err := filepath.Walk(expected, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(expected, path)
		if err != nil {
			return err
		}
		actualPath := filepath.Join(actual, rel)
		actualInfo, err := os.Stat(actualPath)
		if err != nil {
			return fmt.Errorf("missing file or directory in actual: %s (%w)", rel, err)
		}
		if info.IsDir() {
			if !actualInfo.IsDir() {
				return fmt.Errorf("expected directory, got file: %s", rel)
			}
			return nil
		}
		if actualInfo.IsDir() {
			return fmt.Errorf("expected file, got directory: %s", rel)
		}
		expectedContent, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		actualContent, err := os.ReadFile(actualPath)
		if err != nil {
			return err
		}
		if !bytes.Equal(bytes.TrimSpace(expectedContent), bytes.TrimSpace(actualContent)) {
			return fmt.Errorf("content mismatch in file: %s", rel)
		}
		return nil
	})
	if err != nil {
		return err
	}
	// Check for extra generated files
	return filepath.Walk(actual, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(actual, path)
		if err != nil {
			return err
		}
		expectedPath := filepath.Join(expected, rel)
		if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
			return fmt.Errorf("extra file or directory in actual: %s", rel)
		}
		return nil
	})
}
