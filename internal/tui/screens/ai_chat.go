package screens

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/mindyournow/myn-cli/internal/api"
	"github.com/mindyournow/myn-cli/internal/app"
	"github.com/mindyournow/myn-cli/internal/tui/components"
)

type aiConvsLoadedMsg struct{ convs []api.AIConversation }
type aiConvsErrMsg struct{ err error }
type aiConvResponseMsg struct{ content string }
type aiStreamErrMsg struct{ err error }

// chatMessage holds a single message in the local chat history.
type chatMessage struct {
	role    string // "user" or "assistant"
	content string
}

// AIChatScreen shows the AI chat interface with SSE streaming.
type AIChatScreen struct {
	app    *app.App
	width  int
	height int

	conversations   []api.AIConversation
	activeConvID    string
	activeConvTitle string

	input    components.Input
	reader   components.SSEReader
	messages []chatMessage

	loadingConvs bool
	showConvList bool
	err          error
	toast        components.Toast
}

// NewAIChatScreen creates the AI Chat screen model.
func NewAIChatScreen(application *app.App) AIChatScreen {
	inp := components.NewInput("Type a message... (Enter to send)")
	inp.Focus()
	return AIChatScreen{
		app:          application,
		input:        inp,
		reader:       components.NewSSEReader(),
		loadingConvs: true,
		showConvList: true,
		toast:        components.NewToast(),
	}
}

// Init implements tea.Model.
func (s AIChatScreen) Init() tea.Cmd {
	return s.loadConversations()
}

func (s AIChatScreen) loadConversations() tea.Cmd {
	return func() tea.Msg {
		if s.app == nil {
			return aiConvsLoadedMsg{nil}
		}
		ctx := context.Background()
		convs, err := s.app.Client.ListAIConversations(ctx)
		if err != nil {
			return aiConvsErrMsg{err}
		}
		return aiConvsLoadedMsg{convs}
	}
}

func (s AIChatScreen) sendMessage(msg string) tea.Cmd {
	return func() tea.Msg {
		if s.app == nil {
			return aiConvResponseMsg{"(no app — cannot send message)"}
		}
		ctx := context.Background()
		var buf strings.Builder
		req := api.AIChatRequest{
			CurrentMessage: msg,
			ConversationID: s.activeConvID,
		}
		err := s.app.Client.AIChatStream(ctx, req, func(event api.SSEEvent) error {
			if event.Data != "[DONE]" {
				buf.WriteString(event.Data)
			}
			return nil
		})
		if err != nil {
			return aiStreamErrMsg{err}
		}
		return aiConvResponseMsg{buf.String()}
	}
}

// Update implements tea.Model.
func (s AIChatScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case aiConvsLoadedMsg:
		s.conversations = msg.convs
		s.loadingConvs = false
		if len(s.conversations) > 0 && s.activeConvID == "" {
			s.activeConvID = s.conversations[0].ID
			s.activeConvTitle = s.conversations[0].Title
		}

	case aiConvsErrMsg:
		s.err = msg.err
		s.loadingConvs = false

	case aiConvResponseMsg:
		s.messages = append(s.messages, chatMessage{"assistant", msg.content})
		s.reader.AppendChunk(msg.content)
		s.reader.SetDone()

	case aiStreamErrMsg:
		s.reader.SetError(msg.err)
		s.toast.Show(fmt.Sprintf("Error: %v", msg.err), "error")
		var toastCmd tea.Cmd
		s.toast, toastCmd = s.toast.Update(nil)
		return s, tea.Batch(s.toast.Tick(), toastCmd)

	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			s.showConvList = !s.showConvList
			return s, nil
		case "ctrl+n":
			s.activeConvID = ""
			s.activeConvTitle = "New conversation"
			s.messages = nil
			s.reader.Reset()
			return s, nil
		case "enter":
			if text := s.input.Value(); text != "" {
				s.messages = append(s.messages, chatMessage{"user", text})
				s.input.SetValue("")
				s.reader.Reset()
				return s, s.sendMessage(text)
			}
		default:
			var inputCmd tea.Cmd
			s.input, inputCmd = s.input.Update(msg)
			var toastCmd tea.Cmd
			s.toast, toastCmd = s.toast.Update(msg)
			return s, tea.Batch(inputCmd, toastCmd)
		}
	}

	// Always pass through to toast so it can handle its internal tick.
	var toastCmd tea.Cmd
	s.toast, toastCmd = s.toast.Update(msg)
	return s, toastCmd
}

// View implements tea.Model.
func (s AIChatScreen) View() string {
	title := titleStyle.Render("AI CHAT — KAIA")

	if s.loadingConvs {
		return lipgloss.JoinVertical(lipgloss.Left, title, dimStyle.Render("  Loading conversations..."))
	}

	var rows []string
	rows = append(rows, title)

	// Conversation header
	var convHeader string
	if s.activeConvTitle != "" {
		convHeader = dimStyle.Render("  Conversation: ") + itemStyle.Render(s.activeConvTitle)
	} else {
		convHeader = dimStyle.Render("  New conversation")
	}
	rows = append(rows, convHeader, "")

	// Message history
	contentWidth := s.width - 10
	if contentWidth < 20 {
		contentWidth = 20
	}
	for _, m := range s.messages {
		if m.role == "user" {
			rows = append(rows, selectedItemStyle.Render("  You: "+truncate(m.content, contentWidth)))
		} else {
			rows = append(rows, itemStyle.Render("  Kaia: "+truncate(m.content, contentWidth)))
		}
	}

	// SSE reader (streaming content)
	rows = append(rows, "", dimStyle.Render("  Kaia: ")+s.reader.View())

	rows = append(rows, "")

	// Input
	rows = append(rows, s.input.View())
	rows = append(rows, "")
	rows = append(rows, dimStyle.Render("  Enter=send  ctrl+n=new conv  Esc=toggle list  h=history"))

	if s.toast.Visible() {
		rows = append(rows, s.toast.View())
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}
