package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dbd/ght/components"
	"github.com/dbd/ght/components/pullRequestSearch"
	"github.com/dbd/ght/components/tab"
	"github.com/dbd/ght/internal/api"
	"golang.org/x/term"
)

type Model struct {
	Tabs      []tab.Model
	activeTab int
	config    components.Config
	viewport  viewport.Model
	ready     bool
	context   *components.Context
	focused   bool
	showHelp  bool
}

var fullHelp = [][]key.Binding{
	{components.DefaultKeyMap.Up, components.DefaultKeyMap.Down, components.DefaultKeyMap.Left, components.DefaultKeyMap.Right},
	{components.DefaultKeyMap.Help, components.DefaultKeyMap.Close, components.DefaultKeyMap.Exit},
}

func initializeModel() Model {
	config := components.SetupConfig()
	h := help.New()
	h.Styles = components.HelpStyles
	ctx := components.Context{KeyMap: components.DefaultKeyMap, Help: h}
	m := Model{
		config:   config,
		context:  &ctx,
		focused:  true,
		showHelp: false,
	}
	headerHeight := lipgloss.Height(m.headerView())
	footerHeight := lipgloss.Height(m.footerView())
	verticalMarginHeight := headerHeight + footerHeight

	mw, mh, _ := term.GetSize(int(os.Stdout.Fd()))
	m.viewport = viewport.New(mw, mh-verticalMarginHeight)
	m.viewport.YPosition = headerHeight - 1
	m.context.ViewportWidth = m.viewport.Width
	m.context.ViewportHeight = m.viewport.Height
	m.context.ViewportYOffset = m.viewport.YOffset
	m.context.ViewportYPosition = m.viewport.YPosition
	tabs := []tab.Model{}
	for _, search := range config.Pr.Searches {
		t := tab.NewModel(&ctx, search.Name)
		t.Page = pullRequestSearch.NewModel([]api.PullRequestResponse{}, search.Query, &ctx)
		tabs = append(tabs, t)
	}
	m.Tabs = tabs
	return m
}

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, tab := range m.Tabs {
		cmds = append(cmds, tab.Init())
	}
	return tea.Batch(cmds...)
}

func main() {
	p := tea.NewProgram(initializeModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
