package api

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/cli/go-gh/v2/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
)

func GetMilestones(repo string) (milestones []MilestoneListResponse, err error) {
	client, err := api.DefaultGraphQLClient()
	if err != nil {
		return nil, err
	}
	var queryResponse struct {
		Repository struct {
			Milestones struct {
				Nodes []MilestoneListResponse
			} `graphql:"milestones(first: 50, states: [OPEN, CLOSED], orderBy: { field: DUE_DATE, direction: ASC })"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	parts := splitRepo(repo)
	variables := map[string]interface{}{
		"owner": graphql.String(parts[0]),
		"name":  graphql.String(parts[1]),
	}
	err = client.Query("Milestones", &queryResponse, variables)
	if err != nil {
		return nil, err
	}
	return queryResponse.Repository.Milestones.Nodes, nil
}

func GetMilestonesCmd(repo string) tea.Cmd {
	return func() tea.Msg {
		ms, err := GetMilestones(repo)
		if err != nil {
			log.Error("GetMilestones", "error", err)
		}
		return Milestones{Repo: repo, Milestones: ms, Error: err}
	}
}

func GetMilestone(repo string, number int64) (MilestoneResponse, error) {
	client, err := api.DefaultGraphQLClient()
	if err != nil {
		return MilestoneResponse{}, err
	}
	var queryResponse struct {
		Repository struct {
			Milestone MilestoneResponse `graphql:"milestone(number: $number)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	parts := splitRepo(repo)
	variables := map[string]interface{}{
		"owner":  graphql.String(parts[0]),
		"name":   graphql.String(parts[1]),
		"number": graphql.Int(number),
	}
	err = client.Query("Milestone", &queryResponse, variables)
	if err != nil {
		return MilestoneResponse{}, err
	}
	return queryResponse.Repository.Milestone, nil
}

func GetMilestoneCmd(repo string, number int64) tea.Cmd {
	return func() tea.Msg {
		m, err := GetMilestone(repo, number)
		return MilestoneRefresh{Milestone: m, Error: err}
	}
}
