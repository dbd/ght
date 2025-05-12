package tab

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dbd/ght/components"
)

type Model struct {
	Name    string
	context *components.Context
	page    components.Page
}

func NewModel(ctx *components.Context, name string, page components.Page) Model {
	return Model{
		Name:    name,
		context: ctx,
		page:    page,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	m.page.Focus()
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.context.KeyMap.Help):
			m.page.ToggleHelp()
			cmds = append(cmds, m.page.ToggleHelp)
		case key.Matches(msg, components.DefaultKeyMap.Close):
			cmds = append(cmds, m.page.Blur)
		case key.Matches(msg, m.context.KeyMap.Exit):
			return m, tea.Quit
		}
	}
	m.page, cmd = m.page.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.page.View()
}

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, m.page.Init())
	return tea.Batch(cmds...)
}
