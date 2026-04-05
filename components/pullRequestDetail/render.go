package pullRequestDetail

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dbd/ght/components"
	"github.com/dbd/ght/internal/api"
	"github.com/dbd/ght/utils"
)

func renderTimeline(pr api.PullRequestResponse, width int) string {
	doc := strings.Builder{}
	doc.WriteString(renderTimelineItems(pr.TimelineItems.Nodes, width))
	return doc.String()
}

func renderTimelineItems(tis []api.PullRequestTimelineItem, width int) string {
	doc := strings.Builder{}
	type rr struct {
		body   string
		prefix bool
		ti     api.PullRequestTimelineItem
	}
	var s string
	var c string
	var ss []rr
	last := false
	renderLater := []string{"IssueComment", "PullRequestReview"}
	for _, ti := range tis {
		prefix := true
		s = formatTimelineItem(ti, width, true, true)
		if s != "" {
			if slices.Contains(renderLater, ti.Type) {
				prefix = false
				s = ""
			}
			ss = append(ss, rr{body: s, prefix: prefix, ti: ti})
		}
	}
	for i, s := range ss {
		body := s.body
		if i == len(ss)-1 {
			last = true
		}
		doc.WriteString(components.LineStyle.Render("│"))
		doc.WriteString("\n")
		if s.prefix {
			if last {
				c = "╰ "
			} else {
				c = "├ "
			}
			doc.WriteString(components.LineStyle.Render(c))
		}
		if body == "" {
			body = formatTimelineItem(s.ti, width, true, !last)
		}
		doc.WriteString(body)
		doc.WriteString("\n")
	}
	return doc.String()
}

func formatTimelineItem(p api.PullRequestTimelineItem, width int, topCorner, bottomCorner bool) (s string) {
	switch p.Type {
	case "AddedToProjectEvent":
		s = fmt.Sprintf("%s added to project %s %s ago", formatActor(p.AddedToProjectEvent.Actor), p.AddedToProjectEvent.Project.ID, p.AddedToProjectEvent.CreatedAt.ShortSince())
	case "AssignedEvent":
		s = fmt.Sprintf("%s assigned %s %s ago", formatActor(p.AssignedEvent.Actor), formatAssignee(p.AssignedEvent.Assignee), p.AssignedEvent.CreatedAt.ShortSince())
	case "ClosedEvent":
		s = fmt.Sprintf("%s closed the PullRequest %s ago", formatActor(p.ClosedEvent.Actor), p.ClosedEvent.CreatedAt.ShortSince())
	case "IssueComment":
		t := fmt.Sprintf("%s commented %s ago", p.IssueComment.Author.Login, p.IssueComment.CreatedAt.ShortSince())
		s = components.RenderBoxWithTitleCorner(t, p.IssueComment.Body, width, topCorner, bottomCorner)
	case "PullRequestReview":
		t := formatReviewState(p.PullRequestReview)
		var b string
		if len(p.PullRequestReview.Comments.Nodes) >= 1 {
			b = formatReviewComment(p.PullRequestReview.Comments.Nodes[0], width/2)
		} else {
			b = p.IssueComment.Body
		}
		s = components.RenderBoxWithTitleCorner(t, b, width, topCorner, bottomCorner)
	case "LabeledEvent":
		label := lipgloss.NewStyle().Background(lipgloss.Color("#" + p.LabeledEvent.Label.Color)).Render(p.LabeledEvent.Label.Name)
		s = fmt.Sprintf("%s added %s label %s ago", formatActor(p.LabeledEvent.Actor), label, p.LabeledEvent.CreatedAt.ShortSince())
	case "MergedEvent":
		s = fmt.Sprintf("%s merged %s ago", formatActor(p.MergedEvent.Actor), p.MergedEvent.CreatedAt.ShortSince())
	case "RemovedFromProjectEvent":
		s = fmt.Sprintf("%s removed from project %s %s ago", formatActor(p.RemovedFromProjectEvent.Actor), p.RemovedFromProjectEvent.Project.ID, p.RemovedFromProjectEvent.CreatedAt.ShortSince())
	case "RenamedTitleEvent":
		s = fmt.Sprintf("%s change tile from '%s' to '%s' %s ago", formatActor(p.RenamedTitleEvent.Actor), p.RenamedTitleEvent.PreviousTitle, p.RenamedTitleEvent.CurrentTitle, p.RenamedTitleEvent.CreatedAt.ShortSince())
	case "ReopenedEvent":
		s = fmt.Sprintf("%s reopened %s ago", formatActor(p.ReopenedEvent.Actor), p.MergedEvent.CreatedAt.ShortSince())
	case "ReviewRequestedEvent":
		s = fmt.Sprintf("%s requested review from %s %s ago", formatActor(p.ReviewRequestedEvent.Actor), formatReviewer(p.ReviewRequestedEvent.RequestedReviewer), p.ReviewRequestedEvent.CreatedAt.ShortSince())
	case "ReviewRequestRemovedEvent":
		s = fmt.Sprintf("%s removed requested review from %s %s ago", formatActor(p.ReviewRequestRemovedEvent.Actor), formatReviewer(p.ReviewRequestRemovedEvent.RequestedReviewer), p.ReviewRequestRemovedEvent.CreatedAt.ShortSince())
	case "PullRequestCommit":
		s = fmt.Sprintf("%s committed \"%s\" (%s) %s ago", formatCommitAuthor(p.PullRequestCommit.Commit.Author), p.PullRequestCommit.Commit.Message, p.PullRequestCommit.Commit.AbbreviatedOid, p.PullRequestCommit.Commit.AuthoredDate.ShortSince())
	case "UnassignedEvent":
		s = fmt.Sprintf("%s removed assignee %s %s ago", formatActor(p.AssignedEvent.Actor), formatAssignee(p.AssignedEvent.Assignee), p.AssignedEvent.CreatedAt.ShortSince())
	case "UnlabeledEvent":
		label := lipgloss.NewStyle().Background(lipgloss.Color("#" + p.LabeledEvent.Label.Color)).Render(p.LabeledEvent.Label.Name)
		s = fmt.Sprintf("%s removed %s label %s ago", formatActor(p.LabeledEvent.Actor), label, p.LabeledEvent.CreatedAt.ShortSince())
	}
	return s
}

