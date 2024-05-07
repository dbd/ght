package api

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/cli/go-gh/v2/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
)

func GetPullRequests(query, assignee string) (prs []PullRequestResponse) {
	client, err := api.DefaultGraphQLClient()
	if err != nil {
		log.Fatal(err)
	}
	query = query + " assignee:" + assignee
	var queryResponse struct {
		Search struct {
			Nodes []struct {
				PullRequest PullRequestResponse `graphql:"... on PullRequest"`
			}
			IssueCount int
		} `graphql:"search(type: ISSUE, first: 10, query: $query)"`
	}
	variables := map[string]interface{}{
		"query": graphql.String(query),
	}
	err = client.Query("PullRequests", &queryResponse, variables)
	if err != nil {
		log.Fatal(err)
	}
	for _, pr := range queryResponse.Search.Nodes {
		prs = append(prs, pr.PullRequest)

	}
	return

}

func GetPullRequestsCmd(query, assignee string) tea.Cmd {
	return func() tea.Msg {
		return PullRequests{Query: query, PullRequests: GetPullRequests(query, assignee)}
	}
}
