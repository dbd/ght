package main

import (
	"fmt"
	"github.com/dbd/ght/internal/api"
)

func main() {
	prs := api.GetPullRequests("assignee:dbd")
	for _, pr := range prs {
		fmt.Printf("%v+\n", pr)
	}
}
