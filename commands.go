package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dbd/ght/components"
)

var cmdMap = map[string]interface{}{
	"merge":   components.CmdMerge{},
	"refresh": components.CmdRefresh{},
}

func (m Model) sendCommandMessage(command string) tea.Cmd {
	msg, ok := cmdMap[command]
	if !ok {
		m.context.StatusText = fmt.Sprintf("Unknown command: %s", command)
		return nil
	}
	switch msg := msg.(type) {
	case components.CmdMerge:
		m.context.StatusText = fmt.Sprintf("Merge Command, %s", msg.Org)
		return func() tea.Msg { return components.CmdMerge{} }

	}
	return nil
}
