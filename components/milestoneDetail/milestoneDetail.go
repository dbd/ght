package milestoneDetail

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dbd/ght/components"
	"github.com/dbd/ght/components/issueSearch"
	"github.com/dbd/ght/internal/api"
	"golang.org/x/term"
)

type Model struct {
	context   *components.Context
	repo      string
	number    int64
	milestone api.MilestoneResponse
	table     table.Model
	viewport  viewport.Model
	loaded    bool
	showHelp  bool
	focused   bool
}

var (
	fullHelp = [][]key.Binding{
		{components.DefaultKeyMap.Up, components.DefaultKeyMap.Down, components.DefaultKeyMap.Enter},
		{components.DefaultKeyMap.Close, components.DefaultKeyMap.Exit},
	}
	fetchingStatus = "Fetching milestone..."
)

func NewModel(repo string, number int64, ctx *components.Context) *Model {
	t := newIssueTable(ctx)
	ts := table.Styles{Header: components.TableHeader, Selected: components.TableSelected}
	t.SetStyles(ts)
	vp := viewport.New(ctx.ViewportWidth, ctx.ViewportHeight-1)
	vp.YPosition = ctx.ViewportYPosition
	return &Model{context: ctx, repo: repo, number: number, table: t, viewport: vp}
}

func (m Model) Init() tea.Cmd {
	m.context.StatusText = fetchingStatus
	return api.GetMilestoneCmd(m.repo, m.number)
}

func getIssueColumns(maxWidth int) []table.Column {
	fixed := []table.Column{
		{Title: "Number", Width: 8},
		{Title: "Author", Width: 20},
		{Title: "State", Width: 8},
		{Title: "Labels", Width: 20},
	}
	used := 0
	for _, c := range fixed {
		used += c.Width
	}
	titleWidth := maxWidth - used
	if titleWidth < 10 {
		titleWidth = 10
	}
	columns := append(fixed, table.Column{Title: "Title", Width: titleWidth})
	return columns
}

func newIssueTable(ctx *components.Context) table.Model {
	maxWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	columns := getIssueColumns(maxWidth)
	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(false),
		table.WithHeight(10),
	)
	return t
}

func (m Model) Update(msg tea.Msg) (components.Page, tea.Cmd) {
	var cmds []tea.Cmd
	m.viewport.Width = m.context.ViewportWidth
	m.viewport.Height = m.context.ViewportHeight - 1

	switch msg := msg.(type) {
	case api.MilestoneRefresh:
		if msg.Error == nil {
			m.milestone = msg.Milestone
			m.loaded = true
			if m.context.StatusText == fetchingStatus {
				m.context.StatusText = ""
			}
			// Populate issue table
			var rows []table.Row
			for _, issue := range m.milestone.Issues.Nodes {
				var labelNames []string
				for _, l := range issue.Labels.Nodes {
					labelNames = append(labelNames, l.Name)
				}
				labels := strings.Join(labelNames, ",")
				if len(labels) > 19 {
					labels = labels[:16] + "..."
				}
				rows = append(rows, table.Row{
					strconv.FormatInt(issue.Number, 10),
					issue.Author.Login,
					issue.State,
					labels,
					issue.Title,
				})
			}
			m.table.SetRows(rows)
			m.table.Focus()
		} else {
			m.context.StatusText = "Failed to load milestone: " + msg.Error.Error()
		}
		m.viewport.SetContent(m.renderContent())
	case tea.WindowSizeMsg:
		m.viewport.Width = m.context.ViewportWidth
		m.viewport.Height = m.context.ViewportHeight - 1
		if m.table.Width() != m.context.ViewportWidth {
			maxWidth := m.context.ViewportWidth
			m.table.SetWidth(maxWidth)
			columns := getIssueColumns(maxWidth)
			m.table.SetColumns(columns)
		}
		tableH := m.context.ViewportHeight - headerLines - 2
		if tableH < 3 {
			tableH = 3
		}
		m.table.SetHeight(tableH)
		t, cmd := m.table.Update(msg)
		m.table = t
		cmds = append(cmds, cmd)
	case tea.KeyMsg:
		if m.loaded {
			switch {
			case key.Matches(msg, m.context.KeyMap.Enter):
				cmds = append(cmds, m.openSelectedIssue())
			case key.Matches(msg, m.context.KeyMap.Up):
				if m.table.Cursor() == 0 {
					cmds = append(cmds, m.Blur)
					return &m, tea.Batch(cmds...)
				}
			}
			t, cmd := m.table.Update(msg)
			m.table = t
			cmds = append(cmds, cmd)
			m.viewport.SetContent(m.renderContent())
		}
	}

	v, vCmd := m.viewport.Update(msg)
	m.viewport = v
	cmds = append(cmds, vCmd)
	return &m, tea.Batch(cmds...)
}

