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
	CreatedAt          TimeStamp
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
	CreatedAt TimeStamp
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
type Assignees struct {
	Nodes []Assignee
}
type Assignee struct {
	Login string
	Name  string
}

type IssueComments struct {
	Nodes []IssueComment
}
type IssueComment struct {
	Author    Actor
	Body      string
	CreatedAt TimeStamp
	UpdatedAt string
}
type PullRequestCommits struct {
	Nodes []PullRequestCommit
}

type PullRequestCommit struct {
	Commit Commit
}

type Commit struct {
	AbbreviatedOid string
	Author         Author
	Additions      int64
	Deletions      int64
	Message        string
	AuthoredDate   string
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

type Author struct {
	Name string
}

type Actor struct {
	Login string
}

type PullRequestReviewComments struct {
	Nodes []Comment
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

func (c ReviewThreads) getPRCommentsMap() map[string]map[string]map[LineRange][]Comment {
	var res = map[string]map[string]map[LineRange][]Comment{}
	for _, rt := range c.Nodes {
		lr := LineRange{StartLine: rt.StartLine, EndLine: rt.Line}
		if _, ok := res[rt.Path]; !ok {
			res[rt.Path] = map[string]map[LineRange][]Comment{}
		}
		if _, ok := res[rt.Path][rt.DiffSide]; !ok {
			res[rt.Path][rt.DiffSide] = map[LineRange][]Comment{}
		}
		res[rt.Path][rt.DiffSide][lr] = rt.Comments.Nodes
	}
	return res
}

type TimeStamp string

func (t TimeStamp) ShortSince() (s string) {
	layout := "2006-01-02T15:04:05Z"
	ts, _ := time.Parse(layout, string(t))
	now := time.Now()
	d := now.Sub(ts)
	if d.Hours() > 23 {
		s = fmt.Sprintf("%s days", strconv.FormatFloat(d.Hours()/24, 'f', 0, 64))
	} else {
		s = fmt.Sprintf("%s hours", strconv.FormatFloat(d.Hours(), 'f', 0, 64))
	}
	return s
}
