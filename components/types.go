package components

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

type Blur bool

type Context struct {
	ViewportWidth     int
	ViewportHeight    int
	ViewportYOffset   int
	ViewportYPosition int
	StatusText        string
	KeyMap            KeyMap
	Help              help.Model
}

type Page interface {
	Init() tea.Cmd
	Update(tea.Msg) (Page, tea.Cmd)
	View() string
	Blur() tea.Msg
	Focus() tea.Msg
	ToggleHelp() tea.Msg
	IsInTextInput() bool
}
