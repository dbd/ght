package main

import (
	"github.com/dbd/ght/internal/api"
)

func main() {
	api.GetPullRequests("assignee:dbd")
}
