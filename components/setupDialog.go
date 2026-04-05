package components

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type SetupDialogModel struct {
	context   *Context
	focused   bool
	issueMode bool
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

	if m.issueMode {
		doc.WriteString(BoldStyle.Render("Issue Mode Setup") + "\n\n")
		doc.WriteString("No issue searches configured. Create default issue searches?\n\n")
		doc.WriteString("  • " + RenderColoredText("My Issues", "blue") + " - Issues assigned to you\n")
		doc.WriteString("  • " + RenderColoredText("My Open Issues", "blue") + " - Open issues you authored\n")
	} else {
		doc.WriteString(BoldStyle.Render("Welcome to ght!") + "\n\n")
		doc.WriteString("No configuration file found. Would you like to create a\n")
		doc.WriteString("default configuration with these searches?\n\n")
		searches := []struct {
			name string
			desc string
		}{
			{"Assigned", "PRs assigned to you"},
			{"Review Requested", "PRs requesting your review"},
			{"Author", "PRs you've authored"},
		}
		for _, s := range searches {
			doc.WriteString("  • " + RenderColoredText(s.name, "blue") + " - " + s.desc + "\n")
		}
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
			if m.issueMode {
				return &m, tea.Quit
			}
			return &m, tea.Quit
		}
	}
	return &m, nil
}

func (m *SetupDialogModel) createDefaultConfig() tea.Cmd {
	issueMode := m.issueMode
	return func() tea.Msg {
		if issueMode {
			searches := []Search{
				{Name: "My Issues", Query: "is:issue assignee:@me is:open"},
				{Name: "My Open Issues", Query: "is:issue author:@me is:open"},
			}
			config := GetConfig()
			config.Issue.Searches = searches
			if err := writeConfig(config); err != nil {
				return tea.Quit
			}
			return IssueConfigCreated{}
		}

		// PR setup
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

// NewIssueSetupDialogModel creates a setup dialog for configuring issue searches.
func NewIssueSetupDialogModel(ctx *Context) *SetupDialogModel {
	return &SetupDialogModel{
		context:   ctx,
		focused:   true,
		issueMode: true,
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
