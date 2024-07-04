package components

import (
	"github.com/charmbracelet/bubbles/help"
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
