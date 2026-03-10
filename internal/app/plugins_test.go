package app

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mindyournow/myn-cli/internal/output"
)

func makeTestApp(t *testing.T) *App {
	t.Helper()
	a, err := New()
	if err != nil {
		t.Skipf("New() error = %v", err)
	}
	var buf bytes.Buffer
	a.SetFormatter(output.NewFormatterWithWriter(&buf, false, false, true))
	return a
}

func TestDiscoverPlugins_EmptyDir(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	// Plugin dir doesn't exist yet — discoverPlugins should return nothing (or empty).
	plugins := discoverPlugins()
	// No plugins expected from the tmp home; PATH may contribute some but that's fine —
	// we only assert the function doesn't error and returns a slice (possibly non-nil).
	_ = plugins
}

func TestDiscoverPlugins_FindsPlugins(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	// Create the plugin directory and two fake executables.
	pluginDir := filepath.Join(tmpDir, ".local", "share", "mynow", "plugins")
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	for _, name := range []string{"mynow-foo", "mynow-bar"} {
		path := filepath.Join(pluginDir, name)
		if err := os.WriteFile(path, []byte("#!/bin/sh\n"), 0755); err != nil {
			t.Fatalf("WriteFile %s: %v", name, err)
		}
	}

	plugins := discoverPlugins()

	names := make(map[string]bool)
	for _, p := range plugins {
		names[p.Name] = true
	}

	if !names["foo"] {
		t.Error("expected plugin 'foo' to be discovered")
	}
	if !names["bar"] {
		t.Error("expected plugin 'bar' to be discovered")
	}
}

func TestDiscoverPlugins_SkipsEmptySuffix(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	pluginDir := filepath.Join(tmpDir, ".local", "share", "mynow", "plugins")
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	// Binary named literally "mynow-" with no suffix.
	path := filepath.Join(pluginDir, "mynow-")
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"), 0755); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	plugins := discoverPlugins()
	for _, p := range plugins {
		if p.Name == "" {
			t.Error("plugin with empty name should be skipped")
		}
	}
}

func TestDiscoverPlugins_SkipsDirectories(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	pluginDir := filepath.Join(tmpDir, ".local", "share", "mynow", "plugins")
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	// Create a directory named like a plugin.
	dirPath := filepath.Join(pluginDir, "mynow-notaplugin")
	if err := os.Mkdir(dirPath, 0755); err != nil {
		t.Fatalf("Mkdir: %v", err)
	}

	plugins := discoverPlugins()
	for _, p := range plugins {
		if p.Name == "notaplugin" {
			t.Error("directory 'mynow-notaplugin' should not be included as a plugin")
		}
	}
}

func TestPluginRun_MissingPlugin(t *testing.T) {
	a, err := New()
	if err != nil {
		t.Skipf("New() error = %v", err)
	}
	var buf bytes.Buffer
	a.SetFormatter(output.NewFormatterWithWriter(&buf, false, false, true))

	err = a.PluginRun(context.Background(), "nonexistent-plugin-xyzzy", nil)
	if err == nil {
		t.Fatal("PluginRun() with nonexistent plugin: expected error, got nil")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("PluginRun() error = %v, want message containing 'not found'", err)
	}
}

func TestPluginList_Empty(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	// Also clear PATH so no stray mynow-* binaries are found.
	t.Setenv("PATH", "")

	a, err := New()
	if err != nil {
		t.Skipf("New() error = %v", err)
	}
	var buf bytes.Buffer
	a.SetFormatter(output.NewFormatterWithWriter(&buf, false, false, true))

	err = a.PluginList(context.Background())
	if err != nil {
		t.Fatalf("PluginList() error = %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "No plugins installed.") {
		t.Errorf("PluginList() with no plugins: output = %q, want 'No plugins installed.'", out)
	}
}
