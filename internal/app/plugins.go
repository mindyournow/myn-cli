package app

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// PluginList discovers and lists available plugins.
func (a *App) PluginList(ctx context.Context) error {
	plugins := discoverPlugins()
	if a.Formatter.JSON {
		return a.Formatter.Print(plugins)
	}
	if len(plugins) == 0 {
		return a.Formatter.Println("No plugins installed.")
	}
	tbl := a.Formatter.NewTable("NAME", "PATH")
	for _, p := range plugins {
		tbl.AddRow(p.Name, p.Path)
	}
	tbl.Render()
	return nil
}

// PluginEnable enables a plugin (placeholder — plugins are auto-discovered).
func (a *App) PluginEnable(ctx context.Context, name string) error {
	return a.Formatter.Println(fmt.Sprintf("Plugin %q enabled (auto-discovered).", name))
}

// PluginRun executes a plugin.
func (a *App) PluginRun(ctx context.Context, name string, args []string) error {
	pluginBin := "mynow-" + name
	path, err := exec.LookPath(pluginBin)
	if err != nil {
		return fmt.Errorf("plugin %q not found (expected executable %q in PATH)", name, pluginBin)
	}

	cmd := exec.CommandContext(ctx, path, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// Pass API token to plugin
	apiKey, keyErr := a.KeyStore.LoadAPIKey()
	if keyErr != nil && apiKey == "" {
		// No API key available, plugin won't have MYN_API_TOKEN
		fmt.Fprintf(os.Stderr, "Note: no API key available for plugin\n")
	}
	if apiKey != "" {
		cmd.Env = append(os.Environ(), "MYN_API_TOKEN="+apiKey)
	}
	return cmd.Run()
}

type pluginInfo struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func discoverPlugins() []pluginInfo {
	var plugins []pluginInfo
	seen := map[string]bool{}

	// Check ~/.local/share/mynow/plugins/
	home, _ := os.UserHomeDir()
	pluginDir := filepath.Join(home, ".local", "share", "mynow", "plugins")
	if entries, err := os.ReadDir(pluginDir); err == nil {
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			if !strings.HasPrefix(e.Name(), "mynow-") {
				continue
			}
			name := strings.TrimPrefix(e.Name(), "mynow-")
			if name == "" {
				continue // skip executables literally named "mynow-"
			}
			if !seen[name] {
				seen[name] = true
				plugins = append(plugins, pluginInfo{Name: name, Path: filepath.Join(pluginDir, e.Name())})
			}
		}
	} else if !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Warning: failed to read plugin directory %s: %v\n", pluginDir, err)
	}

	// Check PATH
	pathDirs := filepath.SplitList(os.Getenv("PATH"))
	for _, dir := range pathDirs {
		entries, _ := os.ReadDir(dir)
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			if !strings.HasPrefix(e.Name(), "mynow-") {
				continue
			}
			name := strings.TrimPrefix(e.Name(), "mynow-")
			if name == "" {
				continue // skip executables literally named "mynow-"
			}
			if !seen[name] {
				seen[name] = true
				plugins = append(plugins, pluginInfo{Name: name, Path: filepath.Join(dir, e.Name())})
			}
		}
	}
	return plugins
}
