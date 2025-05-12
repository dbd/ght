package tab

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dbd/ght/components"
)

type Model struct {
	Name     string
	Context  *components.Context
	IsActive bool
	Page     components.Page
	Focused  bool
}

func NewModel(ctx *components.Context, name string) Model {
	return Model{
		Name:     name,
		Context:  ctx,
		IsActive: false,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	m.Page.Focus()
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Context.KeyMap.Help):
			m.Page.ToggleHelp()
			cmds = append(cmds, m.Page.ToggleHelp)
		case key.Matches(msg, components.DefaultKeyMap.Close):
			cmds = append(cmds, m.Page.Blur)
		case key.Matches(msg, m.Context.KeyMap.Exit):
			return m, tea.Quit
		}
	}
	m.Page, cmd = m.Page.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.Page.View()
}

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, m.Page.Init())
	return tea.Batch(cmds...)
}
