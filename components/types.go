package components

import "github.com/dbd/ght/internal/api"

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
}
