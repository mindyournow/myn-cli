package components

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// sendKey sends a rune-based key message to a CommandPalette and returns the result.
func sendKey(cp CommandPalette, key string) (CommandPalette, tea.Cmd) {
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
	return cp.Update(msg)
}

// --- TabBar tests ---

func TestTabBar_DefaultTabs(t *testing.T) {
	tb := NewTabBar()
	if len(tb.Tabs) != 13 {
		t.Errorf("expected 13 tabs, got %d", len(tb.Tabs))
	}
	if tb.ActiveIdx != 0 {
		t.Errorf("expected ActiveIdx=0, got %d", tb.ActiveIdx)
	}
}

func TestTabBar_SetWidth(t *testing.T) {
	tb := NewTabBar()
	tb.SetWidth(100)
	if tb.Width != 100 {
		t.Errorf("expected Width=100, got %d", tb.Width)
	}
}

func TestTabBar_View_NotEmpty(t *testing.T) {
	tb := NewTabBar()
	v := tb.View()
	if v == "" {
		t.Error("View() returned empty string")
	}
}

func TestTabBar_View_ActiveTabHighlighted(t *testing.T) {
	tb := NewTabBar()
	tb.ActiveIdx = 2 // "Tasks" is index 2
	v := tb.View()
	if !strings.Contains(v, "Tasks") {
		t.Errorf("View() with ActiveIdx=2 = %q, want to contain 'Tasks'", v)
	}
}

func TestTabBar_SetActiveIdx(t *testing.T) {
	tb := NewTabBar()

	// Verify that setting ActiveIdx is reflected on the struct.
	tb.ActiveIdx = 2
	if tb.ActiveIdx != 2 {
		t.Errorf("expected ActiveIdx=2, got %d", tb.ActiveIdx)
	}

	// View must contain all tab labels regardless of the active index.
	v := tb.View()
	for _, tab := range tb.Tabs {
		if !strings.Contains(v, tab.Label) {
			t.Errorf("View() missing tab label %q", tab.Label)
		}
	}
}

// --- StatusBar tests ---

func TestStatusBar_View_NotEmpty(t *testing.T) {
	sb := NewStatusBar()
	v := sb.View()
	if v == "" {
		t.Error("View() returned empty string")
	}
}

func TestStatusBar_View_UserName(t *testing.T) {
	sb := NewStatusBar()
	sb.UserName = "alice"
	v := sb.View()
	if !strings.Contains(v, "alice") {
		t.Errorf("View() = %q, want to contain 'alice'", v)
	}
}

func TestStatusBar_View_TimerRunning(t *testing.T) {
	sb := NewStatusBar()
	sb.TimerRunning = true
	v := sb.View()
	if !strings.Contains(v, "⏱") {
		t.Errorf("View() = %q, want to contain '⏱'", v)
	}
}

func TestStatusBar_View_NotifCount(t *testing.T) {
	sb := NewStatusBar()
	sb.NotifCount = 3
	v := sb.View()
	if !strings.Contains(v, "3") {
		t.Errorf("View() = %q, want to contain '3'", v)
	}
}

// --- CommandPalette tests ---

func TestCommandPalette_Filter_Empty(t *testing.T) {
	cp := NewCommandPalette()
	cp.Reset()
	if len(cp.filtered) != len(allCommands) {
		t.Errorf("after Reset(), expected %d commands, got %d", len(allCommands), len(cp.filtered))
	}
}

func TestCommandPalette_Filter_Matches(t *testing.T) {
	cp := NewCommandPalette()
	cp.Reset()
	for _, ch := range "inbox" {
		cp, _ = sendKey(cp, string(ch))
	}
	found := false
	for _, entry := range cp.filtered {
		if strings.Contains(entry.Name, "inbox") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("after typing 'inbox', expected a result containing 'inbox', got %v", cp.filtered)
	}
}

func TestCommandPalette_Filter_NoMatch(t *testing.T) {
	cp := NewCommandPalette()
	cp.Reset()
	for _, ch := range "zzz" {
		cp, _ = sendKey(cp, string(ch))
	}
	if len(cp.filtered) != 0 {
		t.Errorf("after typing 'zzz', expected 0 results, got %d: %v", len(cp.filtered), cp.filtered)
	}
}

func TestCommandPalette_Enter_SelectsCommand(t *testing.T) {
	cp := NewCommandPalette()
	cp.Reset()
	// Type "quit" to narrow the list down to the quit command.
	for _, ch := range "quit" {
		cp, _ = sendKey(cp, string(ch))
	}
	// Press enter to select.
	_, cmd := cp.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected a cmd after pressing enter, got nil")
	}
	msg := cmd()
	if sel, ok := msg.(CommandSelectedMsg); !ok {
		t.Errorf("expected CommandSelectedMsg, got %T", msg)
	} else if sel.Command != "quit" {
		t.Errorf("expected command 'quit', got %q", sel.Command)
	}
}

func TestCommandPalette_Esc_Dismisses(t *testing.T) {
	cp := NewCommandPalette()
	_, cmd := cp.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Fatal("expected a cmd after pressing esc, got nil")
	}
	msg := cmd()
	if _, ok := msg.(CommandPaletteDismissMsg); !ok {
		t.Errorf("expected CommandPaletteDismissMsg, got %T", msg)
	}
}

// --- PriorityBadge tests ---

func TestPriorityBadge_CRITICAL(t *testing.T) {
	got := PriorityBadge("CRITICAL")
	if !strings.Contains(got, "●") {
		t.Errorf("PriorityBadge('CRITICAL') = %q, want to contain '●'", got)
	}
}

func TestPriorityBadge_OPPORTUNITY_NOW(t *testing.T) {
	got := PriorityBadge("OPPORTUNITY_NOW")
	if !strings.Contains(got, "○") {
		t.Errorf("PriorityBadge('OPPORTUNITY_NOW') = %q, want to contain '○'", got)
	}
}

func TestPriorityBadge_OVER_THE_HORIZON(t *testing.T) {
	got := PriorityBadge("OVER_THE_HORIZON")
	if !strings.Contains(got, "◌") {
		t.Errorf("PriorityBadge('OVER_THE_HORIZON') = %q, want to contain '◌'", got)
	}
}

func TestPriorityBadge_PARKING_LOT(t *testing.T) {
	got := PriorityBadge("PARKING_LOT")
	if !strings.Contains(got, "·") {
		t.Errorf("PriorityBadge('PARKING_LOT') = %q, want to contain '·'", got)
	}
}

func TestPriorityBadge_Unknown(t *testing.T) {
	got := PriorityBadge("UNKNOWN")
	if got != " " {
		t.Errorf("PriorityBadge('UNKNOWN') = %q, want ' '", got)
	}
}