func formatActor(a api.Actor) (s string) {
	switch a.Type {
	case "User":
		s = components.BoldStyle.Render(fmt.Sprintf("%s (%s)", a.User.Login, a.User.Name))
	case "Mannequin":
		s = components.BoldStyle.Render(fmt.Sprintf("%s (%s)", a.Mannequin.Login, a.Mannequin.Claimant.Name))
	case "Organization":
		s = components.BoldStyle.Render(fmt.Sprintf("%s", a.Organization.Name))
	case "Bot":
		s = components.BoldStyle.Render(fmt.Sprintf("%s", a.Bot.Login))
	}
	return
}

func formatAssignee(a api.Assignee) (s string) {
	switch a.Type {
	case "User":
		s = components.BoldStyle.Render(fmt.Sprintf("%s (%s)", a.User.Login, a.User.Name))
	}
	return
}

func formatCommitAuthor(a api.CommitAuthor) (s string) {
	s = components.BoldStyle.Render(fmt.Sprintf("%s (%s)", a.User.Login, a.User.Name))
	return
}

func formatReviewState(a api.PullRequestReview) (s string) {
	switch a.State {
	case "COMMENTED":
		u := components.BoldStyle.Render(formatActor(a.Author))
		b := components.RenderColoredText("commented", "yellow")
		s = fmt.Sprintf("%s %s %s ago", u, b, a.SubmittedAt.ShortSince())
	case "PENDING":
		u := components.BoldStyle.Render(formatActor(a.Author))
		s = fmt.Sprintf("%s is pending review from %s ago", u, a.SubmittedAt.ShortSince())
	case "APPROVED":
		u := components.BoldStyle.Render(formatActor(a.Author))
		b := components.RenderColoredText("approved", "green")
		s = fmt.Sprintf("%s %s %s ago", u, b, a.SubmittedAt.ShortSince())
	case "CHANGES_REQUESTED":
		u := components.BoldStyle.Render(formatActor(a.Author))
		b := components.RenderColoredText("requested changes", "red")
		s = fmt.Sprintf("%s %s %s ago", u, b, a.SubmittedAt.ShortSince())
	case "DISMISSED":
		u := components.BoldStyle.Render(formatActor(a.Author))
		b := components.RenderColoredText("dimissed", "red")
		s = fmt.Sprintf("%s %s the review %s ago", u, b, a.SubmittedAt.ShortSince())
	}
	return
}

func formatReviewer(a api.RequestedReviewer) (s string) {
	switch a.Type {
	case "User":
		s = components.BoldStyle.Render(fmt.Sprintf("%s (%s)", a.User.Login, a.User.Name))
	case "Mannequin":
		s = components.BoldStyle.Render(fmt.Sprintf("%s (%s)", a.Mannequin.Login, a.Mannequin.Claimant.Name))
	case "Team":
		s = components.BoldStyle.Render(fmt.Sprintf("%s", a.Team.Name))
	case "Bot":
		s = components.BoldStyle.Render(fmt.Sprintf("%s", a.Bot.Login))
	}
	return
}

func formatReviewComment(c api.PullRequestReviewComment, width int) (s string) {
	doc := strings.Builder{}
	hunk := utils.ParseHunkDiff(c.DiffHunk)
	for _, line := range hunk.Lines {
		var lf string
		var side string
		var lr int64
		if line.Left {
			lf = components.DeletionsStyle.Render(line.Raw)
			side = "LEFT"
			lr = line.LeftNumber
		} else if line.Right {
			lf = components.AdditionsStyle.Render(line.Raw)
			side = "RIGHT"
			lr = line.RightNumber
		} else {
			lf = line.Raw
			lr = line.LeftNumber
		}
		_, _ = side, lr
		lp := components.DiffLineNumberStyle.Render(fmt.Sprintf("%d,%d", line.LeftNumber, line.RightNumber))
		doc.WriteString(fmt.Sprintf("%s %s\n", lp, lf))
		if c.OriginalLine == lr {
			b := fmt.Sprintf("%s: %s", formatActor(c.Author), c.Body)
			doc.WriteString(components.RenderBox(b, width))
			doc.WriteString("\n")
		}
	}
	return doc.String()
}

