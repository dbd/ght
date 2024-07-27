package api

import (
	"fmt"
	"strconv"
	"time"
)

type PullRequests struct {
	Query        string
	PullRequests []PullRequestResponse
	Error        error
}

type PullRequestResponse struct {
	Assignees          Assignees `graphql:"assignees(first: 3)"`
	Author             Actor
	Title              string
	Body               string
	CreatedAt          Timestamp
	Repository         Repository
	Number             int64
	Additions          int64
	Deletions          int64
	ChangedFiles       int64
	State              string
	Mergeable          string
	BaseRefName        string
	HeadRefName        string
	Comments           IssueComments      `graphql:"comments(last: 20, orderBy: { field: UPDATED_AT, direction: DESC })"`
	PullRequestCommits PullRequestCommits `graphql:"commits(first: 10)"`
	ReviewThreads      ReviewThreads      `graphql:"reviewThreads(first: 10)"`
	Labels             Labels             `graphql:"labels(first: 10)"`
	Reviews            Reviews            `graphql:"reviews(first: 10)"`
	ReviewRequests     ReviewRequests     `graphql:"reviewRequests(first: 10)"`
	TimelineItems      TimelineItems      `graphql:"timelineItems(first: 30)"`
}

type Repository struct {
	Name          string
	NameWithOwner string
}

type Reviews struct {
	Nodes []Review
}

type Review struct {
	Author    Actor
	Body      string
	CreatedAt Timestamp
	State     string
}

type Labels struct {
	Nodes []Label
}

type Label struct {
	Color       string
	Name        string
	Description string
}
type Project struct {
	ID string
}

// People and Bots

type CommitAuthor struct {
	Name string
	User User
}

type Actor struct {
	Login        string
	Type         string       `graphql:"__typename"`
	User         User         `graphql:"... on User"`
	Bot          Bot          `graphql:"... on Bot"`
	Mannequin    Mannequin    `graphql:"... on Mannequin"`
	Organization Organization `graphql:"... on Organization"`
}

type Assignee struct {
	Type         string       `graphql:"__typename"`
	User         User         `graphql:"... on User"`
	Bot          Bot          `graphql:"... on Bot"`
	Mannequin    Mannequin    `graphql:"... on Mannequin"`
	Organization Organization `graphql:"... on Organization"`
}

type Assignees struct {
	Nodes []User
}

type User struct {
	Login string
	Name  string
}

type Organization struct {
	Login string
	Name  string
}

type Bot struct {
	Login string
}

type Mannequin struct {
	Login    string
	Claimant User
}

type Team struct {
	Name string
}

type ReviewRequests struct {
	Nodes []ReviewRequest
}
type ReviewRequest struct {
	RequestedReviewer RequestedReviewer
}

type RequestedReviewer struct {
	Type      string    `graphql:"__typename"`
	User      User      `graphql:"... on User"`
	Bot       Bot       `graphql:"... on Bot"`
	Mannequin Mannequin `graphql:"... on Mannequin"`
	Team      Team      `graphql:"... on Team"`
}

// Events

type PullRequestTimelineItem struct {
	Type                      string                    `graphql:"__typename"`
	AddedToProjectEvent       AddedToProjectEvent       `graphql:"... on AddedToProjectEvent"`
	AssignedEvent             AssignedEvent             `graphql:"... on AssignedEvent"`
	ClosedEvent               ClosedEvent               `graphql:"... on ClosedEvent"`
	IssueComment              IssueComment              `graphql:"... on IssueComment"`
	LabeledEvent              LabeledEvent              `graphql:"... on LabeledEvent"`
	MergedEvent               MergedEvent               `graphql:"... on MergedEvent"`
	PullRequestCommit         PullRequestCommit         `graphql:"... on PullRequestCommit"`
	PullRequestReview         PullRequestReview         `graphql:"... on PullRequestReview"`
	RemovedFromProjectEvent   RemovedFromProjectEvent   `graphql:"... on RemovedFromProjectEvent"`
	RenamedTitleEvent         RenamedTitleEvent         `graphql:"... on RenamedTitleEvent"`
	ReopenedEvent             ReopenedEvent             `graphql:"... on ReopenedEvent"`
	ReviewRequestRemovedEvent ReviewRequestRemovedEvent `graphql:"... on ReviewRequestRemovedEvent"`
	ReviewRequestedEvent      ReviewRequestedEvent      `graphql:"... on ReviewRequestedEvent"`
	UnassignedEvent           AssignedEvent             `graphql:"... on UnassignedEvent"`
	UnlabeledEvent            LabeledEvent              `graphql:"... on UnlabeledEvent"`
}

