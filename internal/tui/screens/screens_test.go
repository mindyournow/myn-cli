package screens

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/mindyournow/myn-cli/internal/api"
)

// --- TasksScreen tests ---

func TestTasksScreen_Init_NotNil(t *testing.T) {
	s := NewTasksScreen(nil)
	// Init must not panic; the returned cmd may be nil or non-nil — both are fine.
	_ = s.Init()
}

func TestTasksScreen_View_Loading(t *testing.T) {
	s := NewTasksScreen(nil)
	// loading=true at construction
	v := s.View()
	if !strings.Contains(v, "Loading") {
		t.Errorf("View() = %q, want to contain 'Loading'", v)
	}
}

func TestTasksScreen_View_NoTasks(t *testing.T) {
	s := NewTasksScreen(nil)
	updated, _ := s.Update(tasksLoadedMsg{tasks: nil})
	v := updated.(TasksScreen).View()
	if !strings.Contains(v, "No tasks") {
		t.Errorf("View() = %q, want to contain 'No tasks'", v)
	}
}

func TestTasksScreen_Update_WindowSize(t *testing.T) {
	s := NewTasksScreen(nil)
	updated, _ := s.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	ts := updated.(TasksScreen)
	if ts.width != 80 || ts.height != 24 {
		t.Errorf("expected width=80 height=24, got width=%d height=%d", ts.width, ts.height)
	}
}

func TestTasksScreen_Update_KeyNav(t *testing.T) {
	s := NewTasksScreen(nil)
	tasks := []api.UnifiedTask{
		{ID: "1", Title: "Task A"},
		{ID: "2", Title: "Task B"},
		{ID: "3", Title: "Task C"},
	}
	updated, _ := s.Update(tasksLoadedMsg{tasks: tasks})
	ts := updated.(TasksScreen)
	// cursor starts at 0; pressing j should move it to 1
	updated2, _ := ts.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	ts2 := updated2.(TasksScreen)
	if ts2.cursor != 1 {
		t.Errorf("expected cursor=1 after 'j', got %d", ts2.cursor)
	}
}

// --- HabitsScreen tests ---

func TestHabitsScreen_Init_NotNil(t *testing.T) {
	s := NewHabitsScreen(nil)
	_ = s.Init()
}

func TestHabitsScreen_View_Loading(t *testing.T) {
	s := NewHabitsScreen(nil)
	v := s.View()
	if !strings.Contains(v, "Loading") {
		t.Errorf("View() = %q, want to contain 'Loading'", v)
	}
}

func TestHabitsScreen_View_NoHabits(t *testing.T) {
	s := NewHabitsScreen(nil)
	updated, _ := s.Update(habitsLoadedMsg{tasks: nil})
	v := updated.(HabitsScreen).View()
	if !strings.Contains(v, "No habits") {
		t.Errorf("View() = %q, want to contain 'No habits'", v)
	}
}

// --- NotificationsScreen tests ---

func TestNotificationsScreen_Init(t *testing.T) {
	s := NewNotificationsScreen(nil)
	_ = s.Init()
}

func TestNotificationsScreen_View_Loading(t *testing.T) {
	s := NewNotificationsScreen(nil)
	v := s.View()
	if !strings.Contains(v, "Loading") {
		t.Errorf("View() = %q, want to contain 'Loading'", v)
	}
}

func TestNotificationsScreen_View_Empty(t *testing.T) {
	s := NewNotificationsScreen(nil)
	updated, _ := s.Update(notifsLoadedMsg{notifs: nil})
	v := updated.(NotificationsScreen).View()
	if !strings.Contains(v, "No notifications") {
		t.Errorf("View() = %q, want to contain 'No notifications'", v)
	}
}

// --- AIChatScreen tests ---

func TestAIChatScreen_Init(t *testing.T) {
	s := NewAIChatScreen(nil)
	_ = s.Init()
}

func TestAIChatScreen_View_NotEmpty(t *testing.T) {
	s := NewAIChatScreen(nil)
	v := s.View()
	if v == "" {
		t.Error("View() returned empty string")
	}
}

// --- PomodoroScreen tests ---

func TestPomodoroScreen_Init(t *testing.T) {
	s := NewPomodoroScreen(nil)
	_ = s.Init()
}

func TestPomodoroScreen_View_NotEmpty(t *testing.T) {
	s := NewPomodoroScreen(nil)
	v := s.View()
	if !strings.Contains(v, "POMODORO") {
		t.Errorf("View() = %q, want to contain 'POMODORO'", v)
	}
}

func TestPomodoroScreen_Update_Pause(t *testing.T) {
	s := NewPomodoroScreen(nil)
	// Send the settings message so state is fully initialised.
	updated, _ := s.Update(pomodoroSettingsMsg{})
	s = updated.(PomodoroScreen)
	// Send the "p" key to pause.
	updated, _ = s.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	s = updated.(PomodoroScreen)
	if !s.paused {
		t.Error("expected paused=true after 'p' key")
	}
	v := s.View()
	if !strings.Contains(v, "PAUSED") {
		t.Errorf("View() after pause = %q, want to contain 'PAUSED'", v)
	}
}

// --- CompassScreen tests ---

func TestCompassScreen_Init(t *testing.T) {
	s := NewCompassScreen(nil)
	_ = s.Init()
}

func TestCompassScreen_View_Loading(t *testing.T) {
	s := NewCompassScreen(nil)
	v := s.View()
	if !strings.Contains(v, "Loading") {
		t.Errorf("View() = %q, want to contain 'Loading'", v)
	}
}

// --- StatsScreen tests ---

func TestStatsScreen_Init(t *testing.T) {
	s := NewStatsScreen(nil)
	_ = s.Init()
}

func TestStatsScreen_View_Loading(t *testing.T) {
	s := NewStatsScreen(nil)
	v := s.View()
	if !strings.Contains(v, "Loading") {
		t.Errorf("View() = %q, want to contain 'Loading'", v)
	}
}

// --- SearchScreen tests ---

func TestSearchScreen_Init(t *testing.T) {
	s := NewSearchScreen(nil)
	_ = s.Init()
}

func TestSearchScreen_View_NotEmpty(t *testing.T) {
	s := NewSearchScreen(nil)
	v := s.View()
	if v == "" {
		t.Error("View() returned empty string")
	}
}
