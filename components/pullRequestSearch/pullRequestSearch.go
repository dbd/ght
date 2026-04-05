package pullRequestSearch

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dbd/ght/components"
	"github.com/dbd/ght/internal/api"
	"golang.org/x/term"
)

type Model struct {
	context      *components.Context
	pullRequests []api.PullRequestResponse
	focused      bool
	showHelp     bool
	table        table.Model
	query        string
	filter       textinput.Model
	allRows      []table.Row
}

type OpenPR struct {
	PR api.PullRequestResponse
}

var (
	fullHelp = [][]key.Binding{
		{components.DefaultKeyMap.Up, components.DefaultKeyMap.Down, components.DefaultKeyMap.Enter},
		{components.DefaultKeyMap.Filter, components.DefaultKeyMap.Close, components.DefaultKeyMap.Exit},
	}
	fetchingStatus = "Fetching pull requests..."
)

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	m.table.SetHeight(m.context.ViewportHeight)
	m.context.StatusText = fetchingStatus
	cmds = append(cmds, api.GetPullRequestsCmd(m.query))
	return tea.Batch(cmds...)
}

func NewModel(prs []api.PullRequestResponse, query string, ctx *components.Context) *Model {
	t := newEmptyTable(prs, ctx)
	ts := table.Styles{Header: components.TableHeader, Selected: components.TableSelected}
	t.SetStyles(ts)
	m := Model{context: ctx, table: t, query: query, filter: textinput.New()}
	return &m
}

func getColumns(maxWidth int) []table.Column {
	columns := []table.Column{
		{Title: "Age", Width: 10},
		{Title: "Repo", Width: 20},
		{Title: "Number", Width: 8},
		{Title: "Author", Width: 20},
		{Title: "Approved", Width: 10},
		{Title: "CI", Width: 10},
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
	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(ctx.ViewportHeight),
	)
	return t

}

func (m Model) Update(msg tea.Msg) (components.Page, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.table.SetHeight(m.context.ViewportHeight)
		t, cmd := m.table.Update(msg)
		m.table = t
		cmds = append(cmds, cmd)
	case api.PullRequests:
		if msg.Query == m.query {
			if m.context.StatusText == fetchingStatus {
				m.context.StatusText = ""
			}
			m.pullRequests = msg.PullRequests
			var rows []table.Row
			for _, pr := range m.pullRequests {
				rows = append(rows, table.Row{pr.CreatedAt.ShortSince(), pr.Repository.Name, strconv.FormatInt(pr.Number, 10), pr.Author.Login, "false", formatCIState(pr), pr.Title})
			}
			m.table.SetRows(rows)
			m.allRows = rows
		}
	case tea.KeyMsg:
		if m.filter.Focused() {
			if key.Matches(msg, m.context.KeyMap.Exit) {
				m.filter.Blur()
				m.filter.SetValue("")
				m.table.SetRows(m.allRows)
				break
			}
			filter, cmd := m.filter.Update(msg)
			cmds = append(cmds, cmd)
			m.filter = filter
			if key.Matches(msg, m.context.KeyMap.Enter) {
				m.filter.Blur()
				break
			}
			m.table.SetRows(m.filteredRows())
			break
		}
		switch {
		case key.Matches(msg, m.context.KeyMap.Enter):
			cmds = append(cmds, m.openPR(m.table.SelectedRow()))
		case key.Matches(msg, m.context.KeyMap.Filter):
			m.filter.Focus()
			m.filter.Placeholder = "Filter..."
		case key.Matches(msg, m.context.KeyMap.Up):
			if m.table.Cursor() == 0 {
				cmds = append(cmds, m.Blur)
			}
		}
	}
	if m.table.Width() != m.context.ViewportWidth {
		maxWidth := m.context.ViewportWidth
		m.table.SetWidth(maxWidth)
		columns := getColumns(maxWidth)
		m.table.SetColumns(columns)
	}
	t, cmd := m.table.Update(msg)
	m.table = t
	cmds = append(cmds, cmd)
	return &m, tea.Batch(cmds...)
}

// ciColStart is the character offset of the CI column in a rendered table row.
// It is the sum of all preceding column widths: Age(10)+Repo(10)+Number(8)+Author(20)+Approved(10).
const ciColStart = 58
const ciColWidth = 10

