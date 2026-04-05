package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type HelpDialogModel struct {
	context  *Context
	focused  bool
	viewport viewport.Model
}

func NewHelpDialogModel(ctx *Context) *HelpDialogModel {
	vp := viewport.New(60, 20)
	vp.SetContent(getHelpContent())
	
	m := HelpDialogModel{
		context:  ctx,
		viewport: vp,
	}
	return &m
}

func getHelpContent() string {
	doc := strings.Builder{}
	doc.WriteString(BoldStyle.Render("Commands") + "\n\n")
	
	doc.WriteString(BoldStyle.Render("Navigation") + "\n\n")
	navCmds := []struct{ usage, desc string }{
		{"I / P", "Switch to Issue / PR mode"},
		{":issues  :prs", "Switch mode via command"},
		{"h/l or ←/→", "Switch tabs"},
		{"j/k or ↑/↓", "Navigate up/down"},
		{"Enter", "Select/Open"},
		{"/", "Filter (search tabs)"},
		{"?", "Toggle help"},
		{"q or Ctrl+W or Esc", "Back / close tab"},
		{"Ctrl+C", "Exit"},
		{":", "Enter command mode"},
		{"Ctrl+Z", "Suspend"},
	}
	for _, c := range navCmds {
		doc.WriteString(RenderColoredText(c.usage, "green") + " - " + c.desc + "\n")
	}

	doc.WriteString("\n" + BoldStyle.Render("Commands") + "\n\n")
	commands := []struct{ usage, desc string }{
		{":newtab", "New PR search tab"},
		{":new-issue-tab", "New issue search tab"},
		{":milestones <owner/repo>", "Open milestone list for a repo"},
		{":save-tab <name>", "Save current tab to config (PR search, issue search, or milestones)"},
		{":refresh", "Refresh current tab"},
		{":merge", "Merge current PR"},
		{":add-assignee <user>", "Add assignee to current PR"},
		{":add-reviewer <user>", "Add reviewer to current PR"},
		{":comment <message>", "Add comment to current PR"},
		{":approve [message]", "Approve current PR"},
		{":request-changes <msg>", "Request changes on current PR"},
		{":help", "Show this help dialog"},
		{":quit", "Exit"},
	}
	for _, cmd := range commands {
		doc.WriteString(RenderColoredText(cmd.usage, "blue") + "\n")
		doc.WriteString("  " + cmd.desc + "\n\n")
	}

	doc.WriteString(BoldStyle.Render("PR Detail Keys") + "\n\n")
	prKeys := []struct{ key, desc string }{
		{"c", "Show/hide inline comments"},
		{"m", "Open merge dialog"},
		{"C", "Add comment"},
		{"a", "Approve PR"},
		{"x", "Request changes"},
		{"r", "Add reviewer"},
		{"A", "Add assignee"},
	}
	for _, kb := range prKeys {
		doc.WriteString(RenderColoredText(kb.key, "green") + " - " + kb.desc + "\n")
	}

	doc.WriteString("\n" + BoldStyle.Render("Issue Detail Keys") + "\n\n")
	issueKeys := []struct{ key, desc string }{
		{"c", "Add comment"},
		{"A", "Add assignee"},
		{"x", "Close / Reopen issue"},
		{"o", "Open in browser"},
		{"M", "Open milestone detail"},
	}
	for _, kb := range issueKeys {
		doc.WriteString(RenderColoredText(kb.key, "green") + " - " + kb.desc + "\n")
	}
	
	return doc.String()
}

func (m HelpDialogModel) Init() tea.Cmd {
	return nil
}

func (m HelpDialogModel) View() string {
	body := m.viewport.View() + "\n\n"
	helpText := DiffLineNumberStyle.Render("↑/↓ to scroll • Esc to close")
	body += helpText
	
	return RenderBoxWithTitle("Help", body, 70)
}

func (m HelpDialogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return &m, cmd
}

func (m *HelpDialogModel) Focus() {
	m.focused = true
}

func (m *HelpDialogModel) Blur() {
	m.focused = false
}

func (m *HelpDialogModel) Focused() bool {
	return m.focused
}
