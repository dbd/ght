package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dbd/ght/components"
	"golang.org/x/term"
)

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
	if m.showHelp {
		width, _, _ := term.GetSize(int(os.Stdout.Fd()))
		width = width / 2
		vc := m.viewport.Height/2 - lipgloss.Height(m.headerView())
		body = components.RenderHelpBox(m.context.Help.FullHelpView(fullHelp), body, width, vc, 0)
	}

	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), body, m.footerView())
}