type AddedToProjectEvent struct {
	Actor     Actor
	CreatedAt Timestamp
	Project   Project
}

type AssignedEvent struct {
	Actor     Actor
	Assignee  Assignee
	CreatedAt Timestamp
}

type ClosedEvent struct {
	Actor     Actor
	CreatedAt Timestamp
}

type LabeledEvent struct {
	Actor     Actor
	CreatedAt Timestamp
	Label     Label
}

type MergedEvent struct {
	Actor     Actor
	Commit    Commit
	CreatedAt Timestamp
}

type RemovedFromProjectEvent struct {
	Actor     Actor
	CreatedAt Timestamp
	Project   Project
}

type RenamedTitleEvent struct {
	Actor         Actor
	CreatedAt     Timestamp
	CurrentTitle  string
	PreviousTitle string
}

type ReviewRequestedEvent struct {
	Actor             Actor
	CreatedAt         Timestamp
	RequestedReviewer RequestedReviewer
}

type ReviewRequestRemovedEvent struct {
	Actor             Actor
	CreatedAt         Timestamp
	RequestedReviewer RequestedReviewer
}

type ReviewDimissedEvent struct {
	Actor             Actor
	CreatedAt         Timestamp
	RequestedReviewer RequestedReviewer
}

type ReopenedEvent struct {
	Actor     Actor
	CreatedAt Timestamp
}

type IssueComments struct {
	Nodes []IssueComment
}

type TimelineItems struct {
	Nodes []PullRequestTimelineItem
}

type IssueComment struct {
	Author    Actor
	Body      string
	CreatedAt Timestamp
	UpdatedAt string
}
type PullRequestCommits struct {
	Nodes []PullRequestCommit
}

type PullRequestCommit struct {
	Commit Commit
}

type PullRequestReview struct {
	Author      Actor
	Body        string
	Commit      Commit
	State       string
	SubmittedAt Timestamp
	Comments    PullRequestReviewEventComments `graphql:"comments(first: 10)"`
}

type Commit struct {
	AbbreviatedOid string
	Author         CommitAuthor
	Additions      int64
	Deletions      int64
	Message        string
	AuthoredDate   Timestamp
}

type ReviewThreads struct {
	Nodes []ReviewThread
}

type ReviewThread struct {
	DiffSide   string
	IsResolved bool
	Path       string
	StartLine  int64
	Line       int64
	Comments   PullRequestReviewComments `graphql:"comments(first: 10)"`
}

type PullRequestReviewComments struct {
	Nodes []Comment
}

type PullRequestReviewEventComments struct {
	Nodes []PullRequestReviewComment
}

type PullRequestReviewComment struct {
	Author            Actor
	Body              string
	DiffHunk          string
	Outdated          bool
	Line              int64
	OriginalLine      int64
	StartLine         int64
	OriginalStartLine int64
}

type Comment struct {
	Id       string
	Author   Actor
	Body     string
	Outdated bool
}

type LineRange struct {
	StartLine int64
	EndLine   int64
}

type Timestamp string

func (t Timestamp) ShortSince() (s string) {
	layout := "2006-01-02T15:04:05Z"
	ts, _ := time.Parse(layout, string(t))
	now := time.Now()
	d := now.Sub(ts)
	if d.Hours() > 48 {
		s = fmt.Sprintf("%s days", strconv.FormatFloat(d.Hours()/24, 'f', 0, 64))
	} else if d.Hours() > 24 {
		s = fmt.Sprintf("%s day", strconv.FormatFloat(d.Hours(), 'f', 0, 64))
	} else if d.Minutes() == 1 {
		s = fmt.Sprintf("%s minute", strconv.FormatFloat(d.Minutes(), 'f', 0, 64))
	} else if d.Minutes() < 60 {
		s = fmt.Sprintf("%s minutes", strconv.FormatFloat(d.Minutes(), 'f', 0, 64))
	} else if d.Hours() == 1 {
		s = fmt.Sprintf("%s hour", strconv.FormatFloat(d.Hours(), 'f', 0, 64))
	} else {
		s = fmt.Sprintf("%s hours", strconv.FormatFloat(d.Hours(), 'f', 0, 64))

	}
	return s
}
