package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dbd/ght/internal/api"
)

type InputDialogType int

const (
	InputDialogAssignee InputDialogType = iota
	InputDialogReviewer
)

type InputDialogModel struct {
	context     *Context
	pullRequest api.PullRequestResponse
	focused     bool
	dialogType  InputDialogType
	textinput   textinput.Model
}

func NewInputDialogModel(ctx *Context, pr api.PullRequestResponse) *InputDialogModel {
	ti := textinput.New()
	ti.Placeholder = "Enter username..."
	ti.Width = 40

	m := InputDialogModel{
		context:     ctx,
		pullRequest: pr,
		textinput:   ti,
	}
	return &m
}

func (m InputDialogModel) Init() tea.Cmd {
	return nil
}

func (m InputDialogModel) View() string {
	var title string
	switch m.dialogType {
	case InputDialogAssignee:
		title = "Add Assignee"
	case InputDialogReviewer:
		title = "Add Reviewer"
	}

	body := m.textinput.View() + "\n\n"
	helpText := DiffLineNumberStyle.Render("Enter to submit • Esc to cancel")
	body += helpText

	return RenderBoxWithTitle(title, body, 50)
}

func (m InputDialogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			username := m.textinput.Value()
			if username == "" {
				m.context.StatusText = "Username cannot be empty"
				return &m, nil
			}

			switch m.dialogType {
			case InputDialogAssignee:
				cmds = append(cmds, api.AddAssigneeCmd(m.pullRequest, username))
			case InputDialogReviewer:
				cmds = append(cmds, api.AddReviewerCmd(m.pullRequest, username))
			}
			m.Blur()
			return &m, tea.Batch(cmds...)
		}
	}

	var cmd tea.Cmd
	m.textinput, cmd = m.textinput.Update(msg)
	cmds = append(cmds, cmd)

	return &m, tea.Batch(cmds...)
}

func (m *InputDialogModel) Focus() {
	m.focused = true
	m.textinput.Focus()
}

func (m *InputDialogModel) FocusWithType(dialogType InputDialogType) {
	m.dialogType = dialogType
	m.Focus()
	m.textinput.SetValue("")

	switch dialogType {
	case InputDialogAssignee:
		m.textinput.Placeholder = "Enter assignee username..."
	case InputDialogReviewer:
		m.textinput.Placeholder = "Enter reviewer username..."
	}
}

func (m *InputDialogModel) Blur() {
	m.focused = false
	m.textinput.Blur()
}

func (m *InputDialogModel) Focused() bool {
	return m.focused
}
