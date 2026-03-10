package screens

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/mindyournow/myn-cli/internal/api"
	"github.com/mindyournow/myn-cli/internal/app"
	"github.com/mindyournow/myn-cli/internal/tui/components"
)

type notifsLoadedMsg struct{ notifs []api.Notification }
type notifsErrMsg struct{ err error }
type notifActionDoneMsg struct{ msg string }
type notifActionErrMsg struct{ err error }

// NotificationsScreen shows the user's notifications.
type NotificationsScreen struct {
	app    *app.App
	width  int
	height int
	cursor int

	notifications []api.Notification
	loading       bool
	err           error
	toast         components.Toast
}

// NewNotificationsScreen creates the Notifications screen model.
func NewNotificationsScreen(application *app.App) NotificationsScreen {
	return NotificationsScreen{
		app:     application,
		loading: true,
		toast:   components.NewToast(),
	}
}

// Init implements tea.Model.
func (s NotificationsScreen) Init() tea.Cmd {
	return s.loadData()
}

func (s NotificationsScreen) loadData() tea.Cmd {
	return func() tea.Msg {
		if s.app == nil {
			return notifsLoadedMsg{nil}
		}
		ctx := context.Background()
		notifs, err := s.app.Client.ListNotifications(ctx, 0, 50)
		if err != nil {
			return notifsErrMsg{err}
		}
		return notifsLoadedMsg{notifs}
	}
}

func (s NotificationsScreen) markRead(id string) tea.Cmd {
	return func() tea.Msg {
		if s.app == nil {
			return notifActionDoneMsg{"ok"}
		}
		ctx := context.Background()
		if err := s.app.Client.MarkNotificationRead(ctx, id); err != nil {
			return notifActionErrMsg{err}
		}
		return notifActionDoneMsg{"Marked as read"}
	}
}

func (s NotificationsScreen) markAllRead() tea.Cmd {
	return func() tea.Msg {
		if s.app == nil {
			return notifActionDoneMsg{"ok"}
		}
		ctx := context.Background()
		if err := s.app.Client.MarkAllNotificationsRead(ctx); err != nil {
			return notifActionErrMsg{err}
		}
		return notifActionDoneMsg{"All marked as read"}
	}
}

func (s NotificationsScreen) deleteNotif(id string) tea.Cmd {
	return func() tea.Msg {
		if s.app == nil {
			return notifActionDoneMsg{"ok"}
		}
		ctx := context.Background()
		if err := s.app.Client.DeleteNotification(ctx, id); err != nil {
			return notifActionErrMsg{err}
		}
		return notifActionDoneMsg{"Deleted"}
	}
}

// Update implements tea.Model.
func (s NotificationsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case notifsLoadedMsg:
		s.notifications = msg.notifs
		s.loading = false
		if s.cursor >= len(s.notifications) {
			s.cursor = max(0, len(s.notifications)-1)
		}

	case notifsErrMsg:
		s.err = msg.err
		s.loading = false

	case notifActionDoneMsg:
		s.toast.Show(msg.msg, "success")
		var toastCmd tea.Cmd
		s.toast, toastCmd = s.toast.Update(nil)
		return s, tea.Batch(s.toast.Tick(), s.loadData(), toastCmd)

	case notifActionErrMsg:
		s.toast.Show(fmt.Sprintf("Error: %v", msg.err), "error")
		var toastCmd tea.Cmd
		s.toast, toastCmd = s.toast.Update(nil)
		return s, tea.Batch(s.toast.Tick(), toastCmd)

	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if s.cursor < len(s.notifications)-1 {
				s.cursor++
			}
		case "k", "up":
			if s.cursor > 0 {
				s.cursor--
			}
		case "r":
			if s.cursor < len(s.notifications) {
				return s, s.markRead(s.notifications[s.cursor].ID)
			}
		case "R":
			return s, s.markAllRead()
		case "d":
			if s.cursor < len(s.notifications) {
				return s, s.deleteNotif(s.notifications[s.cursor].ID)
			}
		case "g":
			s.loading = true
			return s, s.loadData()
		}
	}

	// Always pass through to toast so it can handle its internal tick.
	var toastCmd tea.Cmd
	s.toast, toastCmd = s.toast.Update(msg)
	return s, toastCmd
}

// View implements tea.Model.
func (s NotificationsScreen) View() string {
	title := titleStyle.Render("NOTIFICATIONS")

	if s.loading {
		return lipgloss.JoinVertical(lipgloss.Left, title, dimStyle.Render("  Loading..."))
	}
	if s.err != nil {
		return lipgloss.JoinVertical(lipgloss.Left, title, formatError(s.err))
	}
	if len(s.notifications) == 0 {
		return lipgloss.JoinVertical(lipgloss.Left,
			title,
			dimStyle.Render("  No notifications."),
			dimStyle.Render("  g=refresh"),
			s.toast.View(),
		)
	}

	var rows []string
	rows = append(rows, title)

	for i, n := range s.notifications {
		readIcon := "●" // unread
		if n.IsRead {
			readIcon = "○"
		}

		titleStr := truncate(n.Title, 50)
		bodyStr := ""
		if n.Body != "" {
			bodyStr = "  " + truncate(n.Body, 40)
		}
		row := fmt.Sprintf("%s %-50s%s", readIcon, titleStr, bodyStr)

		if i == s.cursor {
			rows = append(rows, selectedItemStyle.Render("► "+row))
		} else if n.IsRead {
			rows = append(rows, dimStyle.Render("  "+row))
		} else {
			rows = append(rows, itemStyle.Render("  "+row))
		}
	}

	rows = append(rows, "")
	rows = append(rows, dimStyle.Render("  j/k navigate  r=mark read  R=mark all read  d=delete  g=refresh"))
	rows = append(rows, s.toast.View())

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}
