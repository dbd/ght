package components

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dbd/ght/internal/api"
)

type ReviewDialogMode int

const (
	ReviewModeComment ReviewDialogMode = iota
	ReviewModeApprove
	ReviewModeRequestChanges
)

type ReviewDialogModel struct {
	context     *Context
	pullRequest api.PullRequestResponse
	focused     bool
	mode        ReviewDialogMode
	textarea    textarea.Model
}

func NewReviewDialogModel(ctx *Context, pr api.PullRequestResponse) *ReviewDialogModel {
	ta := textarea.New()
	ta.Placeholder = "Write a comment..."
	ta.ShowLineNumbers = false
	ta.SetWidth(50)
	ta.SetHeight(5)

	m := ReviewDialogModel{
		context:     ctx,
		pullRequest: pr,
		textarea:    ta,
		mode:        ReviewModeComment,
	}
	return &m
}

func (m ReviewDialogModel) Init() tea.Cmd {
	return nil
}

func (m ReviewDialogModel) View() string {
	var title string
	var actionColor lipgloss.Color
	switch m.mode {
	case ReviewModeComment:
		title = "Comment"
		actionColor = Yellow
	case ReviewModeApprove:
		title = "Approve"
		actionColor = Green
	case ReviewModeRequestChanges:
		title = "Request Changes"
		actionColor = Red
	}

	title = title + " " + m.pullRequest.Repository.NameWithOwner + "#" + strconv.FormatInt(m.pullRequest.Number, 10)

	body := strings.Builder{}
	if m.mode == ReviewModeComment {
		body.WriteString("Comment:\n")
	} else {
		body.WriteString("Comment (optional):\n")
	}
	body.WriteString(m.textarea.View() + "\n\n")

	submitStyle := lipgloss.NewStyle().Background(actionColor).Foreground(Black).Bold(true).Padding(0, 2)
	var submitText string
	switch m.mode {
	case ReviewModeComment:
		submitText = "Submit Comment"
	case ReviewModeApprove:
		submitText = "Approve"
	case ReviewModeRequestChanges:
		submitText = "Request Changes"
	}
	body.WriteString(submitStyle.Render(submitText) + "\n\n")

	body.WriteString(DiffLineNumberStyle.Render("Ctrl+S submit • Esc cancel"))

	return RenderBoxWithTitle(title, body.String(), 60)
}

func (m ReviewDialogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Submit):
			body := strings.TrimSpace(m.textarea.Value())

			if m.mode == ReviewModeComment {
				if body == "" {
					m.context.StatusText = "Comment cannot be empty"
					return &m, nil
				}
				cmds = append(cmds, api.AddCommentCmd(m.pullRequest, body))
			} else {
				var action api.ReviewAction
				switch m.mode {
				case ReviewModeApprove:
					action = api.ReviewActionApprove
				case ReviewModeRequestChanges:
					if body == "" {
						m.context.StatusText = "Please provide feedback when requesting changes"
						return &m, nil
					}
					action = api.ReviewActionRequestChanges
				}
				cmds = append(cmds, api.SubmitReviewCmd(m.pullRequest, action, body))
			}
			m.Blur()
			return &m, tea.Batch(cmds...)
		}
	}

	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)

	return &m, tea.Batch(cmds...)
}

func (m *ReviewDialogModel) Focus() {
	m.focused = true
	m.textarea.Focus()
	m.textarea.SetValue("")
}

func (m *ReviewDialogModel) FocusWithMode(mode ReviewDialogMode) {
	m.mode = mode
	m.Focus()

	switch mode {
	case ReviewModeComment:
		m.textarea.Placeholder = "Write a comment..."
	case ReviewModeApprove:
		m.textarea.Placeholder = "Leave a comment (optional)..."
	case ReviewModeRequestChanges:
		m.textarea.Placeholder = "Describe the changes needed..."
	}
}

func (m *ReviewDialogModel) Blur() {
	m.focused = false
	m.textarea.Blur()
}

func (m *ReviewDialogModel) Focused() bool {
	return m.focused
}
