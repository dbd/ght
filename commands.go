package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dbd/ght/components"
)

var cmdMap = map[string]interface{}{
	"merge":           components.CmdMerge{},
	"refresh":         components.CmdRefresh{},
	"newtab":          components.CmdNewTab{},
	"save-tab":        components.CmdSaveTab{},
	"add-assignee":    components.CmdAddAssignee{},
	"add-reviewer":    components.CmdAddReviewer{},
	"comment":         components.CmdComment{},
	"approve":         components.CmdApprove{},
	"request-changes": components.CmdRequestChanges{},
	"help":            components.CmdHelp{},
	"quit":            components.CmdQuit{},
}

func (m Model) sendCommandMessage(command string) tea.Cmd {
	parts := strings.SplitN(command, " ", 2)
	cmdName := parts[0]
	var cmdArg string
	if len(parts) > 1 {
		cmdArg = parts[1]
	}

	msg, ok := cmdMap[cmdName]
	if !ok {
		count := 0
		for key, _ := range cmdMap {
			if strings.HasPrefix(key, cmdName) {
				msg = cmdMap[key]
				ok = true
				count += 1
			}
		}
		if count > 1 {
			m.context.StatusText = fmt.Sprintf("Ambiguous command: %s", cmdName)
			return nil
		}
		if !ok {
			m.context.StatusText = fmt.Sprintf("Unknown command: %s", cmdName)
			return nil
		}
	}
	switch msg.(type) {
	case components.CmdMerge:
		return func() tea.Msg { return components.CmdMerge{} }
	case components.CmdNewTab:
		return func() tea.Msg { return components.CmdNewTab{} }
	case components.CmdSaveTab:
		if cmdArg == "" {
			m.context.StatusText = "Usage: save-tab <name>"
			return nil
		}
		return func() tea.Msg { return components.CmdSaveTab{Name: cmdArg} }
	case components.CmdRefresh:
		return func() tea.Msg { return components.CmdRefresh{} }
	case components.CmdAddAssignee:
		if cmdArg == "" {
			m.context.StatusText = "Usage: add-assignee <username>"
			return nil
		}
		return func() tea.Msg { return components.CmdAddAssignee{Username: cmdArg} }
	case components.CmdAddReviewer:
		if cmdArg == "" {
			m.context.StatusText = "Usage: add-reviewer <username>"
			return nil
		}
		return func() tea.Msg { return components.CmdAddReviewer{Username: cmdArg} }
	case components.CmdComment:
		if cmdArg == "" {
			m.context.StatusText = "Usage: comment <message>"
			return nil
		}
		return func() tea.Msg { return components.CmdComment{Body: cmdArg} }
	case components.CmdApprove:
		return func() tea.Msg { return components.CmdApprove{Body: cmdArg} }
	case components.CmdRequestChanges:
		if cmdArg == "" {
			m.context.StatusText = "Usage: request-changes <message>"
			return nil
		}
		return func() tea.Msg { return components.CmdRequestChanges{Body: cmdArg} }
	case components.CmdHelp:
		return func() tea.Msg { return components.CmdHelp{} }
	case components.CmdQuit:
		return func() tea.Msg { return components.CmdQuit{} }
	}
	return nil
}
