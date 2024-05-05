package tab

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dbd/ght/components"
)

type Model struct {
	Name     string
	Context  *components.Context
	IsActive bool
	Page     tea.Model
	Focused  bool
}

func NewModel(ctx *components.Context, name string) Model {
	return Model{
		Name:     name,
		Context:  ctx,
		IsActive: true,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Page, cmd = m.Page.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.Page.View()
}

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, m.Page.Init())
	return tea.Batch(cmds...)
}
