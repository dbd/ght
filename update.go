package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dbd/ght/components"
	"github.com/dbd/ght/components/pullRequestDetail"
	"github.com/dbd/ght/components/pullRequestSearch"
	"github.com/dbd/ght/components/tab"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	activeTab := m.Tabs[m.activeTab]
	switch msg := msg.(type) {
	case pullRequestSearch.OpenPR:
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
			t := tab.NewModel(m.context, msg.PR.Title, pullRequestDetail.NewModel(msg.PR, m.context))
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

		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - verticalMarginHeight
		m.context.ViewportWidth = m.viewport.Width
		m.context.ViewportHeight = m.viewport.Height
		m.context.ViewportYOffset = m.viewport.YOffset
		m.context.ViewportYPosition = m.viewport.YPosition
		activeTab, cmd := activeTab.Update(msg)
		cmds = append(cmds, cmd)
		m.Tabs[m.activeTab] = activeTab
	case tea.KeyMsg:
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
			activeTab, cmd := activeTab.Update(msg)
			cmds = append(cmds, cmd)
			m.Tabs[m.activeTab] = activeTab
		} else {
			switch {
			case key.Matches(msg, m.context.KeyMap.Help):
				m.showHelp = !m.showHelp
			case key.Matches(msg, m.context.KeyMap.Down):
				m.focused = false
			case key.Matches(msg, m.context.KeyMap.Exit):
				return m, tea.Quit
			case key.Matches(msg, m.context.KeyMap.Close):
				var tt []tab.Model
				for counter, tab := range m.Tabs {
					if counter != m.activeTab {
						tt = append(tt, tab)
					}
				}
				if len(tt) == 0 {
					m.context.StatusText = "Unable to close last tab. Exit instead."
					break
				}
				m.Tabs = tt
				if m.activeTab < len(m.Tabs)-1 {
					m.activeTab++
				} else if m.activeTab != 0 {
					m.activeTab--
				}
			case key.Matches(msg, m.context.KeyMap.Left):
				if m.activeTab > 0 {
					m.activeTab--
				}
			case key.Matches(msg, m.context.KeyMap.Right):
				if m.activeTab < len(m.Tabs)-1 {
					m.activeTab++
				}
			}
		}
	default:
		tabs := []tab.Model{}
		for _, tab := range m.Tabs {
			tab, cmd := tab.Update(msg)
			cmds = append(cmds, cmd)
			tabs = append(tabs, tab)
		}
		m.Tabs = tabs
	}
	return m, tea.Batch(cmds...)
}

func (m Model) canEnterCommandMode() bool {
	for _, tab := range m.Tabs {
		if tab.IsInTextInput() {
			return false
		}
	}
	return true
}
