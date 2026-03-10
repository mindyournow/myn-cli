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

func TestApp_InboxAdd(t *testing.T) {
	app, err := New()
	if err != nil {
		t.Skipf("Skipping: New() error = %v", err)
	}

	var buf bytes.Buffer
	app.SetFormatter(output.NewFormatterWithWriter(&buf, false, false, true))

	ctx := context.Background()
	err = app.InboxAdd(ctx, "Test task")
	if err != nil {
		t.Fatalf("InboxAdd() error = %v", err)
	}

	if !strings.Contains(buf.String(), "Test task") {
		t.Errorf("Output = %q, should contain 'Test task'", buf.String())
	}
}

func TestApp_InboxList(t *testing.T) {
	app, err := New()
	if err != nil {
		t.Skipf("Skipping: New() error = %v", err)
	}

	var buf bytes.Buffer
	app.SetFormatter(output.NewFormatterWithWriter(&buf, false, false, true))

	ctx := context.Background()
	err = app.InboxList(ctx)
	if err != nil {
		t.Fatalf("InboxList() error = %v", err)
	}

	if buf.String() == "" {
		t.Error("Output should not be empty")
	}
}

func TestApp_NowList(t *testing.T) {
	app, err := New()
	if err != nil {
		t.Skipf("Skipping: New() error = %v", err)
	}

	var buf bytes.Buffer
	app.SetFormatter(output.NewFormatterWithWriter(&buf, false, false, true))

	ctx := context.Background()
	err = app.NowList(ctx)
	if err != nil {
		t.Fatalf("NowList() error = %v", err)
	}

	if buf.String() == "" {
		t.Error("Output should not be empty")
	}
}

func TestApp_NowFocus(t *testing.T) {
	app, err := New()
	if err != nil {
		t.Skipf("Skipping: New() error = %v", err)
	}

	var buf bytes.Buffer
	app.SetFormatter(output.NewFormatterWithWriter(&buf, false, false, true))

	ctx := context.Background()
	err = app.NowFocus(ctx)
	if err != nil {
		t.Fatalf("NowFocus() error = %v", err)
	}

	if buf.String() == "" {
		t.Error("Output should not be empty")
	}
}

func TestApp_TaskDone(t *testing.T) {
	app, err := New()
	if err != nil {
		t.Skipf("Skipping: New() error = %v", err)
	}

	var buf bytes.Buffer
	app.SetFormatter(output.NewFormatterWithWriter(&buf, false, false, true))

	ctx := context.Background()
	err = app.TaskDone(ctx, "task-123")
	if err != nil {
		t.Fatalf("TaskDone() error = %v", err)
	}

	if !strings.Contains(buf.String(), "task-123") {
		t.Errorf("Output = %q, should contain 'task-123'", buf.String())
	}
}

func TestApp_TaskSnooze(t *testing.T) {
	app, err := New()
	if err != nil {
		t.Skipf("Skipping: New() error = %v", err)
	}

	var buf bytes.Buffer
	app.SetFormatter(output.NewFormatterWithWriter(&buf, false, false, true))

	ctx := context.Background()
	err = app.TaskSnooze(ctx, "task-456")
	if err != nil {
		t.Fatalf("TaskSnooze() error = %v", err)
	}

	if !strings.Contains(buf.String(), "task-456") {
		t.Errorf("Output = %q, should contain 'task-456'", buf.String())
	}
}

func TestApp_ReviewDaily(t *testing.T) {
	app, err := New()
	if err != nil {
		t.Skipf("Skipping: New() error = %v", err)
	}

	var buf bytes.Buffer
	app.SetFormatter(output.NewFormatterWithWriter(&buf, false, false, true))

	ctx := context.Background()
	err = app.ReviewDaily(ctx)
	if err != nil {
		t.Fatalf("ReviewDaily() error = %v", err)
	}

	if buf.String() == "" {
		t.Error("Output should not be empty")
	}
}

func TestApp_RunTUI(t *testing.T) {
	app, err := New()
	if err != nil {
		t.Skipf("Skipping: New() error = %v", err)
	}

	var buf bytes.Buffer
	app.SetFormatter(output.NewFormatterWithWriter(&buf, false, false, true))

	ctx := context.Background()
	err = app.RunTUI(ctx)
	if err != nil {
		t.Fatalf("RunTUI() error = %v", err)
	}

	if buf.String() == "" {
		t.Error("Output should not be empty")
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

func TestApp_PluginEnable(t *testing.T) {
	app, err := New()
	if err != nil {
		t.Skipf("Skipping: New() error = %v", err)
	}

	var buf bytes.Buffer
	app.SetFormatter(output.NewFormatterWithWriter(&buf, false, false, true))

	ctx := context.Background()
	err = app.PluginEnable(ctx, "my-plugin")
	if err != nil {
		t.Fatalf("PluginEnable() error = %v", err)
	}

	if !strings.Contains(buf.String(), "my-plugin") {
		t.Errorf("Output = %q, should contain 'my-plugin'", buf.String())
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
	err = app.Login(ctx, true) // device flow
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}

	if buf.String() == "" {
		t.Error("Output should not be empty")
	}
}