// headerLines is the approximate number of lines in the milestone header.
const headerLines = 6

func (m Model) renderContent() string {
	if !m.loaded {
		return fetchingStatus
	}
	doc := strings.Builder{}
	doc.WriteString(renderMilestoneHeader(m.milestone) + "\n")
	doc.WriteString(m.table.View())
	return doc.String()
}

func (m Model) View() string {
	m.viewport.SetContent(m.renderContent())
	body := m.viewport.View()
	if m.showHelp {
		width, height, _ := term.GetSize(int(os.Stdout.Fd()))
		width = width / 2
		vc := height / 2
		body = components.RenderHelpBox(m.context.Help.FullHelpView(fullHelp), body, width, vc, 0)
	}
	return body
}

func renderMilestoneHeader(ms api.MilestoneResponse) string {
	doc := strings.Builder{}
	doc.WriteString(components.PrTitleStyle.Render(ms.Title) + "\n")

	var stateColor lipgloss.Color
	if ms.State == "OPEN" {
		stateColor = components.Green
	} else {
		stateColor = components.Grey
	}
	state := lipgloss.NewStyle().Foreground(stateColor).Bold(true).Render(ms.State)

	var due string
	if string(ms.DueOn) != "" {
		due = fmt.Sprintf(" · due %s ago", ms.DueOn.ShortSince())
	}
	doc.WriteString(fmt.Sprintf("%s%s\n", state, due))

	if ms.Description != "" {
		doc.WriteString(ms.Description + "\n")
	}

	open := ms.OpenIssues.TotalCount
	closed := ms.ClosedIssues.TotalCount
	total := open + closed
	doc.WriteString(renderProgressBar(closed, total) + fmt.Sprintf(" %d/%d closed\n", closed, total))
	doc.WriteString("\n")
	return doc.String()
}

func renderProgressBar(done, total int) string {
	const width = 20
	if total == 0 {
		bar := strings.Repeat("░", width)
		return "[" + lipgloss.NewStyle().Foreground(components.Grey).Render(bar) + "]"
	}
	filled := (done * width) / total
	if filled > width {
		filled = width
	}
	green := lipgloss.NewStyle().Foreground(components.Green).Render(strings.Repeat("█", filled))
	grey := lipgloss.NewStyle().Foreground(components.Grey).Render(strings.Repeat("░", width-filled))
	return "[" + green + grey + "]"
}

func (m Model) openSelectedIssue() tea.Cmd {
	row := m.table.SelectedRow()
	if len(row) < 1 {
		return nil
	}
	n, err := strconv.ParseInt(row[0], 10, 64)
	if err != nil {
		return nil
	}
	for _, issue := range m.milestone.Issues.Nodes {
		if issue.Number == n {
			iss := issue
			return func() tea.Msg {
				return issueSearch.OpenIssue{Issue: iss}
			}
		}
	}
	return nil
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
	return false
}

func (m *Model) GetMilestone() api.MilestoneResponse {
	return m.milestone
}

func (m *Model) GetRepo() string {
	return m.repo
}

func (m *Model) GetNumber() int64 {
	return m.number
}
