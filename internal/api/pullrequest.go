package api

import (
	"github.com/cli/go-gh/v2"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/cli/go-gh/v2/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
)

type MergeMethod string

const (
	MergeMethodMerge  MergeMethod = "merge"
	MergeMethodSquash MergeMethod = "squash"
	MergeMethodRebase MergeMethod = "rebase"
)

type MergeResult struct {
	Success bool
	Error   error
	PR      PullRequestResponse
}

func GetPullRequests(query string) (prs []PullRequestResponse) {
	client, err := api.DefaultGraphQLClient()
	if err != nil {
		log.Fatal(err)
	}
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

func GetPullRequestsCmd(query string) tea.Cmd {
	return func() tea.Msg {
		return PullRequests{Query: query, PullRequests: GetPullRequests(query)}
	}
}

func MergePullRequestCmd(pr PullRequestResponse, method MergeMethod, deleteBranch bool) tea.Cmd {
	return func() tea.Msg {
		args := []string{"pr", "merge", "--repo", pr.Repository.NameWithOwner}

		switch method {
		case MergeMethodMerge:
			args = append(args, "--merge")
		case MergeMethodSquash:
			args = append(args, "--squash")
		case MergeMethodRebase:
			args = append(args, "--rebase")
		}

		if deleteBranch {
			args = append(args, "--delete-branch")
		}

		args = append(args, pr.HeadRefName)

		_, _, err := gh.Exec(args...)
		if err != nil {
			return MergeResult{Success: false, Error: err, PR: pr}
		}
		return MergeResult{Success: true, Error: nil, PR: pr}
	}
}

type PullRequestRefresh struct {
	PR    PullRequestResponse
	Error error
}

func GetPullRequest(repo string, number int64) (PullRequestResponse, error) {
	client, err := api.DefaultGraphQLClient()
	if err != nil {
		return PullRequestResponse{}, err
	}
	var queryResponse struct {
		Repository struct {
			PullRequest PullRequestResponse `graphql:"pullRequest(number: $number)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	// Split repo into owner/name
	parts := splitRepo(repo)

	variables := map[string]interface{}{
		"owner":  graphql.String(parts[0]),
		"name":   graphql.String(parts[1]),
		"number": graphql.Int(number),
	}
	err = client.Query("PullRequest", &queryResponse, variables)
	if err != nil {
		return PullRequestResponse{}, err
	}
	return queryResponse.Repository.PullRequest, nil
}

func GetPullRequestCmd(repo string, number int64) tea.Cmd {
	return func() tea.Msg {
		pr, err := GetPullRequest(repo, number)
		return PullRequestRefresh{PR: pr, Error: err}
	}
}

func splitRepo(repo string) []string {
	for i, c := range repo {
		if c == '/' {
			return []string{repo[:i], repo[i+1:]}
		}
	}
	return []string{"", repo}
}
