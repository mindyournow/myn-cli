package integration

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// projectRoot returns the absolute path to the project root directory.
func projectRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "..")
}

func TestCLIVersion(t *testing.T) {
	cmd := exec.Command("go", "run", "./cmd/mynow", "version")
	cmd.Dir = projectRoot()
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		t.Fatalf("version command failed: %v\noutput: %s", err, out.String())
	}
	if !strings.Contains(out.String(), "mynow") {
		t.Errorf("expected 'mynow' in version output, got: %s", out.String())
	}
}

func TestCLIHelp(t *testing.T) {
	cmd := exec.Command("go", "run", "./cmd/mynow", "--help")
	cmd.Dir = projectRoot()
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	// --help exits 0 for Cobra commands
	cmd.Run() //nolint:errcheck
	if !strings.Contains(out.String(), "Usage:") {
		t.Errorf("expected 'Usage:' in help output, got: %s", out.String())
	}
}

func TestCLITaskListJSON(t *testing.T) {
	// Verify that running task list --json without auth doesn't panic;
	// it may return a non-zero exit code due to missing credentials, which is expected.
	cmd := exec.Command("go", "run", "./cmd/mynow", "task", "list", "--json")
	cmd.Dir = projectRoot()
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	// We deliberately ignore the exit code — no auth is fine here.
	cmd.Run() //nolint:errcheck
	// Just verify it doesn't produce a Go panic
	if strings.Contains(out.String(), "panic:") {
		t.Errorf("unexpected panic in CLI output: %s", out.String())
	}
}
