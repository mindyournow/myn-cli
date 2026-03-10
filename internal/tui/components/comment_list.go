package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	commentAuthorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7c3aed")).
				Background(lipgloss.Color("#0f172a")).
				Bold(true)

	commentTimestampStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#475569")).
				Background(lipgloss.Color("#0f172a"))

	commentBodyStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#e2e8f0")).
				Background(lipgloss.Color("#0f172a"))

	commentDividerStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#475569")).
				Background(lipgloss.Color("#0f172a"))

	commentBgStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#0f172a"))
)

// TaskComment holds data for a single comment.
type TaskComment struct {
	ID         string
	AuthorName string
	Body       string
	CreatedAt  string
}

// CommentList renders a list of task comments.
type CommentList struct {
	comments []TaskComment
	width    int
	height   int
}

// NewCommentList creates a new CommentList with the given comments.
func NewCommentList(comments []TaskComment) CommentList {
	return CommentList{comments: comments}
}

// SetSize sets the available terminal dimensions.
func (cl *CommentList) SetSize(w, h int) {
	cl.width = w
	cl.height = h
}

// View renders the comment list.
func (cl CommentList) View() string {
	if len(cl.comments) == 0 {
		return commentTimestampStyle.Render("No comments yet.")
	}

	w := cl.width
	if w < 1 {
		w = 80
	}

	divider := commentDividerStyle.Render(strings.Repeat("─", w))

	var rows []string
	for i, c := range cl.comments {
		// Author + timestamp header
		header := commentAuthorStyle.Render(c.AuthorName) +
			commentTimestampStyle.Render("  "+c.CreatedAt)

		// Body text — wrap naive by width
		body := commentBodyStyle.Width(w).Render(c.Body)

		rows = append(rows, header)
		rows = append(rows, body)

		if i < len(cl.comments)-1 {
			rows = append(rows, divider)
		}
	}

	return commentBgStyle.Width(w).Render(strings.Join(rows, "\n"))
}
