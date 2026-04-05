package issueSearch

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
	context  *components.Context
	issues   []api.IssueResponse
	focused  bool
	showHelp bool
	table    table.Model
	query    string
	filter   textinput.Model
	allRows  []table.Row
}

type OpenIssue struct {
	Issue api.IssueResponse
}

var (
	fullHelp = [][]key.Binding{
		{components.DefaultKeyMap.Up, components.DefaultKeyMap.Down, components.DefaultKeyMap.Enter},
		{components.DefaultKeyMap.Filter, components.DefaultKeyMap.Close, components.DefaultKeyMap.Exit},
	}
	fetchingStatus = "Fetching issues..."
)

func (m Model) Init() tea.Cmd {
	m.table.SetHeight(m.context.ViewportHeight)
	m.context.StatusText = fetchingStatus
	return api.GetIssuesCmd(m.query)
}

func NewModel(issues []api.IssueResponse, query string, ctx *components.Context) *Model {
	t := newEmptyTable(ctx)
	ts := table.Styles{Header: components.TableHeader, Selected: components.TableSelected}
	t.SetStyles(ts)
	m := Model{context: ctx, table: t, query: query, filter: textinput.New()}
	return &m
}

func getColumns(maxWidth int) []table.Column {
	columns := []table.Column{
		{Title: "Age", Width: 10},
		{Title: "Repo", Width: 10},
		{Title: "Number", Width: 8},
		{Title: "Author", Width: 20},
		{Title: "State", Width: 8},
		{Title: "Labels", Width: 20},
	}
	for _, c := range columns {
		maxWidth -= c.Width
	}
	columns = append(columns, table.Column{Title: "Title", Width: maxWidth})
	return columns
}

func newEmptyTable(ctx *components.Context) table.Model {
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
	case api.Issues:
		if msg.Query == m.query {
			if msg.Error != nil {
				m.context.StatusText = "Error fetching issues: " + msg.Error.Error()
				break
			}
			if m.context.StatusText == fetchingStatus {
				m.context.StatusText = ""
			}
			m.issues = msg.Issues
			var rows []table.Row
			for _, issue := range m.issues {
				rows = append(rows, table.Row{
					issue.CreatedAt.ShortSince(),
					issue.Repository.Name,
					strconv.FormatInt(issue.Number, 10),
					issue.Author.Login,
					issue.State,
					formatLabels(issue),
					issue.Title,
				})
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
			cmds = append(cmds, m.openIssue(m.table.SelectedRow()))
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

// stateColStart is the character offset of the State column.
// Age(10)+Repo(10)+Number(8)+Author(20) = 48
const stateColStart = 48
const stateColWidth = 8

func colorizeTableState(tableView string) string {
	lines := strings.Split(tableView, "\n")
	for i, line := range lines {
		if strings.ContainsRune(line, '\x1b') || len(line) < stateColStart+stateColWidth {
			continue
		}
		cell := line[stateColStart : stateColStart+stateColWidth]
		value := strings.TrimRight(cell, " ")
		var color lipgloss.Color
		switch value {
		case "OPEN":
			color = components.Green
		case "CLOSED":
			color = components.Grey
		default:
			continue
		}
		colored := lipgloss.NewStyle().Foreground(color).Render(value)
		padding := strings.Repeat(" ", stateColWidth-len(value))
		lines[i] = line[:stateColStart] + colored + padding + line[stateColStart+stateColWidth:]
	}
	return strings.Join(lines, "\n")
}

func (m Model) View() string {
	doc := strings.Builder{}
	doc.WriteString(components.BoxBorderStyle.Width(m.context.ViewportWidth-2).Align(lipgloss.Left).Render(m.query) + "\n")
	doc.WriteString(colorizeTableState(m.table.View()))
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

func (m Model) openIssue(row table.Row) tea.Cmd {
	if len(row) < 3 {
		return nil
	}
	i, err := strconv.Atoi(row[2])
	if err != nil {
		return nil
	}
	var issue api.IssueResponse
	for _, iss := range m.issues {
		if row[1] == iss.Repository.Name && int64(i) == iss.Number {
			issue = iss
			break
		}
	}
	return func() tea.Msg {
		return OpenIssue{Issue: issue}
	}
}

func (m *Model) GetQuery() string {
	return m.query
}

func (m Model) filteredRows() []table.Row {
	if m.filter.Value() == "" {
		return m.allRows
	}
	val := strings.ToLower(m.filter.Value())
	var rows []table.Row
	for _, row := range m.allRows {
		if strings.Contains(strings.ToLower(fmt.Sprintf("%v", row)), val) {
			rows = append(rows, row)
		}
	}
	return rows
}

func formatLabels(issue api.IssueResponse) string {
	var names []string
	for _, l := range issue.Labels.Nodes {
		names = append(names, l.Name)
	}
	s := strings.Join(names, ",")
	if len(s) > 19 {
		s = s[:16] + "..."
	}
	return fmt.Sprintf("%-19s", s)
}
