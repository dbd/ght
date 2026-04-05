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
	context       *components.Context
	pullRequest   api.PullRequestResponse
	mergeDialog   components.MergeDialogModel
	reviewDialog  components.ReviewDialogModel
	inputDialog   components.InputDialogModel
	viewport      viewport.Model
	ready         bool
	showComments  bool
	paginator     paginator.Model
	diff          string
	showHelp      bool
	isInTextInput bool
}

var (
	showComments = key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "show comments"))

	openMerge = key.NewBinding(
		key.WithKeys("m"),
		key.WithHelp("m", "merge PR"))

	openComment = key.NewBinding(
		key.WithKeys("C"),
		key.WithHelp("C", "add comment"))

	openApprove = key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "approve PR"))

	openRequestChanges = key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "request changes"))

	openAddReviewer = key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "add reviewer"))

	openAddAssignee = key.NewBinding(
		key.WithKeys("A"),
		key.WithHelp("A", "add assignee"))

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
	m.isInTextInput = false
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
	fullHelp = [][]key.Binding{{components.DefaultKeyMap.Up, components.DefaultKeyMap.Down, showComments, openMerge}, {openComment, openApprove, openRequestChanges}, {openAddReviewer, openAddAssignee}, {m.viewport.KeyMap.PageDown, m.viewport.KeyMap.PageUp, m.viewport.KeyMap.HalfPageUp, m.viewport.KeyMap.HalfPageDown}}
	m.mergeDialog = *components.NewMergeDialogModel(ctx, pr)
	m.reviewDialog = *components.NewReviewDialogModel(ctx, pr)
	m.inputDialog = *components.NewInputDialogModel(ctx, pr)
	return &m
}

