package pullRequestSearch

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dbd/ght/components"
	"github.com/dbd/ght/internal/api"
)

type Model struct {
	Context      *components.Context
	PullRequests []api.PullRequestResponse
	Focused      bool
	table        table.Model
}

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	fmt.Println("In pr search init")
	cmds = append(cmds, api.GetPullRequestsCmd("assignee:dbd"))
	return tea.Batch(cmds...)
}

func NewModel(prs []api.PullRequestResponse, ctx *components.Context) Model {
	table := newEmptyTable(prs, ctx)
	return Model{Context: ctx, table: table}
}

func newEmptyTable(prs []api.PullRequestResponse, ctx *components.Context) table.Model {
	columns := []table.Column{
		{Title: "Age", Width: 5},
		{Title: "Repo", Width: 10},
		{Title: "Number", Width: 4},
		{Title: "Title", Width: 25},
		{Title: "Author", Width: 10},
		{Title: "Approved", Width: 10},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(7),
	)
	return t

}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case api.PullRequests:
		m.PullRequests = msg.PullRequests
		var rows []table.Row
		for _, pr := range m.PullRequests {
			rows = append(rows, table.Row{"1", pr.Repository.Name, strconv.FormatInt(pr.Number, 10), pr.Title, pr.Author.Login, "false"})
		}
		m.table.SetRows(rows)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, components.DefaultKeyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, components.DefaultKeyMap.Escape):
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case key.Matches(msg, components.DefaultKeyMap.Enter):
			cmds = append(cmds, m.openPR(m.table.SelectedRow()))
		case key.Matches(msg, components.DefaultKeyMap.Up):
			if m.table.Cursor() == 0 {
				cmds = append(cmds, m.Blur)
			}
		}
	}
	t, cmd := m.table.Update(msg)
	m.table = t
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.table.View()
}

func (m Model) Blur() tea.Msg {
	return components.Blur(true)
}

func (m Model) openPR(row table.Row) tea.Cmd {
	var pr api.PullRequestResponse
	i, err := strconv.Atoi(row[2])
	if err != nil {
		panic(err)
	}
	for _, pr = range m.PullRequests {
		if row[1] == pr.Repository.Name && int64(i) == pr.Number {
			break
		}

	}
	return func() tea.Msg {
		return components.OpenPR{PR: pr}
	}
}
