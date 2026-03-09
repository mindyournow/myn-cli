package integration

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

const (
	backendPort     = 17000
	composeProject  = "myn-cli-test"
	startupTimeout  = 120 * time.Second
	healthPollDelay = 2 * time.Second
)

// BackendURL returns the base URL for the test MYN backend.
func BackendURL() string {
	if url := os.Getenv("MYN_TEST_BACKEND_URL"); url != "" {
		return url
	}
	return fmt.Sprintf("http://localhost:%d", backendPort)
}

// composeFilePath returns the absolute path to the test docker-compose.yml.
func composeFilePath() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "docker-compose.yml")
}

// EnsureBackend starts the MYN backend stack if not already running.
// If MYN_TEST_BACKEND_URL is set, it assumes an external backend and skips Docker.
func EnsureBackend(t *testing.T) {
	t.Helper()

	if os.Getenv("MYN_INTEGRATION_TEST") == "" {
		t.Skip("skipping integration test (set MYN_INTEGRATION_TEST=1 to run)")
	}

	// If pointing at an external backend, just health-check it
	if os.Getenv("MYN_TEST_BACKEND_URL") != "" {
		waitForHealth(t, BackendURL())
		return
	}

	// Resolve MYN backend source path for Docker build
	backendPath := os.Getenv("MYN_BACKEND_PATH")
	if backendPath == "" {
		// Default: sibling directory relative to common project layout
		backendPath = filepath.Join(os.Getenv("HOME"), "Projects", "myn", "api")
	}
	if _, err := os.Stat(filepath.Join(backendPath, "Dockerfile")); os.IsNotExist(err) {
		t.Fatalf("MYN backend not found at %s (set MYN_BACKEND_PATH)", backendPath)
	}

	composePath := composeFilePath()

	// Start the stack
	cmd := exec.Command("docker", "compose",
		"-f", composePath,
		"-p", composeProject,
		"up", "-d", "--build", "--wait",
	)
	cmd.Env = append(os.Environ(), fmt.Sprintf("MYN_BACKEND_PATH=%s", backendPath))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	t.Logf("Starting MYN backend stack (source: %s)...", backendPath)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to start backend stack: %v", err)
	}

	// Register cleanup
	t.Cleanup(func() {
		teardown := exec.Command("docker", "compose",
			"-f", composePath,
			"-p", composeProject,
			"down", "-v",
		)
		teardown.Stdout = os.Stdout
		teardown.Stderr = os.Stderr
		teardown.Run()
	})

	waitForHealth(t, BackendURL())
}

func waitForHealth(t *testing.T, baseURL string) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), startupTimeout)
	defer cancel()

	healthURL := baseURL + "/actuator/health"
	t.Logf("Waiting for backend health at %s...", healthURL)

	for {
		select {
		case <-ctx.Done():
			t.Fatalf("Backend did not become healthy within %s", startupTimeout)
		default:
			resp, err := http.Get(healthURL)
			if err == nil && resp.StatusCode == 200 {
				resp.Body.Close()
				t.Log("Backend is healthy.")
				return
			}
			if resp != nil {
				resp.Body.Close()
			}
			time.Sleep(healthPollDelay)
		}
	}
}
