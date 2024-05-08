package pullRequestSearch

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dbd/ght/components"
	"github.com/dbd/ght/internal/api"
	"golang.org/x/term"
)

type Model struct {
	Context      *components.Context
	PullRequests []api.PullRequestResponse
	Focused      bool
	table        table.Model
	query        string
}

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	fmt.Println("In pr search init")
	cmds = append(cmds, api.GetPullRequestsCmd("is:pr assignee:@me"))
	return tea.Batch(cmds...)
}

func NewModel(prs []api.PullRequestResponse, query string, ctx *components.Context) Model {
	t := newEmptyTable(prs, ctx)
	ts := table.Styles{Header: components.TableHeader, Selected: components.TableSelected}
	t.SetStyles(ts)
	return Model{Context: ctx, table: t, query: query}
}

func getColumns(maxWidth int) []table.Column {
	columns := []table.Column{
		{Title: "Age", Width: 5},
		{Title: "Repo", Width: 10},
		{Title: "Number", Width: 8},
		{Title: "Author", Width: 20},
		{Title: "Approved", Width: 10},
	}
	for _, c := range columns {
		maxWidth -= c.Width
	}
	columns = append(columns, table.Column{Title: "Title", Width: maxWidth})
	return columns
}

func newEmptyTable(prs []api.PullRequestResponse, ctx *components.Context) table.Model {
	maxWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	columns := getColumns(maxWidth)
	ctx.StatusText = fmt.Sprintf("%+v", maxWidth)
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
			rows = append(rows, table.Row{"1", pr.Repository.Name, strconv.FormatInt(pr.Number, 10), pr.Author.Login, "false", pr.Title})
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
	if m.table.Width() != m.Context.ViewportWidth {
		maxWidth := m.Context.ViewportWidth
		m.table.SetWidth(maxWidth)
		columns := getColumns(maxWidth)
		m.table.SetColumns(columns)
	}
	t, cmd := m.table.Update(msg)
	m.table = t
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	doc := strings.Builder{}
	doc.WriteString(components.BoxBorderStyle.Width(m.Context.ViewportWidth-2).Align(lipgloss.Left).Render(m.query) + "\n")
	doc.WriteString(m.table.View())
	return doc.String()
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
