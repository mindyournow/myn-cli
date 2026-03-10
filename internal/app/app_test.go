package app

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/mindyournow/myn-cli/internal/output"
)

func TestNew(t *testing.T) {
	app, err := New()
	if err != nil {
		// This might fail in CI if HOME/XDG_CONFIG_HOME is not set
		t.Skipf("New() error = %v (skipping - may need HOME/XDG_CONFIG_HOME)", err)
	}

	if app == nil {
		t.Fatal("New() returned nil app")
	}

	if app.Config == nil {
		t.Error("Config should not be nil")
	}

	if app.Client == nil {
		t.Error("Client should not be nil")
	}

	if app.Keyring == nil {
		t.Error("Keyring should not be nil")
	}

	if app.Formatter == nil {
		t.Error("Formatter should not be nil")
	}
}

func TestApp_SetFormatter(t *testing.T) {
	app, err := New()
	if err != nil {
		t.Skipf("Skipping: New() error = %v", err)
	}

	var buf bytes.Buffer
	newFormatter := output.NewFormatterWithWriter(&buf, false, false, false)

	app.SetFormatter(newFormatter)

	if app.Formatter != newFormatter {
		t.Error("Formatter should be updated")
	}
}

func TestApp_InboxAdd_RequiresAuth(t *testing.T) {
	a, err := New()
	if err != nil {
		t.Skipf("Skipping: New() error = %v", err)
	}
	a.SetFormatter(output.NewFormatterWithWriter(&bytes.Buffer{}, false, false, true))
	// Without credentials, InboxAdd should fail with auth or network error
	err = a.InboxAdd(context.Background(), "Test task")
	if err == nil {
		t.Fatal("InboxAdd() without auth: expected error, got nil")
	}
}

func TestApp_InboxList_RequiresAuth(t *testing.T) {
	a, err := New()
	if err != nil {
		t.Skipf("Skipping: New() error = %v", err)
	}
	a.SetFormatter(output.NewFormatterWithWriter(&bytes.Buffer{}, false, false, true))
	err = a.InboxList(context.Background())
	if err == nil {
		t.Fatal("InboxList() without auth: expected error, got nil")
	}
}

func TestApp_NowList_RequiresAuth(t *testing.T) {
	a, err := New()
	if err != nil {
		t.Skipf("Skipping: New() error = %v", err)
	}
	a.SetFormatter(output.NewFormatterWithWriter(&bytes.Buffer{}, false, false, true))
	err = a.NowList(context.Background())
	if err == nil {
		t.Fatal("NowList() without auth: expected error, got nil")
	}
}

func TestApp_NowFocus_RequiresAuth(t *testing.T) {
	a, err := New()
	if err != nil {
		t.Skipf("Skipping: New() error = %v", err)
	}
	a.SetFormatter(output.NewFormatterWithWriter(&bytes.Buffer{}, false, false, true))
	err = a.NowFocus(context.Background())
	if err == nil {
		t.Fatal("NowFocus() without auth: expected error, got nil")
	}
}

func TestApp_TaskComplete_RequiresAuth(t *testing.T) {
	a, err := New()
	if err != nil {
		t.Skipf("Skipping: New() error = %v", err)
	}

	var buf bytes.Buffer
	a.SetFormatter(output.NewFormatterWithWriter(&buf, false, false, true))

	// Without credentials, TaskComplete should return an auth error
	ctx := context.Background()
	err = a.TaskComplete(ctx, "task-123")
	if err == nil {
		t.Fatal("TaskComplete() without auth: expected error, got nil")
	}
	// Should be an auth error or network error (no backend available)
}

func TestApp_TaskSnoozeTask_RequiresAuth(t *testing.T) {
	a, err := New()
	if err != nil {
		t.Skipf("Skipping: New() error = %v", err)
	}

	var buf bytes.Buffer
	a.SetFormatter(output.NewFormatterWithWriter(&buf, false, false, true))

	ctx := context.Background()
	err = a.TaskSnoozeTask(ctx, "task-456", TaskSnoozeOpt{})
	if err == nil {
		t.Fatal("TaskSnoozeTask() without auth: expected error, got nil")
	}
}

func TestApp_ReviewDaily(t *testing.T) {
	a, err := New()
	if err != nil {
		t.Skipf("Skipping: New() error = %v", err)
	}

	a.SetFormatter(output.NewFormatterWithWriter(&bytes.Buffer{}, false, false, true))

	ctx := context.Background()
	err = a.ReviewDaily(ctx)
	// Without credentials, ReviewDaily should fail with an auth error.
	if err == nil {
		t.Fatal("ReviewDaily() without auth: expected error, got nil")
	}
}

func TestApp_RunTUI_Deprecated(t *testing.T) {
	app, err := New()
	if err != nil {
		t.Skipf("Skipping: New() error = %v", err)
	}

	ctx := context.Background()
	err = app.RunTUI(ctx)
	if err == nil {
		t.Fatal("RunTUI() should return error (deprecated method)")
	}
	if !strings.Contains(err.Error(), "deprecated") {
		t.Errorf("RunTUI() error = %v, want 'deprecated' message", err)
	}
}

func TestApp_PluginList(t *testing.T) {
	app, err := New()
	if err != nil {
		t.Skipf("Skipping: New() error = %v", err)
	}

	var buf bytes.Buffer
	app.SetFormatter(output.NewFormatterWithWriter(&buf, false, false, true))

	ctx := context.Background()
	err = app.PluginList(ctx)
	if err != nil {
		t.Fatalf("PluginList() error = %v", err)
	}

	if buf.String() == "" {
		t.Error("Output should not be empty")
	}
}

func TestApp_PluginEnable_NotFound(t *testing.T) {
	app, err := New()
	if err != nil {
		t.Skipf("Skipping: New() error = %v", err)
	}

	var buf bytes.Buffer
	app.SetFormatter(output.NewFormatterWithWriter(&buf, false, false, true))

	ctx := context.Background()
	err = app.PluginEnable(ctx, "nonexistent-plugin")
	if err == nil {
		t.Fatal("PluginEnable() should error for nonexistent plugin")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("Error = %v, should contain 'not found'", err)
	}
}

func TestApp_Logout(t *testing.T) {
	app, err := New()
	if err != nil {
		t.Skipf("Skipping: New() error = %v", err)
	}

	var buf bytes.Buffer
	app.SetFormatter(output.NewFormatterWithWriter(&buf, false, false, true))

	ctx := context.Background()
	err = app.Logout(ctx)
	if err != nil {
		t.Fatalf("Logout() error = %v", err)
	}

	if buf.String() == "" {
		t.Error("Output should not be empty")
	}
}

func TestApp_Login_DeviceNotImplemented(t *testing.T) {
	app, err := New()
	if err != nil {
		t.Skipf("Skipping: New() error = %v", err)
	}

	var buf bytes.Buffer
	app.SetFormatter(output.NewFormatterWithWriter(&buf, false, false, true))

	ctx := context.Background()
	err = app.Login(ctx, true) // device flow — expected to fail with "not supported"
	if err == nil {
		t.Fatal("Login() expected error for device flow, got nil")
	}
	if !strings.Contains(err.Error(), "not yet supported") {
		t.Errorf("Login() error = %v, want 'not yet supported'", err)
	}
}
