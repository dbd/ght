package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dbd/ght/components"
	"github.com/dbd/ght/components/pullRequestDetail"
	"github.com/dbd/ght/components/pullRequestSearch"
	"github.com/dbd/ght/components/tab"
	"github.com/dbd/ght/internal/api"
	"golang.org/x/term"
)

type Model struct {
	Tabs      []tab.Model
	activeTab int
	viewport  viewport.Model
	ready     bool
	context   *components.Context
	focused   bool
	keymap    components.KeyMap
	help      help.Model
}

func initializeModel() Model {
	ctx := components.Context{}
	tabs := []tab.Model{}
	for _, name := range []string{"Assigned"} {
		t := tab.NewModel(&ctx, name)
		if name == "Assigned" {
			t.Page = pullRequestSearch.NewModel([]api.PullRequestResponse{}, "is:pr assignee:@me", &ctx)
		}
		tabs = append(tabs, t)
	}
	h := help.New()
	h.Styles.FullKey.UnsetForeground()
	h.Styles.FullDesc.UnsetForeground()
	h.Styles.FullKey.UnsetForeground()
	return Model{
		Tabs:    tabs,
		context: &ctx,
		focused: true,
		keymap:  components.DefaultKeyMap,
		help:    h,
	}
}

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, tab := range m.Tabs {
		cmds = append(cmds, tab.Init())
	}
	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	activeTab := m.Tabs[m.activeTab]
	switch msg := msg.(type) {
	case components.OpenPR:
		var alreadyOpened bool
		for counter, tab := range m.Tabs {
			if tab.Name == msg.PR.Title {
				alreadyOpened = true
				_ = counter
			}
		}
		if alreadyOpened {
			m.activeTab = len(m.Tabs) - 1
			m.focused = true
		} else {
			t := tab.NewModel(m.context, msg.PR.Title)
			t.Page = pullRequestDetail.NewModel(msg.PR, m.context)
			m.Tabs = append(m.Tabs, t)
			m.activeTab = len(m.Tabs) - 1
			m.focused = true
		}
	case components.Blur:
		m.focused = true
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight - 1
			m.ready = true
			m.context.ViewportWidth = m.viewport.Width
			m.context.ViewportHeight = m.viewport.Height
			m.context.ViewportYOffset = m.viewport.YOffset
			m.context.ViewportYPosition = m.viewport.YPosition
			m.viewport.SetContent(activeTab.Page.View())
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
			m.context.ViewportWidth = m.viewport.Width
			m.context.ViewportHeight = m.viewport.Height
			m.context.ViewportYOffset = m.viewport.YOffset
			m.context.ViewportYPosition = m.viewport.YPosition
			activeTab, cmd := activeTab.Update(msg)
			cmds = append(cmds, cmd)
			m.Tabs[m.activeTab] = activeTab
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Help):
			m.help.ShowAll = !m.help.ShowAll
		}
		if !m.focused {
			activeTab, cmd := activeTab.Update(msg)
			cmds = append(cmds, cmd)
			m.Tabs[m.activeTab] = activeTab
		} else {
			switch {
			case key.Matches(msg, components.DefaultKeyMap.Down):
				m.focused = false
			case key.Matches(msg, components.DefaultKeyMap.Quit):
				return m, tea.Quit
			case key.Matches(msg, components.DefaultKeyMap.Left):
				if m.activeTab > 0 {
					m.activeTab--
				}
			case key.Matches(msg, components.DefaultKeyMap.Right):
				if m.activeTab < len(m.Tabs)-1 {
					m.activeTab++
				}

			}
		}
	default:
		activeTab, cmd := activeTab.Update(msg)
		cmds = append(cmds, cmd)
		m.Tabs[m.activeTab] = activeTab
	}
	return m, tea.Batch(cmds...)
}

func (m Model) headerView() string {
	width, _, _ := term.GetSize(int(os.Stdout.Fd()))
	doc := strings.Builder{}
	tabs := []string{}
	for i, t := range m.Tabs {
		if i == m.activeTab {
			if m.focused {
				tabs = append(tabs, components.ActiveTabStyle.Render(t.Name))
			} else {
				tabs = append(tabs, components.ActiveTabBlurStyle.Render(t.Name))
			}
		} else {
			tabs = append(tabs, components.InactiveTabStyle.Render(t.Name))
		}
	}

	row := lipgloss.JoinHorizontal(
		lipgloss.Top,
		tabs...,
	)
	gap := components.TabGap.Render(strings.Repeat(" ", max(0, width-lipgloss.Width(row)-2)))
	row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
	doc.WriteString(row)
	return doc.String()
}
func (m Model) footerView() string {
	width, _, _ := term.GetSize(int(os.Stdout.Fd()))
	doc := strings.Builder{}
	w := lipgloss.Width

	statusKey := components.StatusSectionStyle.Render("PR")
	statusHelp := components.StatusHelpStyle.Render("? Help")
	focused := strconv.FormatBool(m.focused)
	dimensions := fmt.Sprintf("%d,%d,%d,%d", m.viewport.Width, m.viewport.Height, m.viewport.YOffset, m.viewport.YPosition)
	_, _ = focused, dimensions
	statusVal := components.StatusText.Copy().
		Width(width - w(statusKey) - w(statusHelp)).
		Render(m.context.StatusText)

	bar := lipgloss.JoinHorizontal(lipgloss.Top,
		statusKey,
		statusVal,
		statusHelp,
	)
	doc.WriteString(components.StatusBarStyle.Width(width).Render(bar))
	return doc.String()
}
func (m Model) View() string {
	var body string
	m.viewport.SetContent(m.Tabs[m.activeTab].View())
	body = m.viewport.View()
	if m.help.ShowAll {
		width, _, _ := term.GetSize(int(os.Stdout.Fd()))
		width = width / 2
		//width := lipgloss.Height(m.footerView())
		//width = m.viewport.Width / 2
		vc := m.viewport.Height/2 - lipgloss.Height(m.headerView())
		body = components.RenderHelpBox(m.help.View(m.keymap), body, width, vc, 0)
	}

	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), body, m.footerView())
}

func main() {
	p := tea.NewProgram(initializeModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
