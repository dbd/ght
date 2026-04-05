package issueDetail

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dbd/ght/components"
	"github.com/dbd/ght/internal/api"
	"golang.org/x/term"
)

type Model struct {
	context       *components.Context
	issue         api.IssueResponse
	inputDialog   components.InputDialogModel
	viewport      viewport.Model
	ready         bool
	showHelp      bool
	isInTextInput bool
	// inputMode tracks what the current input dialog is for
	inputMode string // "comment" or "assignee"
}

var (
	openMilestone = key.NewBinding(
		key.WithKeys("M"),
		key.WithHelp("M", "open milestone"))

	openComment = key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "add comment"))

	openAddAssignee = key.NewBinding(
		key.WithKeys("A"),
		key.WithHelp("A", "add assignee"))

	closeIssue = key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "close/reopen issue"))

	openBrowser = key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "open in browser"))

	fullHelp = [][]key.Binding{
		{components.DefaultKeyMap.Up, components.DefaultKeyMap.Down, openMilestone},
		{openComment, openAddAssignee, closeIssue, openBrowser},
	}
)

func (m Model) Init() tea.Cmd {
	return nil
}

func NewModel(issue api.IssueResponse, ctx *components.Context) *Model {
	var m Model
	m.issue = issue
	m.context = ctx
	m.viewport = viewport.New(ctx.ViewportWidth, ctx.ViewportHeight-1)
	m.viewport.SetContent(RenderIssueDetail(issue, ctx.ViewportWidth-2))
	m.viewport.YPosition = ctx.ViewportYPosition
	m.inputDialog = *components.NewInputDialogModel(ctx, api.PullRequestResponse{})
	m.ready = true
	return &m
}

func (m Model) Update(msg tea.Msg) (components.Page, tea.Cmd) {
	var cmds []tea.Cmd
	m.viewport.Width = m.context.ViewportWidth
	m.viewport.Height = m.context.ViewportHeight - 1

	if m.inputDialog.Focused() {
		m.isInTextInput = true
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			if key.Matches(keyMsg, components.DefaultKeyMap.Exit) {
				m.inputDialog.Blur()
				m.isInTextInput = false
				return &m, nil
			}
			if key.Matches(keyMsg, components.DefaultKeyMap.Enter) {
				value := m.inputDialog.Value()
				m.inputDialog.Blur()
				m.isInTextInput = false
				if value != "" {
					switch m.inputMode {
					case "comment":
						m.context.StatusText = "Adding comment..."
						cmds = append(cmds, api.AddIssueCommentCmd(m.issue, value))
					case "assignee":
						m.context.StatusText = "Adding assignee..."
						cmds = append(cmds, api.AddIssueAssigneeCmd(m.issue, value))
					}
				}
				return &m, tea.Batch(cmds...)
			}
		}
		// Pass other keys to the textinput inside the dialog
		id, idCmd := m.inputDialog.UpdateTextOnly(msg)
		m.inputDialog = *id
		cmds = append(cmds, idCmd)
		return &m, tea.Batch(cmds...)
	}
	m.isInTextInput = false

	switch msg := msg.(type) {
	case api.IssueRefresh:
		if msg.Error == nil {
			m.issue = msg.Issue
			m.viewport.SetContent(RenderIssueDetail(m.issue, m.context.ViewportWidth-2))
		} else {
			m.context.StatusText = "Failed to refresh issue: " + msg.Error.Error()
		}
	case api.IssueCommentResult:
		if msg.Success {
			m.context.StatusText = "Comment added to issue #" + strconv.FormatInt(msg.Issue.Number, 10)
			cmds = append(cmds, api.GetIssueCmd(msg.Issue.Repository.NameWithOwner, msg.Issue.Number))
		} else {
			m.context.StatusText = "Comment failed: " + msg.Error.Error()
		}
	case api.IssueAssigneeResult:
		if msg.Success {
			m.context.StatusText = "Assignee added to issue #" + strconv.FormatInt(msg.Issue.Number, 10)
			cmds = append(cmds, api.GetIssueCmd(msg.Issue.Repository.NameWithOwner, msg.Issue.Number))
		} else {
			m.context.StatusText = fmt.Sprintf("Failed to add assignee: %v", msg.Error)
		}
	case api.IssueCloseResult:
		if msg.Success {
			cmds = append(cmds, api.GetIssueCmd(msg.Issue.Repository.NameWithOwner, msg.Issue.Number))
		} else {
			m.context.StatusText = fmt.Sprintf("Failed to update issue state: %v", msg.Error)
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, openMilestone):
			if m.issue.Milestone != nil {
				return &m, func() tea.Msg {
					return components.OpenMilestoneByNumber{
						Repo:   m.issue.Repository.NameWithOwner,
						Number: m.issue.Milestone.Number,
					}
				}
			}
		case key.Matches(msg, openComment):
			m.inputDialog.FocusWithPlaceholder("Enter comment...")
			m.inputMode = "comment"
			m.isInTextInput = true
		case key.Matches(msg, openAddAssignee):
			m.inputDialog.FocusWithPlaceholder("Enter username...")
			m.inputMode = "assignee"
			m.isInTextInput = true
		case key.Matches(msg, closeIssue):
			if m.issue.State == "OPEN" {
				m.context.StatusText = "Closing issue..."
				cmds = append(cmds, api.CloseIssueCmd(m.issue))
			} else {
				m.context.StatusText = "Reopening issue..."
				cmds = append(cmds, api.ReopenIssueCmd(m.issue))
			}
		case key.Matches(msg, openBrowser):
			cmds = append(cmds, api.OpenIssueInBrowserCmd(m.issue))
		case key.Matches(msg, components.DefaultKeyMap.Up):
			if m.viewport.AtTop() {
				cmds = append(cmds, m.Blur)
			}
		}
	}

	m.viewport.SetContent(RenderIssueDetail(m.issue, m.context.ViewportWidth-2))
	v, vCmd := m.viewport.Update(msg)
	m.viewport = v
	cmds = append(cmds, vCmd)
	return &m, tea.Batch(cmds...)
}

func (m Model) View() string {
	doc := strings.Builder{}
	doc.WriteString(m.viewport.View())
	body := doc.String()
	if m.showHelp {
		width, height, _ := term.GetSize(int(os.Stdout.Fd()))
		width = width / 2
		vc := height / 2
		body = components.RenderHelpBox(m.context.Help.FullHelpView(fullHelp), body, width, vc, 0)
	}
	if m.inputDialog.Focused() {
		width, _, _ := term.GetSize(int(os.Stdout.Fd()))
		vc := m.context.ViewportHeight / 2
		body = components.RenderOverlay(m.inputDialog.View(), body, width/4, vc)
	}
	return body
}

func (m *Model) Blur() tea.Msg {
	return components.Blur(true)
}

func (m *Model) Focus() tea.Msg {
	return components.Blur(false)
}

func (m *Model) ToggleHelp() tea.Msg {
	m.showHelp = !m.showHelp
	return m.showHelp
}

func (m *Model) IsInTextInput() bool {
	return m.isInTextInput
}

func (m *Model) GetIssue() api.IssueResponse {
	return m.issue
}
