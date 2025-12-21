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
	
	commands := []struct {
		name  string
		usage string
		desc  string
	}{
		{"newtab", ":newtab", "Create a new search tab"},
		{"save-tab", ":save-tab <name>", "Save current search tab to config"},
		{"refresh", ":refresh", "Refresh current tab (search or PR)"},
		{"merge", ":merge", "Open merge dialog for current PR"},
		{"add-assignee", ":add-assignee <username>", "Add assignee to current PR"},
		{"add-reviewer", ":add-reviewer <username>", "Add reviewer to current PR"},
		{"comment", ":comment <message>", "Add comment to current PR"},
		{"approve", ":approve [message]", "Approve current PR (optional message)"},
		{"request-changes", ":request-changes <message>", "Request changes on current PR"},
		{"help", ":help", "Show this help dialog"},
	}
	
	for _, cmd := range commands {
		doc.WriteString(RenderColoredText(cmd.usage, "blue") + "\n")
		doc.WriteString("  " + cmd.desc + "\n\n")
	}
	
	doc.WriteString("\n" + BoldStyle.Render("Key Bindings") + "\n\n")
	
	bindings := []struct {
		key  string
		desc string
	}{
		{"j/k or ↑/↓", "Navigate up/down"},
		{"h/l or ←/→", "Switch tabs"},
		{"Enter", "Select/Open"},
		{"/", "Search/Filter"},
		{"?", "Toggle help"},
		{"q or Ctrl+W", "Close tab"},
		{"Esc or Ctrl+C", "Exit/Cancel"},
		{":", "Enter command mode"},
		{"Ctrl+Z", "Suspend"},
	}
	
	doc.WriteString(BoldStyle.Render("Global Keys") + "\n")
	for _, kb := range bindings {
		doc.WriteString(RenderColoredText(kb.key, "green") + " - " + kb.desc + "\n")
	}
	
	doc.WriteString("\n" + BoldStyle.Render("PR Detail Keys") + "\n")
	prKeys := []struct {
		key  string
		desc string
	}{
		{"c", "Show comments"},
		{"m", "Open merge dialog"},
		{"C", "Add comment"},
		{"a", "Approve PR"},
		{"x", "Request changes"},
		{"r", "Add reviewer"},
		{"A", "Add assignee (Shift+A)"},
	}
	
	for _, kb := range prKeys {
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
