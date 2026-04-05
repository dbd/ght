package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dbd/ght/components"
	"github.com/dbd/ght/components/issueSearch"
	"github.com/dbd/ght/components/milestoneList"
	"github.com/dbd/ght/components/pullRequestSearch"
	"github.com/dbd/ght/components/tab"
	"github.com/dbd/ght/internal/api"
	"golang.org/x/term"
)

type Model struct {
	prTabs        []tab.Model
	issueTabs     []tab.Model
	prActiveTab   int
	issueActiveTab int
	mode          string // "pr" or "issue"
	config        components.Config
	viewport      viewport.Model
	ready         bool
	context       *components.Context
	command       textinput.Model
	focused       bool
	showHelp      bool
	helpDialog    components.HelpDialogModel
	setupDialog   *components.SetupDialogModel
	needsSetup    bool
	issueNeedsSetup bool
}

// activeTabs returns the tabs for the current mode.
func (m *Model) activeTabs() []tab.Model {
	if m.mode == "issue" {
		return m.issueTabs
	}
	return m.prTabs
}

// setActiveTabs sets the tabs for the current mode.
func (m *Model) setActiveTabs(tabs []tab.Model) {
	if m.mode == "issue" {
		m.issueTabs = tabs
	} else {
		m.prTabs = tabs
	}
}

// activeTabIdx returns the active tab index for the current mode.
func (m *Model) activeTabIdx() int {
	if m.mode == "issue" {
		return m.issueActiveTab
	}
	return m.prActiveTab
}

// setActiveTabIdx sets the active tab index for the current mode.
func (m *Model) setActiveTabIdx(i int) {
	if m.mode == "issue" {
		m.issueActiveTab = i
	} else {
		m.prActiveTab = i
	}
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
		config:  config,
		context: &ctx,
		focused: true,
		mode:    "pr",
		command: textinput.New(),
	}

	// Check if we need to run PR setup
	if len(config.Pr.Searches) == 0 {
		m.needsSetup = true
		m.setupDialog = components.NewSetupDialogModel(&ctx)
		return m
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

	// Initialize PR tabs
	prTabs := []tab.Model{}
	for _, search := range config.Pr.Searches {
		t := tab.NewModel(&ctx, search.Name, pullRequestSearch.NewModel([]api.PullRequestResponse{}, search.Query, &ctx))
		prTabs = append(prTabs, t)
	}
	m.prTabs = prTabs

	// Initialize Issue tabs
	issueTabs := []tab.Model{}
	for _, search := range config.Issue.Searches {
		t := tab.NewModel(&ctx, search.Name, issueSearch.NewModel([]api.IssueResponse{}, search.Query, &ctx))
		issueTabs = append(issueTabs, t)
	}
	for _, ms := range config.Issue.Milestones {
		t := tab.NewModel(&ctx, ms.Name, milestoneList.NewModel(ms.Repo, &ctx))
		issueTabs = append(issueTabs, t)
	}
	m.issueTabs = issueTabs

	m.helpDialog = *components.NewHelpDialogModel(&ctx)
	return m
}

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, t := range m.prTabs {
		cmds = append(cmds, t.Init())
	}
	for _, t := range m.issueTabs {
		cmds = append(cmds, t.Init())
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
