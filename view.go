package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dbd/ght/components"
	"golang.org/x/term"
)

func (m Model) headerView() string {
	width, _, _ := term.GetSize(int(os.Stdout.Fd()))
	doc := strings.Builder{}
	tabs := m.activeTabs()
	idx := m.activeTabIdx()
	tabViews := []string{}
	for i, t := range tabs {
		if i == idx {
			if m.focused {
				tabViews = append(tabViews, components.ActiveTabStyle.Render(t.Name))
			} else {
				tabViews = append(tabViews, components.ActiveTabBlurStyle.Render(t.Name))
			}
		} else {
			tabViews = append(tabViews, components.InactiveTabStyle.Render(t.Name))
		}
	}

	row := lipgloss.JoinHorizontal(
		lipgloss.Top,
		tabViews...,
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

	modeLabel := "PR"
	if m.mode == "issue" {
		modeLabel = "Issue"
	}

	statusKey := components.StatusSectionStyle.Render(modeLabel)
	statusHelp := components.StatusHelpStyle.Render("? Help")
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
	// PR setup dialog
	if m.needsSetup {
		width, height, _ := term.GetSize(int(os.Stdout.Fd()))
		background := strings.Repeat(strings.Repeat(" ", width)+"\n", height)
		return components.RenderCenteredOverlay(m.setupDialog.View(), background, width/2, height/2)
	}

	// Issue setup dialog
	if m.issueNeedsSetup {
		width, height, _ := term.GetSize(int(os.Stdout.Fd()))
		background := strings.Repeat(strings.Repeat(" ", width)+"\n", height)
		return components.RenderCenteredOverlay(m.setupDialog.View(), background, width/2, height/2)
	}

	tabs := m.activeTabs()
	if len(tabs) == 0 {
		return fmt.Sprintf("%s\n%s\n%s", m.headerView(), "No tabs. Use :new-issue-tab or :newtab", m.footerView())
	}

	var body string
	m.viewport.SetContent(tabs[m.activeTabIdx()].View())
	body = m.viewport.View()
	if m.showHelp {
		width, _, _ := term.GetSize(int(os.Stdout.Fd()))
		width = width / 2
		vc := m.viewport.Height/2 - lipgloss.Height(m.headerView())
		body = components.RenderHelpBox(m.context.Help.FullHelpView(fullHelp), body, width, vc, 0)
	}
	if m.helpDialog.Focused() {
		width, _, _ := term.GetSize(int(os.Stdout.Fd()))
		vc := m.viewport.Height / 2
		hc := width / 2
		body = components.RenderCenteredOverlay(m.helpDialog.View(), body, hc, vc)
	}

	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), body, m.footerView())
}
