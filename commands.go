package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type CmdMerge struct {
	Org           string
	Repo          string
	PullRequestId int
}

type CmdRefresh struct{}

var cmdMap = map[string]interface{}{
	"merge":   CmdMerge{},
	"refresh": CmdRefresh{},
}

func (m Model) sendCommandMessage(command string) tea.Cmd {
	msg, ok := cmdMap[command]
	if !ok {
		m.context.StatusText = fmt.Sprintf("Unknown command: %s", command)
		return nil
	}
	switch msg := msg.(type) {
	case CmdMerge:
		m.context.StatusText = fmt.Sprintf("Merge Command, %s", msg.Org)
		return func() tea.Msg { return CmdMerge{} }

	}
	return nil
}