func formatHeader(pr api.PullRequestResponse) string {
	doc := strings.Builder{}
	ad := components.AdditionsStyle.Render("+"+strconv.FormatInt(pr.Additions, 10)) + " · " + components.DeletionsStyle.Render("-"+strconv.FormatInt(pr.Deletions, 10))
	doc.WriteString(components.PrTitleStyle.Render(pr.Title) + "\n")
	doc.WriteString(pr.Author.Login + " · " + pr.BaseRefName + " ← " + pr.HeadRefName + "\n")
	doc.WriteString(strconv.FormatInt(pr.Number, 10) + " · " + pr.Repository.NameWithOwner + " | " + ad + "\n")
	assignees := formatAssignees(pr)
	if assignees != "" {
		doc.WriteString("Assignees: " + assignees + "\n")
	}
	doc.WriteString("Reviewers: " + formatReviewers(pr) + "\n")
	doc.WriteString("CI: " + formatCIChecks(pr) + "\n")
	return doc.String()
}

func formatCIChecks(pr api.PullRequestResponse) string {
	state := pr.CIState()
	if state == "" {
		return "-"
	}
	checks := pr.CIChecks()
	if len(checks) == 0 {
		return formatCIStateText(state)
	}
	doc := strings.Builder{}
	doc.WriteString(formatCIStateText(state) + " ")
	for i, c := range checks {
		if i > 0 {
			doc.WriteString(" · ")
		}
		var name, conclusion string
		switch c.Type {
		case "CheckRun":
			name = c.CheckRun.Name
			conclusion = c.CheckRun.Conclusion
			if conclusion == "" {
				conclusion = c.CheckRun.Status
			}
		case "StatusContext":
			name = c.StatusContext.Context
			conclusion = c.StatusContext.State
		}
		doc.WriteString(formatCIStateText(conclusion) + " " + name)
	}
	return doc.String()
}

func formatCIStateText(state string) string {
	switch state {
	case "SUCCESS", "success":
		return components.RenderColoredText("✓", "green")
	case "FAILURE", "failure", "FAILED":
		return components.RenderColoredText("✗", "red")
	case "PENDING", "pending", "IN_PROGRESS", "QUEUED":
		return components.RenderColoredText("●", "yellow")
	case "ERROR", "error", "ACTION_REQUIRED", "TIMED_OUT", "CANCELLED", "STARTUP_FAILURE":
		return components.RenderColoredText("!", "red")
	case "SKIPPED", "NEUTRAL":
		return components.RenderColoredText("–", "grey")
	case "EXPECTED":
		return components.RenderColoredText("?", "grey")
	default:
		return components.RenderColoredText(state, "grey")
	}
}

func formatAssignees(pr api.PullRequestResponse) string {
	doc := strings.Builder{}
	for i, assignee := range pr.Assignees.Nodes {
		if i > 0 {
			doc.WriteString(", ")
		}
		doc.WriteString(components.BoldStyle.Render(assignee.Login))
	}
	return doc.String()
}

func formatReviewers(pr api.PullRequestResponse) string {
	doc := strings.Builder{}
	rm := map[string]string{}
	for _, reviewer := range pr.ReviewRequests.Nodes {
		rm[formatReviewer(reviewer.RequestedReviewer)] = ""
	}
	for _, r := range pr.Reviews.Nodes {
		c := ""
		switch r.State {
		case "COMMENTED":
			c = "yellow"
		case "PENDING":
			c = ""
		case "APPROVED":
			c = "green"
		case "CHANGES_REQUESTED":
			c = "red"
		case "DISMISSED":
			c = "grey"
		}
		rm[formatActor(r.Author)] = c
	}
	for k, v := range rm {
		doc.WriteString(components.RenderColoredText(k, v))
	}
	return doc.String()
}

func RenderMergeOverlay(pr api.PullRequestResponse, width int) string {
	var allowedMergeTypes []string
	doc := strings.Builder{}
	if pr.Repository.MergeCommitAllowed {
		allowedMergeTypes = append(allowedMergeTypes, "Merge Commit")
	}
	if pr.Repository.RebaseMergeAllowed {
		allowedMergeTypes = append(allowedMergeTypes, "Rebase")
	}
	if pr.Repository.SquashMergeAllowed {
		allowedMergeTypes = append(allowedMergeTypes, "Squash Commits")
	}

	doc.WriteString("Merging pull request " + pr.Repository.NameWithOwner + "#" + strconv.FormatInt(pr.Number, 10) + "\n")
	doc.WriteString("Allowed merge types\n")
	for _, t := range allowedMergeTypes {
		doc.WriteString("\t - " + t + "\n")
	}
	return doc.String()
}