func (m Model) Update(msg tea.Msg) (components.Page, tea.Cmd) {
	var cmds []tea.Cmd
	m.viewport.Width = m.context.ViewportWidth
	m.viewport.Height = m.context.ViewportHeight - 1
	p, pCmd := m.paginator.Update(msg)
	m.paginator = p
	if m.mergeDialog.Focused() {
		m.isInTextInput = true
		if msg, ok := msg.(tea.KeyMsg); ok {
			if key.Matches(msg, components.DefaultKeyMap.Exit) {
				m.mergeDialog.Blur()
				m.isInTextInput = false
				return &m, nil
			}
		} else {
			m.isInTextInput = true
		}
		md, mdCmd := m.mergeDialog.Update(msg)
		m.mergeDialog = *md.(*components.MergeDialogModel)
		cmds = append(cmds, mdCmd)
		return &m, tea.Batch(cmds...)
	}
	if m.reviewDialog.Focused() {
		m.isInTextInput = true
		if msg, ok := msg.(tea.KeyMsg); ok {
			if key.Matches(msg, components.DefaultKeyMap.Exit) {
				m.reviewDialog.Blur()
				m.isInTextInput = false
				return &m, nil
			}
		} else {
			m.isInTextInput = true
		}
		rd, rdCmd := m.reviewDialog.Update(msg)
		m.reviewDialog = *rd.(*components.ReviewDialogModel)
		cmds = append(cmds, rdCmd)
		return &m, tea.Batch(cmds...)
	}
	if m.inputDialog.Focused() {
		m.isInTextInput = true
		if msg, ok := msg.(tea.KeyMsg); ok {
			if key.Matches(msg, components.DefaultKeyMap.Exit) {
				m.inputDialog.Blur()
				m.isInTextInput = false
				return &m, nil
			}
		} else {
			m.isInTextInput = true
		}
		id, idCmd := m.inputDialog.Update(msg)
		m.inputDialog = *id.(*components.InputDialogModel)
		cmds = append(cmds, idCmd)
		return &m, tea.Batch(cmds...)
	}
	m.isInTextInput = false
	switch msg := msg.(type) {
	case api.MergeResult:
		if msg.Success {
			m.context.StatusText = "Successfully merged PR #" + strconv.FormatInt(msg.PR.Number, 10)
			cmds = append(cmds, api.GetPullRequestCmd(msg.PR.Repository.NameWithOwner, msg.PR.Number))
		} else {
			m.context.StatusText = "Merge failed: " + msg.Error.Error()
		}
	case api.PullRequestRefresh:
		if msg.Error == nil {
			m.pullRequest = msg.PR
		} else {
			m.context.StatusText = "Failed to refresh PR: " + msg.Error.Error()
		}
	case api.AssigneeResult:
		if msg.Success {
			m.context.StatusText = "Successfully added assignee"
			cmds = append(cmds, api.GetPullRequestCmd(msg.PR.Repository.NameWithOwner, msg.PR.Number))
		} else {
			m.context.StatusText = fmt.Sprintf("Failed to add assignee: %v", msg.Error)
		}
	case api.ReviewerResult:
		if msg.Success {
			m.context.StatusText = "Successfully added reviewer"
			cmds = append(cmds, api.GetPullRequestCmd(msg.PR.Repository.NameWithOwner, msg.PR.Number))
		} else {
			m.context.StatusText = fmt.Sprintf("Failed to add reviewer: %v", msg.Error)
		}
	case api.ReviewResult:
		if msg.Success {
			var actionStr string
			switch msg.Action {
			case api.ReviewActionApprove:
				actionStr = "approved"
			case api.ReviewActionRequestChanges:
				actionStr = "requested changes on"
			default:
				actionStr = "reviewed"
			}
			m.context.StatusText = "Successfully " + actionStr + " PR #" + strconv.FormatInt(msg.PR.Number, 10)
			cmds = append(cmds, api.GetPullRequestCmd(msg.PR.Repository.NameWithOwner, msg.PR.Number))
		} else {
			m.context.StatusText = "Review failed: " + msg.Error.Error() + " | " + msg.StdErr.String()
		}
	case api.CommentResult:
		if msg.Success {
			m.context.StatusText = "Comment added to PR #" + strconv.FormatInt(msg.PR.Number, 10)
			cmds = append(cmds, api.GetPullRequestCmd(msg.PR.Repository.NameWithOwner, msg.PR.Number))
		} else {
			m.context.StatusText = "Comment failed: " + msg.Error.Error()
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, showComments):
			m.showComments = !m.showComments
		case key.Matches(msg, openMerge):
			m.mergeDialog.Focus()
			m.isInTextInput = true
		case key.Matches(msg, openComment):
			m.reviewDialog.FocusWithMode(components.ReviewModeComment)
			m.isInTextInput = true
		case key.Matches(msg, openApprove):
			m.reviewDialog.FocusWithMode(components.ReviewModeApprove)
			m.isInTextInput = true
		case key.Matches(msg, openRequestChanges):
			m.reviewDialog.FocusWithMode(components.ReviewModeRequestChanges)
			m.isInTextInput = true
		case key.Matches(msg, openAddReviewer):
			m.inputDialog.FocusWithType(components.InputDialogReviewer)
			m.isInTextInput = true
		case key.Matches(msg, openAddAssignee):
			m.inputDialog.FocusWithType(components.InputDialogAssignee)
			m.isInTextInput = true
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
	if m.mergeDialog.Focused() {
		width, _, _ := term.GetSize(int(os.Stdout.Fd()))
		vc := m.context.ViewportHeight / 2
		body = components.RenderOverlay(m.mergeDialog.View(), body, width/4, vc)
	}
	if m.reviewDialog.Focused() {
		width, _, _ := term.GetSize(int(os.Stdout.Fd()))
		vc := m.context.ViewportHeight / 2
		body = components.RenderOverlay(m.reviewDialog.View(), body, width/4, vc)
	}
	if m.inputDialog.Focused() {
		width, _, _ := term.GetSize(int(os.Stdout.Fd()))
		vc := m.context.ViewportHeight / 2
		body = components.RenderOverlay(m.inputDialog.View(), body, width/4, vc)
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

func (m *Model) ShowingHelp() bool {
	return m.showHelp
}

func (m *Model) IsInTextInput() bool {
	return m.isInTextInput
}

func (m *Model) GetPullRequest() api.PullRequestResponse {
	return m.pullRequest
}
