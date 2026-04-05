package main

import (
	"fmt"
	"os"

	"github.com/dbd/ght/components"
	"github.com/dbd/ght/components/pullRequestDetail"
	"github.com/dbd/ght/internal/api"
	"golang.org/x/term"
)

var staticPR = api.PullRequestResponse{
	ID:           "PR_kwDOJx1abc",
	Number:       89,
	Title:        "[pr] add CI information to pull request detail view",
	Body:         "## Summary\n\nAdds CI check status to the PR detail header. Each check run and status context is shown with a colored symbol (✓ failure ✗ pending ●) alongside its name.\n\n## Changes\n\n- `formatCIChecks` renders the rollup state and individual check conclusions\n- `formatCIStateText` maps GitHub conclusion strings to colored symbols\n- Header now shows `CI: ✓ build · ✓ lint · ✗ test` style output\n",
	State:        "MERGED",
	BaseRefName:  "master",
	HeadRefName:  "feat/ci-checks",
	Additions:    87,
	Deletions:    12,
	ChangedFiles: 3,
	Author: api.Actor{
		Login: "dbd",
		Type:  "User",
		User:  api.User{Login: "dbd", Name: "Derek Daniels"},
	},
	CreatedAt: "2026-03-28T14:22:00Z",
	Repository: api.Repository{
		Name:               "ght",
		NameWithOwner:      "dbd/ght",
		MergeCommitAllowed: true,
		SquashMergeAllowed: true,
		RebaseMergeAllowed: true,
	},
	Labels: api.Labels{
		Nodes: []api.Label{
			{Name: "enhancement", Color: "a2eeef"},
		},
	},
	Reviews: api.Reviews{
		Nodes: []api.Review{
			{
				Author:      api.Actor{Login: "alice", Type: "User", User: api.User{Login: "alice", Name: "Alice Smith"}},
				State: "APPROVED",
			},
		},
	},
	HeadRef: api.HeadRef{
		Target: api.HeadRefTarget{
			Type: "Commit",
			Commit: api.HeadCommit{
				StatusCheckRollup: api.StatusCheckRollup{
					State: "SUCCESS",
					Contexts: api.StatusCheckRollupContextConnection{
						Nodes: []api.StatusCheckRollupContext{
							{Type: "CheckRun", CheckRun: api.CheckRun{Name: "build", Status: "COMPLETED", Conclusion: "SUCCESS"}},
							{Type: "CheckRun", CheckRun: api.CheckRun{Name: "lint", Status: "COMPLETED", Conclusion: "SUCCESS"}},
							{Type: "CheckRun", CheckRun: api.CheckRun{Name: "test", Status: "COMPLETED", Conclusion: "SUCCESS"}},
						},
					},
				},
			},
		},
	},
	TimelineItems: api.TimelineItems{
		Nodes: []api.PullRequestTimelineItem{
			{
				Type: "PullRequestCommit",
				PullRequestCommit: api.PullRequestCommit{
					Commit: api.Commit{
						AbbreviatedOid: "4805e58",
						Message:        "add CI status check rendering to header",
						AuthoredDate:   "2026-03-28T14:20:00Z",
						Author:         api.CommitAuthor{Name: "Derek Daniels", User: api.User{Login: "dbd", Name: "Derek Daniels"}},
					},
				},
			},
			{
				Type: "PullRequestCommit",
				PullRequestCommit: api.PullRequestCommit{
					Commit: api.Commit{
						AbbreviatedOid: "89f54dd",
						Message:        "fix CI state color for SKIPPED checks",
						AuthoredDate:   "2026-03-28T15:05:00Z",
						Author:         api.CommitAuthor{Name: "Derek Daniels", User: api.User{Login: "dbd", Name: "Derek Daniels"}},
					},
				},
			},
			{
				Type: "ReviewRequestedEvent",
				ReviewRequestedEvent: api.ReviewRequestedEvent{
					Actor:     api.Actor{Login: "dbd", Type: "User", User: api.User{Login: "dbd", Name: "Derek Daniels"}},
					CreatedAt: "2026-03-28T14:22:00Z",
					RequestedReviewer: api.RequestedReviewer{
						Type: "User",
						User: api.User{Login: "alice", Name: "Alice Smith"},
					},
				},
			},
			{
				Type: "PullRequestReview",
				PullRequestReview: api.PullRequestReview{
					Author:      api.Actor{Login: "alice", Type: "User", User: api.User{Login: "alice", Name: "Alice Smith"}},
					State:       "APPROVED",
					Body:        "Looks good! The CI formatting is clean and the color coding makes the state obvious at a glance.",
					SubmittedAt: "2026-03-29T09:10:00Z",
				},
			},
			{
				Type: "MergedEvent",
				MergedEvent: api.MergedEvent{
					Actor:     api.Actor{Login: "dbd", Type: "User", User: api.User{Login: "dbd", Name: "Derek Daniels"}},
					CreatedAt: "2026-03-29T10:00:00Z",
				},
			},
		},
	},
}

func main() {
	const detailWidth = 120
	termWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || termWidth <= 0 {
		termWidth = detailWidth
	}
	rendered := pullRequestDetail.RenderPullRequestDetail(staticPR, detailWidth)
	fmt.Print(components.CenterInViewport(rendered, termWidth, detailWidth))
}
