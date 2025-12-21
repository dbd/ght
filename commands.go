package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dbd/ght/components"
)

var cmdMap = map[string]interface{}{
	"merge":    components.CmdMerge{},
	"refresh":  components.CmdRefresh{},
	"newtab":   components.CmdNewTab{},
	"save-tab": components.CmdSaveTab{},
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
		m.context.StatusText = fmt.Sprintf("Unknown command: %s", cmdName)
		return nil
	}
	switch msg.(type) {
	case components.CmdMerge:
		m.context.StatusText = fmt.Sprintf("Merge Command")
		return func() tea.Msg { return components.CmdMerge{} }
	case components.CmdNewTab:
		return func() tea.Msg { return components.CmdNewTab{} }
	case components.CmdSaveTab:
		if cmdArg == "" {
			m.context.StatusText = "Usage: save-tab <name>"
			return nil
		}
		return func() tea.Msg { return components.CmdSaveTab{Name: cmdArg} }
	}
	return nil
}
