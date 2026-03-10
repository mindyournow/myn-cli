package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Formatter handles CLI output in text or JSON format.
type Formatter struct {
	JSON    bool
	Quiet   bool
	NoColor bool
	output    io.Writer // stdout for normal output
	errOutput io.Writer // stderr for errors and warnings (BUG-2 fix)
}

// NewFormatter creates a new Formatter with the specified options.
// Normal output goes to stdout; error/warning output goes to stderr (BUG-2 fix).
func NewFormatter(json, quiet, noColor bool) *Formatter {
	return &Formatter{
		JSON:      json,
		Quiet:     quiet,
		NoColor:   noColor,
		output:    os.Stdout,
		errOutput: os.Stderr,
	}
}

// NewFormatterWithWriter creates a new Formatter with a custom writer.
// Both stdout and stderr output are directed to the same writer (useful for testing).
func NewFormatterWithWriter(w io.Writer, json, quiet, noColor bool) *Formatter {
	return &Formatter{
		JSON:      json,
		Quiet:     quiet,
		NoColor:   noColor,
		output:    w,
		errOutput: w,
	}
}

// SetWriter sets the output writer for normal output.
func (f *Formatter) SetWriter(w io.Writer) {
	f.output = w
}

// Print outputs data in the configured format.
// Returns an error if JSON encoding fails.
func (f *Formatter) Print(data any) error {
	if f.Quiet {
		return nil
	}

	if f.JSON {
		enc := json.NewEncoder(f.output)
		enc.SetIndent("", "  ")
		if err := enc.Encode(data); err != nil {
			return fmt.Errorf("failed to encode JSON: %w", err)
		}
		return nil
	}

	fmt.Fprintln(f.output, data)
	return nil
}

// Printf prints a formatted string.
func (f *Formatter) Printf(format string, args ...any) error {
	if f.Quiet {
		return nil
	}
	_, err := fmt.Fprintf(f.output, format+"\n", args...)
	return err
}

// Println prints a line of text.
func (f *Formatter) Println(args ...any) error {
	if f.Quiet {
		return nil
	}
	_, err := fmt.Fprintln(f.output, args...)
	return err
}

// PrintJSON prints data as JSON regardless of the JSON setting.
// Useful for commands that always output JSON.
func (f *Formatter) PrintJSON(data any) error {
	enc := json.NewEncoder(f.output)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}
	return nil
}

// Success prints a success message.
func (f *Formatter) Success(msg string) error {
	if f.Quiet {
		return nil
	}
	if !f.NoColor {
		msg = "\033[32m" + msg + "\033[0m" // Green
	}
	_, err := fmt.Fprintln(f.output, msg)
	return err
}

// Error prints an error message to stderr (BUG-2 fix: errors go to stderr per Spec §13.2).
func (f *Formatter) Error(msg string) error {
	if f.Quiet {
		return nil
	}
	if !f.NoColor {
		msg = "\033[31mError: " + msg + "\033[0m" // Red
	} else {
		msg = "Error: " + msg
	}
	_, err := fmt.Fprintln(f.errOutput, msg)
	return err
}

// Warning prints a warning message to stderr.
func (f *Formatter) Warning(msg string) error {
	if f.Quiet {
		return nil
	}
	if !f.NoColor {
		msg = "\033[33mWarning: " + msg + "\033[0m" // Yellow
	} else {
		msg = "Warning: " + msg
	}
	_, err := fmt.Fprintln(f.errOutput, msg)
	return err
}

// Info prints an info message.
func (f *Formatter) Info(msg string) error {
	if f.Quiet {
		return nil
	}
	_, err := fmt.Fprintln(f.output, msg)
	return err
}