func colorizeTableCI(tableView string) string {
	lines := strings.Split(tableView, "\n")
	for i, line := range lines {
		// Lines with ANSI codes are the header or the selected row — skip them.
		if strings.ContainsRune(line, '\x1b') || len(line) < ciColStart+ciColWidth {
			continue
		}
		cell := line[ciColStart : ciColStart+ciColWidth]
		value := strings.TrimRight(cell, " ")
		if !ciAllPassing(value) {
			continue
		}
		colored := lipgloss.NewStyle().Foreground(components.Green).Render(value)
		padding := strings.Repeat(" ", ciColWidth-len(value))
		lines[i] = line[:ciColStart] + colored + padding + line[ciColStart+ciColWidth:]
	}
	return strings.Join(lines, "\n")
}

func ciAllPassing(s string) bool {
	idx := strings.IndexByte(s, '|')
	if idx < 1 {
		return false
	}
	a, err1 := strconv.Atoi(s[:idx])
	b, err2 := strconv.Atoi(s[idx+1:])
	return err1 == nil && err2 == nil && a > 0 && a == b
}

func (m Model) View() string {
	doc := strings.Builder{}
	doc.WriteString(components.BoxBorderStyle.Width(m.context.ViewportWidth-2).Align(lipgloss.Left).Render(m.query) + "\n")
	doc.WriteString(colorizeTableCI(m.table.View()))
	body := doc.String()
	if m.showHelp {
		width, height, _ := term.GetSize(int(os.Stdout.Fd()))
		width = width / 2
		vc := height / 2

		body = components.RenderHelpBox(m.context.Help.FullHelpView(fullHelp), body, width, vc, 0)
	}
	if m.filter.Focused() {
		width, _, _ := term.GetSize(int(os.Stdout.Fd()))
		width = width / 2
		vc := m.context.ViewportHeight / 2
		body = components.RenderFilter(m.filter.View(), body, width, vc, width/2)

	}
	return body
}

func (m *Model) Blur() tea.Msg {
	m.focused = false
	m.table.Blur()
	return components.Blur(true)
}

func (m *Model) Focus() tea.Msg {
	m.focused = true
	m.table.Focus()
	return m.focused
}

func (m *Model) ToggleHelp() tea.Msg {
	m.showHelp = !m.showHelp
	return m.showHelp
}

func (m *Model) ShowingHelp() bool {
	return m.showHelp
}

func (m *Model) IsInTextInput() bool {
	return m.filter.Focused()
}

func (m Model) openPR(row table.Row) tea.Cmd {
	var pr api.PullRequestResponse
	i, err := strconv.Atoi(row[2])
	if err != nil {
		panic(err)
	}
	for _, pr = range m.pullRequests {
		if row[1] == pr.Repository.Name && int64(i) == pr.Number {
			break
		}

	}
	return func() tea.Msg {
		return OpenPR{PR: pr}
	}
}

func (m *Model) GetQuery() string {
	return m.query
}

func formatCIState(pr api.PullRequestResponse) string {
	state := pr.CIState()
	if state == "" {
		return "-"
	}
	checks := pr.CIChecks()
	if len(checks) == 0 {
		switch state {
		case "SUCCESS":
			return "pass"
		case "FAILURE":
			return "fail"
		case "PENDING":
			return "running"
		default:
			return "-"
		}
	}
	passing := 0
	for _, c := range checks {
		switch c.Type {
		case "CheckRun":
			if c.CheckRun.Conclusion == "SUCCESS" {
				passing++
			}
		case "StatusContext":
			if c.StatusContext.State == "SUCCESS" {
				passing++
			}
		}
	}
	return fmt.Sprintf("%d|%d", passing, len(checks))
}

func (m Model) filteredRows() []table.Row {
	var rows []table.Row
	if m.filter.Value() != "" {
		val := strings.ToLower(m.filter.Value())
		for _, row := range m.allRows {
			if strings.Contains(strings.ToLower(fmt.Sprintf("%v", row)), val) {
				rows = append(rows, row)
			}
		}
	} else {
		rows = m.allRows
	}
	return rows
}
