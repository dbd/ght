package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dbd/ght/components"
	"github.com/dbd/ght/components/issueDetail"
	"github.com/dbd/ght/components/issueSearch"
	"github.com/dbd/ght/components/milestoneDetail"
	"github.com/dbd/ght/components/milestoneList"
	"github.com/dbd/ght/components/pullRequestDetail"
	"github.com/dbd/ght/components/pullRequestSearch"
	"github.com/dbd/ght/components/tab"
	"github.com/dbd/ght/internal/api"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Handle PR setup dialog
	if m.needsSetup {
		switch msg := msg.(type) {
		case components.ConfigCreated:
			return initializeModel(), tea.ClearScreen
		case tea.KeyMsg:
			sd, sdCmd := m.setupDialog.Update(msg)
			m.setupDialog = sd.(*components.SetupDialogModel)
			return m, sdCmd
		}
		return m, nil
	}

	// Handle issue setup dialog (when switching to issue mode with no searches)
	if m.issueNeedsSetup {
		switch msg := msg.(type) {
		case components.IssueConfigCreated:
			m.issueNeedsSetup = false
			// Reload issue tabs from config
			config := components.GetConfig()
			m.config = config
			newCtx := m.context
			for _, search := range config.Issue.Searches {
				t := tab.NewModel(newCtx, search.Name, issueSearch.NewModel([]api.IssueResponse{}, search.Query, newCtx))
				m.issueTabs = append(m.issueTabs, t)
			}
			var initCmds []tea.Cmd
			for _, t := range m.issueTabs {
				initCmds = append(initCmds, t.Init())
			}
			return m, tea.Batch(initCmds...)
		case tea.KeyMsg:
			sd, sdCmd := m.setupDialog.Update(msg)
			m.setupDialog = sd.(*components.SetupDialogModel)
			return m, sdCmd
		}
		return m, nil
	}

	tabs := m.activeTabs()
	if len(tabs) == 0 {
		// No tabs in current mode — handle basic keys
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if key.Matches(msg, m.context.KeyMap.Exit) {
				return m, tea.Quit
			}
			m = m.handleModeKeys(msg, &cmds)
		case tea.WindowSizeMsg:
			m = m.handleWindowSize(msg, &cmds)
		}
		return m, tea.Batch(cmds...)
	}

	activeTab := tabs[m.activeTabIdx()]

	switch msg := msg.(type) {
	case components.CmdQuit:
		return m, tea.Quit

	case pullRequestSearch.OpenPR:
		tabs := m.prTabs
		for _, t := range tabs {
			if t.Name == msg.PR.Title {
				m.prActiveTab = len(tabs) - 1
				m.focused = true
				return m, nil
			}
		}
		t := tab.NewModel(m.context, msg.PR.Title, pullRequestDetail.NewModel(msg.PR, m.context))
		m.prTabs = append(m.prTabs, t)
		m.prActiveTab = len(m.prTabs) - 1
		m.focused = true

	case issueSearch.OpenIssue:
		// Switch to issue mode if needed
		m.mode = "issue"
		for _, t := range m.issueTabs {
			if t.Name == msg.Issue.Title {
				m.issueActiveTab = len(m.issueTabs) - 1
				m.focused = true
				return m, nil
			}
		}
		t := tab.NewModel(m.context, msg.Issue.Title, issueDetail.NewModel(msg.Issue, m.context))
		m.issueTabs = append(m.issueTabs, t)
		m.issueActiveTab = len(m.issueTabs) - 1
		m.focused = true

	case components.OpenMilestoneByNumber:
		m.mode = "issue"
		tabName := fmt.Sprintf("Milestone #%d", msg.Number)
		for _, t := range m.issueTabs {
			if t.Name == tabName {
				m.issueActiveTab = len(m.issueTabs) - 1
				m.focused = true
				return m, nil
			}
		}
		t := tab.NewModel(m.context, tabName, milestoneDetail.NewModel(msg.Repo, msg.Number, m.context))
		m.issueTabs = append(m.issueTabs, t)
		m.issueActiveTab = len(m.issueTabs) - 1
		m.focused = true
		cmds = append(cmds, func() tea.Cmd {
			allTabs := m.issueTabs
			return allTabs[len(allTabs)-1].Init()
		}())

	case components.CmdSwitchMode:
		if msg.Mode != m.mode {
			m.mode = msg.Mode
			m.focused = true
			if m.mode == "issue" && len(m.issueTabs) == 0 {
				m.issueNeedsSetup = true
				m.setupDialog = components.NewIssueSetupDialogModel(m.context)
			}
		}

	case components.Blur:
		m.focused = true

	case components.CmdNewTab:
		t := tab.NewModel(m.context, "New Search", pullRequestSearch.NewModel([]api.PullRequestResponse{}, "", m.context))
		m.prTabs = append(m.prTabs, t)
		m.prActiveTab = len(m.prTabs) - 1
		m.focused = true

	case components.CmdNewIssueTab:
		m.mode = "issue"
		t := tab.NewModel(m.context, "New Search", issueSearch.NewModel([]api.IssueResponse{}, "", m.context))
		m.issueTabs = append(m.issueTabs, t)
		m.issueActiveTab = len(m.issueTabs) - 1
		m.focused = true

	case components.CmdMilestones:
		m.mode = "issue"
		tabName := "Milestones: " + msg.Repo
		for _, t := range m.issueTabs {
			if t.Name == tabName {
				m.issueActiveTab = len(m.issueTabs) - 1
				m.focused = true
				return m, nil
			}
		}
		t := tab.NewModel(m.context, tabName, milestoneList.NewModel(msg.Repo, m.context))
		m.issueTabs = append(m.issueTabs, t)
		m.issueActiveTab = len(m.issueTabs) - 1
		m.focused = true
		cmds = append(cmds, func() tea.Cmd {
			allTabs := m.issueTabs
			return allTabs[len(allTabs)-1].Init()
		}())

	case components.CmdSaveTab:
		activeTab := m.activeTabs()[m.activeTabIdx()]
		if searchPage, ok := activeTab.GetPage().(*pullRequestSearch.Model); ok {
			query := searchPage.GetQuery()
			if query == "" {
				m.context.StatusText = "Cannot save tab: no query set"
			} else if err := components.SaveSearch(msg.Name, query); err != nil {
				m.context.StatusText = fmt.Sprintf("Failed to save: %v", err)
			} else {
				m.context.StatusText = fmt.Sprintf("Saved '%s'", msg.Name)
				m.prTabs[m.prActiveTab].Name = msg.Name
			}
		} else if issuePage, ok := activeTab.GetPage().(*issueSearch.Model); ok {
			query := issuePage.GetQuery()
			if query == "" {
				m.context.StatusText = "Cannot save tab: no query set"
			} else if err := components.SaveIssueSearch(msg.Name, query); err != nil {
				m.context.StatusText = fmt.Sprintf("Failed to save: %v", err)
			} else {
				m.context.StatusText = fmt.Sprintf("Saved '%s'", msg.Name)
				m.issueTabs[m.issueActiveTab].Name = msg.Name
			}
		} else if msPage, ok := activeTab.GetPage().(*milestoneList.Model); ok {
			repo := msPage.GetRepo()
			if repo == "" {
				m.context.StatusText = "Cannot save tab: no repo set"
			} else if err := components.SaveMilestoneRepo(msg.Name, repo); err != nil {
				m.context.StatusText = fmt.Sprintf("Failed to save: %v", err)
			} else {
				m.context.StatusText = fmt.Sprintf("Saved '%s'", msg.Name)
				m.issueTabs[m.issueActiveTab].Name = msg.Name
			}
		} else {
			m.context.StatusText = "Cannot save this tab type"
		}

	case components.CmdRefresh:
		activeTab := m.activeTabs()[m.activeTabIdx()]
		if searchPage, ok := activeTab.GetPage().(*pullRequestSearch.Model); ok {
			query := searchPage.GetQuery()
			if query == "" {
				m.context.StatusText = "Cannot refresh: no query set"
			} else {
				m.context.StatusText = "Refreshing search..."
				cmds = append(cmds, api.GetPullRequestsCmd(query))
			}
		} else if prPage, ok := activeTab.GetPage().(*pullRequestDetail.Model); ok {
			pr := prPage.GetPullRequest()
			m.context.StatusText = "Refreshing PR..."
			cmds = append(cmds, api.GetPullRequestCmd(pr.Repository.NameWithOwner, pr.Number))
		} else if issueSPage, ok := activeTab.GetPage().(*issueSearch.Model); ok {
			query := issueSPage.GetQuery()
			if query == "" {
				m.context.StatusText = "Cannot refresh: no query set"
			} else {
				m.context.StatusText = "Refreshing issues..."
				cmds = append(cmds, api.GetIssuesCmd(query))
			}
		} else if issueDPage, ok := activeTab.GetPage().(*issueDetail.Model); ok {
			issue := issueDPage.GetIssue()
			m.context.StatusText = "Refreshing issue..."
			cmds = append(cmds, api.GetIssueCmd(issue.Repository.NameWithOwner, issue.Number))
		} else if msDPage, ok := activeTab.GetPage().(*milestoneDetail.Model); ok {
			m.context.StatusText = "Refreshing milestone..."
			cmds = append(cmds, api.GetMilestoneCmd(msDPage.GetRepo(), msDPage.GetNumber()))
		} else if msLPage, ok := activeTab.GetPage().(*milestoneList.Model); ok {
			m.context.StatusText = "Refreshing milestones..."
			cmds = append(cmds, api.GetMilestonesCmd(msLPage.GetRepo()))
		} else {
			m.context.StatusText = "Cannot refresh this tab type"
		}

	case components.CmdAddAssignee:
		activeTab := m.activeTabs()[m.activeTabIdx()]
		if prPage, ok := activeTab.GetPage().(*pullRequestDetail.Model); ok {
			pr := prPage.GetPullRequest()
			m.context.StatusText = "Adding assignee..."
			cmds = append(cmds, api.AddAssigneeCmd(pr, msg.Username))
		} else {
			m.context.StatusText = "Can only add assignees on PR detail tabs"
		}

	case components.CmdAddReviewer:
		activeTab := m.activeTabs()[m.activeTabIdx()]
		if prPage, ok := activeTab.GetPage().(*pullRequestDetail.Model); ok {
			pr := prPage.GetPullRequest()
			m.context.StatusText = "Adding reviewer..."
			cmds = append(cmds, api.AddReviewerCmd(pr, msg.Username))
		} else {
			m.context.StatusText = "Can only add reviewers on PR detail tabs"
		}

	case components.CmdComment:
		activeTab := m.activeTabs()[m.activeTabIdx()]
		if prPage, ok := activeTab.GetPage().(*pullRequestDetail.Model); ok {
			pr := prPage.GetPullRequest()
			m.context.StatusText = "Adding comment..."
			cmds = append(cmds, api.AddCommentCmd(pr, msg.Body))
		} else {
			m.context.StatusText = "Can only comment on PR detail tabs"
		}

	case components.CmdApprove:
		activeTab := m.activeTabs()[m.activeTabIdx()]
		if prPage, ok := activeTab.GetPage().(*pullRequestDetail.Model); ok {
			pr := prPage.GetPullRequest()
			m.context.StatusText = "Approving PR..."
			cmds = append(cmds, api.SubmitReviewCmd(pr, api.ReviewActionApprove, msg.Body))
		} else {
			m.context.StatusText = "Can only approve PR detail tabs"
		}

	case components.CmdRequestChanges:
		activeTab := m.activeTabs()[m.activeTabIdx()]
		if prPage, ok := activeTab.GetPage().(*pullRequestDetail.Model); ok {
			pr := prPage.GetPullRequest()
			m.context.StatusText = "Requesting changes..."
			cmds = append(cmds, api.SubmitReviewCmd(pr, api.ReviewActionRequestChanges, msg.Body))
		} else {
			m.context.StatusText = "Can only request changes on PR detail tabs"
		}

	case components.CmdHelp:
		m.helpDialog.Focus()

	case components.CmdMerge:
		activeTab := m.activeTabs()[m.activeTabIdx()]
		if prPage, ok := activeTab.GetPage().(*pullRequestDetail.Model); ok {
			_ = prPage
			// Merge is handled inside the PR detail page via its dialog
		}

	case tea.WindowSizeMsg:
		m = m.handleWindowSize(msg, &cmds)

	case tea.KeyMsg:
		if m.helpDialog.Focused() {
			if key.Matches(msg, m.context.KeyMap.Exit) {
				m.helpDialog.Blur()
				return m, nil
			}
			hd, hdCmd := m.helpDialog.Update(msg)
			m.helpDialog = *hd.(*components.HelpDialogModel)
			cmds = append(cmds, hdCmd)
			return m, tea.Batch(cmds...)
		}
		if key.Matches(msg, m.context.KeyMap.Suspend) {
			return m, tea.Suspend
		}
		if m.command.Focused() {
			command, cmd := m.command.Update(msg)
			cmds = append(cmds, cmd)
			m.command = command
			m.context.StatusText = fmt.Sprintf(":%s█", m.command.Value())
			if key.Matches(msg, m.context.KeyMap.Enter) {
				m.command.Blur()
				m.context.StatusText = ""
				cmd := m.sendCommandMessage(m.command.Value())
				cmds = append(cmds, cmd)
				m.command.SetValue("")
				break
			}
			break
		}
		if key.Matches(msg, m.context.KeyMap.Leader) && m.canEnterCommandMode() {
			m.command.Focus()
			m.context.StatusText = ":█"
		}
		if !m.focused {
			tabs := m.activeTabs()
			idx := m.activeTabIdx()
			updatedTab, cmd := tabs[idx].Update(msg)
			cmds = append(cmds, cmd)
			tabs[idx] = updatedTab
			m.setActiveTabs(tabs)
		} else {
			m = m.handleModeKeys(msg, &cmds)
		}

	default:
		// Broadcast to all tabs in both modes so async responses reach their tab
		// regardless of which mode is currently active.
		updatedPR := []tab.Model{}
		for _, t := range m.prTabs {
			updated, cmd := t.Update(msg)
			cmds = append(cmds, cmd)
			updatedPR = append(updatedPR, updated)
		}
		m.prTabs = updatedPR
		updatedIssue := []tab.Model{}
		for _, t := range m.issueTabs {
			updated, cmd := t.Update(msg)
			cmds = append(cmds, cmd)
			updatedIssue = append(updatedIssue, updated)
		}
		m.issueTabs = updatedIssue
	}

	_ = activeTab
	return m, tea.Batch(cmds...)
}

