package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestNewFormatter(t *testing.T) {
	f := NewFormatter(true, false, true)

	if !f.JSON {
		t.Error("JSON should be true")
	}
	if f.Quiet {
		t.Error("Quiet should be false")
	}
	if !f.NoColor {
		t.Error("NoColor should be true")
	}
}

func TestFormatter_Print_JSON(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatterWithWriter(&buf, true, false, false)

	data := map[string]string{"key": "value"}
	err := f.Print(data)
	if err != nil {
		t.Fatalf("Print() error = %v", err)
	}

	// Verify it's valid JSON
	var result map[string]string
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}

	if result["key"] != "value" {
		t.Errorf("key = %q, want %q", result["key"], "value")
	}
}

func TestFormatter_Print_Text(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatterWithWriter(&buf, false, false, false)

	data := "hello world"
	err := f.Print(data)
	if err != nil {
		t.Fatalf("Print() error = %v", err)
	}

	if !strings.Contains(buf.String(), "hello world") {
		t.Errorf("Output = %q, should contain 'hello world'", buf.String())
	}
}

func TestFormatter_Print_Quiet(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatterWithWriter(&buf, false, true, false)

	data := "hello world"
	err := f.Print(data)
	if err != nil {
		t.Fatalf("Print() error = %v", err)
	}

	if buf.String() != "" {
		t.Errorf("Output in quiet mode = %q, should be empty", buf.String())
	}
}

func TestFormatter_Printf(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatterWithWriter(&buf, false, false, false)

	err := f.Printf("Hello %s, number %d\n", "world", 42)
	if err != nil {
		t.Fatalf("Printf() error = %v", err)
	}

	expected := "Hello world, number 42\n"
	if buf.String() != expected {
		t.Errorf("Output = %q, want %q", buf.String(), expected)
	}
}

func TestFormatter_Printf_Quiet(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatterWithWriter(&buf, false, true, false)

	err := f.Printf("Hello %s", "world")
	if err != nil {
		t.Fatalf("Printf() error = %v", err)
	}

	if buf.String() != "" {
		t.Errorf("Output in quiet mode = %q, should be empty", buf.String())
	}
}

func TestFormatter_Println(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatterWithWriter(&buf, false, false, false)

	err := f.Println("line1", "line2")
	if err != nil {
		t.Fatalf("Println() error = %v", err)
	}

	expected := "line1 line2\n"
	if buf.String() != expected {
		t.Errorf("Output = %q, want %q", buf.String(), expected)
	}
}

func TestFormatter_Println_Quiet(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatterWithWriter(&buf, false, true, false)

	err := f.Println("hello")
	if err != nil {
		t.Fatalf("Println() error = %v", err)
	}

	if buf.String() != "" {
		t.Errorf("Output in quiet mode = %q, should be empty", buf.String())
	}
}

func TestFormatter_PrintJSON(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatterWithWriter(&buf, false, false, false) // JSON output forced

	data := map[string]int{"number": 42}
	err := f.PrintJSON(data)
	if err != nil {
		t.Fatalf("PrintJSON() error = %v", err)
	}

	var result map[string]int
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}

	if result["number"] != 42 {
		t.Errorf("number = %d, want %d", result["number"], 42)
	}
}

func TestFormatter_PrintJSON_Quiet(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatterWithWriter(&buf, false, true, false)

	// PrintJSON should still output even in quiet mode (it's forced JSON output)
	data := map[string]string{"key": "value"}
	err := f.PrintJSON(data)
	if err != nil {
		t.Fatalf("PrintJSON() error = %v", err)
	}

	if buf.String() == "" {
		t.Error("PrintJSON should output even in quiet mode")
	}
}

func TestFormatter_Success(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatterWithWriter(&buf, false, false, false)

	err := f.Success("Operation completed")
	if err != nil {
		t.Fatalf("Success() error = %v", err)
	}

	// With color enabled, should have ANSI codes
	if !strings.Contains(buf.String(), "Operation completed") {
		t.Errorf("Output = %q, should contain 'Operation completed'", buf.String())
	}
}

