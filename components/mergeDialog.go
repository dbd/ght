package components

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dbd/ght/internal/api"
)

type MergeOption struct {
	Label  string
	Method api.MergeMethod
}

type MergeDialogModel struct {
	context          *Context
	pullRequest      api.PullRequestResponse
	focused          bool
	cursor           int
	deleteAfterMerge bool
	mergeOptions     []MergeOption
}

func NewMergeDialogModel(ctx *Context, pr api.PullRequestResponse) *MergeDialogModel {
	var options []MergeOption
	if pr.Repository.MergeCommitAllowed {
		options = append(options, MergeOption{Label: "Create a merge commit", Method: api.MergeMethodMerge})
	}
	if pr.Repository.SquashMergeAllowed {
		options = append(options, MergeOption{Label: "Squash and merge", Method: api.MergeMethodSquash})
	}
	if pr.Repository.RebaseMergeAllowed {
		options = append(options, MergeOption{Label: "Rebase and merge", Method: api.MergeMethodRebase})
	}

	m := MergeDialogModel{
		context:          ctx,
		pullRequest:      pr,
		mergeOptions:     options,
		cursor:           0,
		deleteAfterMerge: false,
	}
	return &m
}

func (m MergeDialogModel) Init() tea.Cmd {
	return nil
}

func (m MergeDialogModel) View() string {
	doc := strings.Builder{}

	title := BoldStyle.Render("Merge Pull Request")
	doc.WriteString(title + "\n")
	doc.WriteString(m.pullRequest.Repository.NameWithOwner + "#" + strconv.FormatInt(m.pullRequest.Number, 10) + "\n")
	doc.WriteString(strings.Repeat("─", 40) + "\n\n")

	doc.WriteString("Merge method:\n")
	for i, opt := range m.mergeOptions {
		cursor := "  "
		style := lipgloss.NewStyle()
		if i == m.cursor {
			cursor = "> "
			style = style.Foreground(Green).Bold(true)
		}
		doc.WriteString(style.Render(cursor+opt.Label) + "\n")
	}

	doc.WriteString("\n")

	// Delete branch option
	deleteIdx := len(m.mergeOptions)
	cursor := "  "
	style := lipgloss.NewStyle()
	if m.cursor == deleteIdx {
		cursor = "> "
		style = style.Foreground(Green).Bold(true)
	}
	checkbox := "[ ]"
	if m.deleteAfterMerge {
		checkbox = "[x]"
	}
	doc.WriteString(style.Render(cursor+checkbox+" Delete branch after merge") + "\n")

	doc.WriteString("\n" + strings.Repeat("─", 40) + "\n")

	// Confirm option
	confirmIdx := len(m.mergeOptions) + 1
	confirmStyle := lipgloss.NewStyle()
	if m.cursor == confirmIdx {
		confirmStyle = confirmStyle.Background(Green).Foreground(Black).Bold(true).Padding(0, 2)
	} else {
		confirmStyle = confirmStyle.Foreground(Grey).Padding(0, 2)
	}
	doc.WriteString(confirmStyle.Render("Confirm Merge") + "\n")

	doc.WriteString("\n")
	doc.WriteString(DiffLineNumberStyle.Render("↑/↓ navigate • Enter select • Esc cancel"))

	return doc.String()
}

func (m MergeDialogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	totalOptions := len(m.mergeOptions) + 2 // merge options + delete checkbox + confirm button

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.context.KeyMap.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, m.context.KeyMap.Down):
			if m.cursor < totalOptions-1 {
				m.cursor++
			}
		case key.Matches(msg, m.context.KeyMap.Enter):
			deleteIdx := len(m.mergeOptions)
			confirmIdx := len(m.mergeOptions) + 1

			if m.cursor == deleteIdx {
				// Toggle delete branch option
				m.deleteAfterMerge = !m.deleteAfterMerge
			} else if m.cursor == confirmIdx {
				// Confirm merge - use first merge option if cursor was on confirm
				// Find the selected merge method (default to first available)
				method := m.mergeOptions[0].Method
				cmds = append(cmds, api.MergePullRequestCmd(m.pullRequest, method, m.deleteAfterMerge))
				m.Blur()
			} else if m.cursor < len(m.mergeOptions) {
				// Selected a merge method - execute merge with this method
				method := m.mergeOptions[m.cursor].Method
				cmds = append(cmds, api.MergePullRequestCmd(m.pullRequest, method, m.deleteAfterMerge))
				m.Blur()
			}
		}
	}
	return &m, tea.Batch(cmds...)
}

func (m *MergeDialogModel) Focus() {
	m.focused = true
	m.cursor = 0
}

func (m *MergeDialogModel) Blur() {
	m.focused = false
}

func (m *MergeDialogModel) Focused() bool {
	return m.focused
}
