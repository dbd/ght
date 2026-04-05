package components

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type SetupDialogModel struct {
	context *Context
	focused bool
}

type ConfigCreated struct{}

func NewSetupDialogModel(ctx *Context) *SetupDialogModel {
	return &SetupDialogModel{
		context: ctx,
		focused: true,
	}
}

func (m SetupDialogModel) Init() tea.Cmd {
	return nil
}

func (m SetupDialogModel) View() string {
	doc := strings.Builder{}
	
	doc.WriteString(BoldStyle.Render("Welcome to ght!") + "\n\n")
	doc.WriteString("No configuration file found. Would you like to create a\n")
	doc.WriteString("default configuration with these searches?\n\n")
	
	searches := []struct {
		name  string
		query string
		desc  string
	}{
		{"Assigned", "is:pr assignee:@me", "PRs assigned to you"},
		{"Review Requested", "is:pr review-requested:@me", "PRs requesting your review"},
		{"Author", "is:pr author:@me", "PRs you've authored"},
	}
	
	for _, s := range searches {
		doc.WriteString("  • " + RenderColoredText(s.name, "blue") + " - " + s.desc + "\n")
	}
	
	doc.WriteString("\n")
	doc.WriteString(RenderColoredText("Press Enter to create config", "green") + " • ")
	doc.WriteString(RenderColoredText("Esc to exit", "red") + "\n")
	
	return RenderBoxWithTitle("Setup", doc.String(), 60)
}

func (m SetupDialogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return &m, m.createDefaultConfig()
		case "esc", "ctrl+c":
			return &m, tea.Quit
		}
	}
	return &m, nil
}

func (m *SetupDialogModel) createDefaultConfig() tea.Cmd {
	return func() tea.Msg {
		// Create default searches
		searches := []Search{
			{Name: "Assigned", Query: "is:pr assignee:@me"},
			{Name: "Review Requested", Query: "is:pr review-requested:@me"},
			{Name: "Author", Query: "is:pr author:@me"},
		}
		
		config := Config{
			Pr: PrConfig{
				Searches: searches,
			},
		}
		
		if err := writeConfig(config); err != nil {
			return tea.Quit
		}
		
		return ConfigCreated{}
	}
}

func (m *SetupDialogModel) Focus() {
	m.focused = true
}

func (m *SetupDialogModel) Blur() {
	m.focused = false
}

func (m *SetupDialogModel) Focused() bool {
	return m.focused
}