func TestFormatter_Success_NoColor(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatterWithWriter(&buf, false, false, true)

	err := f.Success("Operation completed")
	if err != nil {
		t.Fatalf("Success() error = %v", err)
	}

	// Should not have ANSI codes when NoColor is true
	if strings.Contains(buf.String(), "\033[") {
		t.Error("Output should not contain ANSI codes when NoColor is true")
	}

	if !strings.Contains(buf.String(), "Operation completed") {
		t.Errorf("Output = %q, should contain 'Operation completed'", buf.String())
	}
}

func TestFormatter_Success_Quiet(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatterWithWriter(&buf, false, true, false)

	err := f.Success("Operation completed")
	if err != nil {
		t.Fatalf("Success() error = %v", err)
	}

	if buf.String() != "" {
		t.Errorf("Output in quiet mode = %q, should be empty", buf.String())
	}
}

func TestFormatter_Error(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatterWithWriter(&buf, false, false, false)

	err := f.Error("Something went wrong")
	if err != nil {
		t.Fatalf("Error() error = %v", err)
	}

	if !strings.Contains(buf.String(), "Error:") {
		t.Errorf("Output = %q, should contain 'Error:'", buf.String())
	}
	if !strings.Contains(buf.String(), "Something went wrong") {
		t.Errorf("Output = %q, should contain 'Something went wrong'", buf.String())
	}
}

func TestFormatter_Error_NoColor(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatterWithWriter(&buf, false, false, true)

	err := f.Error("Something went wrong")
	if err != nil {
		t.Fatalf("Error() error = %v", err)
	}

	// Should still have "Error:" prefix without color
	if !strings.HasPrefix(buf.String(), "Error:") {
		t.Errorf("Output = %q, should start with 'Error:'", buf.String())
	}
}

func TestFormatter_Error_Quiet(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatterWithWriter(&buf, false, true, false)

	err := f.Error("Something went wrong")
	if err != nil {
		t.Fatalf("Error() error = %v", err)
	}

	if buf.String() != "" {
		t.Errorf("Output in quiet mode = %q, should be empty", buf.String())
	}
}

func TestFormatter_Warning(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatterWithWriter(&buf, false, false, false)

	err := f.Warning("This is a warning")
	if err != nil {
		t.Fatalf("Warning() error = %v", err)
	}

	if !strings.Contains(buf.String(), "Warning:") {
		t.Errorf("Output = %q, should contain 'Warning:'", buf.String())
	}
	if !strings.Contains(buf.String(), "This is a warning") {
		t.Errorf("Output = %q, should contain 'This is a warning'", buf.String())
	}
}

func TestFormatter_Warning_NoColor(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatterWithWriter(&buf, false, false, true)

	err := f.Warning("This is a warning")
	if err != nil {
		t.Fatalf("Warning() error = %v", err)
	}

	// Should still have "Warning:" prefix without color
	if !strings.HasPrefix(buf.String(), "Warning:") {
		t.Errorf("Output = %q, should start with 'Warning:'", buf.String())
	}
}

func TestFormatter_Info(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatterWithWriter(&buf, false, false, false)

	err := f.Info("Information message")
	if err != nil {
		t.Fatalf("Info() error = %v", err)
	}

	if !strings.Contains(buf.String(), "Information message") {
		t.Errorf("Output = %q, should contain 'Information message'", buf.String())
	}
}

func TestFormatter_SetWriter(t *testing.T) {
	var buf1 bytes.Buffer
	var buf2 bytes.Buffer
	f := NewFormatterWithWriter(&buf1, false, false, false)

	f.Println("first")
	f.SetWriter(&buf2)
	f.Println("second")

	if !strings.Contains(buf1.String(), "first") {
		t.Error("buf1 should contain 'first'")
	}
	if strings.Contains(buf1.String(), "second") {
		t.Error("buf1 should not contain 'second'")
	}
	if !strings.Contains(buf2.String(), "second") {
		t.Error("buf2 should contain 'second'")
	}
}
