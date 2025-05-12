package pullRequestDetail

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/cli/go-gh/v2"
	"github.com/dbd/ght/components"
	"github.com/dbd/ght/internal/api"
	"github.com/dbd/ght/utils"
	"golang.org/x/term"
)

type Model struct {
	context      *components.Context
	pullRequest  api.PullRequestResponse
	viewport     viewport.Model
	ready        bool
	showComments bool
	paginator    paginator.Model
	diff         string
	showHelp     bool
}

var (
	showComments = key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "show comments"))

	fullHelp = [][]key.Binding{{}, {}}
)

func (m Model) Init() tea.Cmd {
	return nil
}

func NewModel(pr api.PullRequestResponse, ctx *components.Context) *Model {
	var m Model
	m.pullRequest = pr
	m.context = ctx
	m.showComments = false
	m.viewport = viewport.New(m.context.ViewportWidth, m.context.ViewportHeight-1)
	m.viewport.SetContent(RenderPullRequestDetail(m.pullRequest, ctx.ViewportWidth-2))
	m.viewport.YPosition = m.context.ViewportYPosition
	m.ready = true

	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = 1
	p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")
	p.SetTotalPages(2)

	m.paginator = p
	diff, _, err := gh.Exec("pr", "diff", "--repo", pr.Repository.NameWithOwner, strconv.FormatInt(pr.Number, 10))
	if err != nil {
		log.Fatal(err)
	}
	m.diff = diff.String()
	fullHelp = [][]key.Binding{{components.DefaultKeyMap.Up, components.DefaultKeyMap.Down, components.DefaultKeyMap.Enter, showComments}, {m.viewport.KeyMap.PageDown, m.viewport.KeyMap.PageUp, m.viewport.KeyMap.HalfPageUp, m.viewport.KeyMap.HalfPageDown}}

	return &m
}

func (m Model) Update(msg tea.Msg) (components.Page, tea.Cmd) {
	var cmds []tea.Cmd
	m.viewport.Width = m.context.ViewportWidth
	m.viewport.Height = m.context.ViewportHeight - 1
	p, pCmd := m.paginator.Update(msg)
	m.paginator = p
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, showComments):
			m.showComments = !m.showComments
		case key.Matches(msg, components.DefaultKeyMap.Up):
			if m.viewport.AtTop() {
				cmds = append(cmds, m.Blur)
			}
		}
	}
	if m.paginator.Page == 0 {
		m.viewport.SetContent(RenderPullRequestDetail(m.pullRequest, m.context.ViewportWidth-2))
	} else {
		m.viewport.SetContent(m.RenderPullDiff())
	}
	v, vCmd := m.viewport.Update(msg)
	m.viewport = v
	cmds = append(cmds, vCmd, pCmd)
	return &m, tea.Batch(cmds...)
}

func (m Model) View() string {
	doc := strings.Builder{}
	doc.WriteString(components.CenterAll.Copy().Width(m.context.ViewportWidth - 2).Render(m.paginator.View()))
	doc.WriteString("\n")
	doc.WriteString(m.viewport.View())
	body := doc.String()
	if m.showHelp {
		width, height, _ := term.GetSize(int(os.Stdout.Fd()))
		width = width / 2
		vc := height / 2

		body = components.RenderHelpBox(m.context.Help.FullHelpView(fullHelp), body, width, vc, 0)
	}
	return body
}

func (m Model) RenderPullDiff() string {
	files := utils.ParseDiffText(m.diff)
	doc := strings.Builder{}
	rtm := m.pullRequest.ReviewThreads.GetPRCommentsMap()
	for _, file := range files {
		body := strings.Builder{}
		header := strings.Builder{}
		header.WriteString(file.Path + "\n")
		for _, line := range file.Preamble {
			header.WriteString(line + "\n")
		}
		for _, hunk := range file.Hunks {
			for _, line := range hunk.Lines {
				var lf string
				var side string
				var lr int64
				if line.Left {
					lf = components.DeletionsStyle.Render(line.Raw)
					side = "LEFT"
					lr = line.LeftNumber
				} else if line.Right {
					lf = components.AdditionsStyle.Render(line.Raw)
					side = "RIGHT"
					lr = line.RightNumber
				} else {
					lf = line.Raw
					lr = line.LeftNumber
				}
				lp := components.DiffLineNumberStyle.Render(fmt.Sprintf("%d,%d", line.LeftNumber, line.RightNumber))
				body.WriteString(fmt.Sprintf("%s %s\n", lp, lf))
				if m.showComments {
					if _, ok := rtm[file.Path]; ok {
						if _, ok := rtm[file.Path][side]; ok {
							for lrange, comments := range rtm[file.Path][side] {
								if lrange.EndLine == lr {
									body.WriteString(components.RenderBoxWithTitle(comments[0].Author.Login, comments[0].Body, 80))
									body.WriteString("\n")
								}
							}
						}
					}
				}
			}
			body.WriteString("\n--------------------------------------\n")
		}
		doc.WriteString(components.RenderBoxWithTitle(header.String(), body.String(), 160))
		doc.WriteString("\n")
	}
	return doc.String()
}
func RenderPullRequestDetail(pr api.PullRequestResponse, width int) string {
	doc := strings.Builder{}
	doc.WriteString(formatHeader(pr))
	body, err := glamour.Render(pr.Body, "dark")
	if err != nil {
		body = "ERROR"
	}
	doc.WriteString(components.RenderBoxWithTitleCorner(pr.Author.Login, body, width, false, true) + "\n")
	doc.WriteString(renderTimeline(pr, width))
	return doc.String()
}

func (m *Model) Blur() tea.Msg {
	return components.Blur(true)
}

func (m *Model) Focus() tea.Msg {
	return components.Blur(false)
}

func (m *Model) ToggleHelp() tea.Msg {
	m.showHelp = !m.showHelp
	return m.showHelp
}

func (m *Model) IsInTextInput() bool {
	return false
}