func (m Model) handleModeKeys(msg tea.KeyMsg, cmds *[]tea.Cmd) Model {
	tabs := m.activeTabs()
	idx := m.activeTabIdx()
	switch {
	case key.Matches(msg, m.context.KeyMap.Help):
		m.showHelp = !m.showHelp
	case key.Matches(msg, m.context.KeyMap.Exit):
		if m.showHelp {
			m.showHelp = false
		} else {
			*cmds = append(*cmds, tea.Quit)
		}
	case msg.String() == "I":
		if m.mode != "issue" {
			m.mode = "issue"
			m.focused = true
			if len(m.issueTabs) == 0 {
				m.issueNeedsSetup = true
				m.setupDialog = components.NewIssueSetupDialogModel(m.context)
			}
		}
	case msg.String() == "P":
		if m.mode != "pr" {
			m.mode = "pr"
			m.focused = true
		}
	case key.Matches(msg, m.context.KeyMap.Down):
		m.focused = false
	case key.Matches(msg, m.context.KeyMap.Close):
		var tt []tab.Model
		for i, t := range tabs {
			if i != idx {
				tt = append(tt, t)
			}
		}
		if len(tt) == 0 {
			m.context.StatusText = "Unable to close last tab. Exit instead."
			break
		}
		m.setActiveTabs(tt)
		newIdx := idx
		if newIdx >= len(tt) {
			newIdx = len(tt) - 1
		}
		m.setActiveTabIdx(newIdx)
	case key.Matches(msg, m.context.KeyMap.Left):
		if idx > 0 {
			m.setActiveTabIdx(idx - 1)
		}
	case key.Matches(msg, m.context.KeyMap.Right):
		if idx < len(tabs)-1 {
			m.setActiveTabIdx(idx + 1)
		}
	}
	return m
}

func (m Model) handleWindowSize(msg tea.WindowSizeMsg, cmds *[]tea.Cmd) Model {
	headerHeight := lipgloss.Height(m.headerView())
	footerHeight := lipgloss.Height(m.footerView())
	verticalMarginHeight := headerHeight + footerHeight

	m.viewport.Width = msg.Width
	m.viewport.Height = msg.Height - verticalMarginHeight
	m.context.ViewportWidth = m.viewport.Width
	m.context.ViewportHeight = m.viewport.Height
	m.context.ViewportYOffset = m.viewport.YOffset
	m.context.ViewportYPosition = m.viewport.YPosition

	tabs := m.activeTabs()
	if len(tabs) > 0 {
		idx := m.activeTabIdx()
		updated, cmd := tabs[idx].Update(msg)
		*cmds = append(*cmds, cmd)
		tabs[idx] = updated
		m.setActiveTabs(tabs)
	}
	return m
}

func (m Model) canEnterCommandMode() bool {
	for _, t := range m.activeTabs() {
		if t.IsInTextInput() {
			return false
		}
	}
	return true
}
