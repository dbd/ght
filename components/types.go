package components

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/dbd/ght/internal/api"
)

type Blur bool

type OpenPR struct {
	PR api.PullRequestResponse
}

type Context struct {
	ViewportWidth     int
	ViewportHeight    int
	ViewportYOffset   int
	ViewportYPosition int
	StatusText        string
	KeyMap            KeyMap
	Help              help.Model
}
