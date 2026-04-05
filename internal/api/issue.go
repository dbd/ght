package api

import (
	"github.com/cli/go-gh/v2"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cli/go-gh/v2/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
)

type IssueCommentResult struct {
	Success bool
	Error   error
	Issue   IssueResponse
}

type IssueAssigneeResult struct {
	Success bool
	Error   error
	Issue   IssueResponse
}

type IssueCloseResult struct {
	Success bool
	Error   error
	Issue   IssueResponse
}

func GetIssues(query string) ([]IssueResponse, error) {
	client, err := api.DefaultGraphQLClient()
	if err != nil {
		return nil, err
	}
	var queryResponse struct {
		Search struct {
			Nodes []struct {
				Issue IssueResponse `graphql:"... on Issue"`
			}
			IssueCount int
		} `graphql:"search(type: ISSUE, first: 50, query: $query)"`
	}
	variables := map[string]interface{}{
		"query": graphql.String(query),
	}
	if err = client.Query("Issues", &queryResponse, variables); err != nil {
		return nil, err
	}
	var issues []IssueResponse
	for _, node := range queryResponse.Search.Nodes {
		if node.Issue.ID != "" {
			issues = append(issues, node.Issue)
		}
	}
	return issues, nil
}

func GetIssuesCmd(query string) tea.Cmd {
	return func() tea.Msg {
		issues, err := GetIssues(query)
		return Issues{Query: query, Issues: issues, Error: err}
	}
}

func GetIssue(repo string, number int64) (IssueResponse, error) {
	client, err := api.DefaultGraphQLClient()
	if err != nil {
		return IssueResponse{}, err
	}
	var queryResponse struct {
		Repository struct {
			Issue IssueResponse `graphql:"issue(number: $number)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	parts := splitRepo(repo)
	variables := map[string]interface{}{
		"owner":  graphql.String(parts[0]),
		"name":   graphql.String(parts[1]),
		"number": graphql.Int(number),
	}
	err = client.Query("Issue", &queryResponse, variables)
	if err != nil {
		return IssueResponse{}, err
	}
	return queryResponse.Repository.Issue, nil
}

func GetIssueCmd(repo string, number int64) tea.Cmd {
	return func() tea.Msg {
		issue, err := GetIssue(repo, number)
		return IssueRefresh{Issue: issue, Error: err}
	}
}

func AddIssueCommentCmd(issue IssueResponse, body string) tea.Cmd {
	return func() tea.Msg {
		args := []string{"issue", "comment", "--repo", issue.Repository.NameWithOwner, "--body", body, issueNumber(issue)}
		_, _, err := gh.Exec(args...)
		if err != nil {
			return IssueCommentResult{Success: false, Error: err, Issue: issue}
		}
		return IssueCommentResult{Success: true, Issue: issue}
	}
}

func AddIssueAssigneeCmd(issue IssueResponse, username string) tea.Cmd {
	return func() tea.Msg {
		args := []string{"issue", "edit", "--repo", issue.Repository.NameWithOwner, "--add-assignee", username, issueNumber(issue)}
		_, _, err := gh.Exec(args...)
		if err != nil {
			return IssueAssigneeResult{Success: false, Error: err, Issue: issue}
		}
		return IssueAssigneeResult{Success: true, Issue: issue}
	}
}

func CloseIssueCmd(issue IssueResponse) tea.Cmd {
	return func() tea.Msg {
		args := []string{"issue", "close", "--repo", issue.Repository.NameWithOwner, issueNumber(issue)}
		_, _, err := gh.Exec(args...)
		if err != nil {
			return IssueCloseResult{Success: false, Error: err, Issue: issue}
		}
		return IssueCloseResult{Success: true, Issue: issue}
	}
}

func ReopenIssueCmd(issue IssueResponse) tea.Cmd {
	return func() tea.Msg {
		args := []string{"issue", "reopen", "--repo", issue.Repository.NameWithOwner, issueNumber(issue)}
		_, _, err := gh.Exec(args...)
		if err != nil {
			return IssueCloseResult{Success: false, Error: err, Issue: issue}
		}
		return IssueCloseResult{Success: true, Issue: issue}
	}
}

func OpenIssueInBrowserCmd(issue IssueResponse) tea.Cmd {
	return func() tea.Msg {
		gh.Exec("issue", "view", "--web", "--repo", issue.Repository.NameWithOwner, issueNumber(issue))
		return nil
	}
}

func issueNumber(issue IssueResponse) string {
	var s string
	n := issue.Number
	if n == 0 {
		return "0"
	}
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}
