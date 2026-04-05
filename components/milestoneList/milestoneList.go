package milestoneList

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
	context    *components.Context
	repo       string
	milestones []api.MilestoneListResponse
	table      table.Model
	showHelp   bool
	focused    bool
}

var (
	fullHelp = [][]key.Binding{
		{components.DefaultKeyMap.Up, components.DefaultKeyMap.Down, components.DefaultKeyMap.Enter},
		{components.DefaultKeyMap.Close, components.DefaultKeyMap.Exit},
	}
	fetchingStatus = "Fetching milestones..."
)

func NewModel(repo string, ctx *components.Context) *Model {
	t := newEmptyTable(ctx)
	ts := table.Styles{Header: components.TableHeader, Selected: components.TableSelected}
	t.SetStyles(ts)
	return &Model{context: ctx, repo: repo, table: t}
}

func (m Model) Init() tea.Cmd {
	m.context.StatusText = fetchingStatus
	return api.GetMilestonesCmd(m.repo)
}

func getColumns(maxWidth int) []table.Column {
	fixed := []table.Column{
		{Title: "Number", Width: 8},
		{Title: "State", Width: 8},
		{Title: "Due Date", Width: 12},
		{Title: "Open", Width: 6},
		{Title: "Closed", Width: 6},
	}
	used := 0
	for _, c := range fixed {
		used += c.Width
	}
	titleWidth := maxWidth - used
	if titleWidth < 10 {
		titleWidth = 10
	}
	columns := []table.Column{{Title: "Title", Width: titleWidth}}
	columns = append(columns, fixed...)
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
	case api.Milestones:
		if msg.Repo == m.repo {
			if m.context.StatusText == fetchingStatus {
				m.context.StatusText = ""
			}
			m.milestones = msg.Milestones
			var rows []table.Row
			for _, ms := range m.milestones {
				rows = append(rows, table.Row{
					ms.Title,
					strconv.FormatInt(ms.Number, 10),
					ms.State,
					formatDueDate(ms.DueOn),
					strconv.Itoa(ms.OpenIssues.TotalCount),
					strconv.Itoa(ms.ClosedIssues.TotalCount),
				})
			}
			m.table.SetRows(rows)
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.context.KeyMap.Enter):
			cmds = append(cmds, m.openMilestone(m.table.SelectedRow()))
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

func (m Model) View() string {
	doc := strings.Builder{}
	doc.WriteString(components.BoxBorderStyle.Width(m.context.ViewportWidth-2).Align(lipgloss.Left).Render("Milestones: "+m.repo) + "\n")
	doc.WriteString(m.table.View())
	body := doc.String()
	if m.showHelp {
		width, height, _ := term.GetSize(int(os.Stdout.Fd()))
		width = width / 2
		vc := height / 2
		body = components.RenderHelpBox(m.context.Help.FullHelpView(fullHelp), body, width, vc, 0)
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

func (m *Model) IsInTextInput() bool {
	return false
}

func (m *Model) GetRepo() string {
	return m.repo
}

func (m Model) openMilestone(row table.Row) tea.Cmd {
	if len(row) < 2 {
		return nil
	}
	n, err := strconv.ParseInt(row[1], 10, 64)
	if err != nil {
		return nil
	}
	repo := m.repo
	return func() tea.Msg {
		return components.OpenMilestoneByNumber{Repo: repo, Number: n}
	}
}

func formatDueDate(t api.Timestamp) string {
	if string(t) == "" {
		return "-"
	}
	return fmt.Sprintf("%s ago", t.ShortSince())
}
