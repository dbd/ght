package issueDetail

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/dbd/ght/components"
	"github.com/dbd/ght/internal/api"
)

func RenderIssueDetail(issue api.IssueResponse, width int) string {
	doc := strings.Builder{}
	doc.WriteString(renderIssueHeader(issue))
	if issue.Milestone != nil {
		doc.WriteString(renderMilestoneBanner(*issue.Milestone) + "\n")
	}
	body, err := renderMarkdown(issue.Body)
	if err != nil {
		body = issue.Body
	}
	doc.WriteString(components.RenderBoxWithTitleCorner(issue.Author.Login, body, width, false, true) + "\n")
	doc.WriteString(renderIssueTimeline(issue, width))
	return doc.String()
}

func renderIssueHeader(issue api.IssueResponse) string {
	doc := strings.Builder{}
	doc.WriteString(components.PrTitleStyle.Render(issue.Title) + "\n")

	var stateColor lipgloss.Color
	if issue.State == "OPEN" {
		stateColor = components.Green
	} else {
		stateColor = components.Grey
	}
	state := lipgloss.NewStyle().Foreground(stateColor).Bold(true).Render(issue.State)
	doc.WriteString(fmt.Sprintf("#%d · %s · %s · %s\n", issue.Number, issue.Repository.NameWithOwner, issue.Author.Login, state))

	if len(issue.Labels.Nodes) > 0 {
		var labelParts []string
		for _, l := range issue.Labels.Nodes {
			labelParts = append(labelParts, lipgloss.NewStyle().Background(lipgloss.Color("#"+l.Color)).Render(l.Name))
		}
		doc.WriteString("Labels: " + strings.Join(labelParts, " ") + "\n")
	}

	if len(issue.Assignees.Nodes) > 0 {
		var names []string
		for _, a := range issue.Assignees.Nodes {
			names = append(names, components.BoldStyle.Render(a.Login))
		}
		doc.WriteString("Assignees: " + strings.Join(names, ", ") + "\n")
	}

	return doc.String()
}

func renderMilestoneBanner(m api.IssueMilestone) string {
	var stateColor lipgloss.Color
	if m.State == "OPEN" {
		stateColor = components.Green
	} else {
		stateColor = components.Grey
	}
	state := lipgloss.NewStyle().Foreground(stateColor).Render(m.State)

	var due string
	if string(m.DueOn) != "" {
		due = " · due " + m.DueOn.ShortSince() + " ago"
	}

	hint := components.RenderColoredText("[M] open milestone", "blue")
	content := fmt.Sprintf("Milestone: %s  %s%s  %s", components.BoldStyle.Render(m.Title), state, due, hint)
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(components.Blue).
		Padding(0, 1).
		Render(content)
}

func renderIssueTimeline(issue api.IssueResponse, width int) string {
	doc := strings.Builder{}
	doc.WriteString(renderIssueTimelineItems(issue.TimelineItems.Nodes, width))
	return doc.String()
}

func renderIssueTimelineItems(tis []api.IssueTimelineItem, width int) string {
	doc := strings.Builder{}
	type rr struct {
		body   string
		prefix bool
		ti     api.IssueTimelineItem
	}
	var ss []rr
	renderLater := []string{"IssueComment"}
	for _, ti := range tis {
		prefix := true
		s := formatIssueTimelineItem(ti, width, true, true)
		if s != "" {
			if slices.Contains(renderLater, ti.Type) {
				prefix = false
				s = ""
			}
			ss = append(ss, rr{body: s, prefix: prefix, ti: ti})
		}
	}
	for i, s := range ss {
		last := i == len(ss)-1
		body := s.body
		doc.WriteString(components.LineStyle.Render("│"))
		doc.WriteString("\n")
		if s.prefix {
			var c string
			if last {
				c = "╰ "
			} else {
				c = "├ "
			}
			doc.WriteString(components.LineStyle.Render(c))
		}
		if body == "" {
			body = formatIssueTimelineItem(s.ti, width, true, !last)
		}
		doc.WriteString(body)
		doc.WriteString("\n")
	}
	return doc.String()
}

func formatIssueTimelineItem(ti api.IssueTimelineItem, width int, topCorner, bottomCorner bool) string {
	switch ti.Type {
	case "IssueComment":
		t := fmt.Sprintf("%s commented %s ago", ti.IssueComment.Author.Login, ti.IssueComment.CreatedAt.ShortSince())
		return components.RenderBoxWithTitleCorner(t, ti.IssueComment.Body, width, topCorner, bottomCorner)
	case "AssignedEvent":
		return fmt.Sprintf("%s assigned %s %s ago", formatActor(ti.AssignedEvent.Actor), formatAssignee(ti.AssignedEvent.Assignee), ti.AssignedEvent.CreatedAt.ShortSince())
	case "UnassignedEvent":
		return fmt.Sprintf("%s removed assignee %s %s ago", formatActor(ti.UnassignedEvent.Actor), formatAssignee(ti.UnassignedEvent.Assignee), ti.UnassignedEvent.CreatedAt.ShortSince())
	case "LabeledEvent":
		label := lipgloss.NewStyle().Background(lipgloss.Color("#" + ti.LabeledEvent.Label.Color)).Render(ti.LabeledEvent.Label.Name)
		return fmt.Sprintf("%s added %s label %s ago", formatActor(ti.LabeledEvent.Actor), label, ti.LabeledEvent.CreatedAt.ShortSince())
	case "UnlabeledEvent":
		label := lipgloss.NewStyle().Background(lipgloss.Color("#" + ti.UnlabeledEvent.Label.Color)).Render(ti.UnlabeledEvent.Label.Name)
		return fmt.Sprintf("%s removed %s label %s ago", formatActor(ti.UnlabeledEvent.Actor), label, ti.UnlabeledEvent.CreatedAt.ShortSince())
	case "ClosedEvent":
		return formatClosedEvent(ti.ClosedEvent)
	case "CrossReferencedEvent":
		return formatCrossReferencedEvent(ti.CrossReferencedEvent)
	case "ReopenedEvent":
		return fmt.Sprintf("%s reopened the issue %s ago", formatActor(ti.ReopenedEvent.Actor), ti.ReopenedEvent.CreatedAt.ShortSince())
	case "RenamedTitleEvent":
		return fmt.Sprintf("%s changed title from '%s' to '%s' %s ago", formatActor(ti.RenamedTitleEvent.Actor), ti.RenamedTitleEvent.PreviousTitle, ti.RenamedTitleEvent.CurrentTitle, ti.RenamedTitleEvent.CreatedAt.ShortSince())
	case "MilestonedEvent":
		return fmt.Sprintf("%s added to milestone '%s' %s ago", formatActor(ti.MilestonedEvent.Actor), ti.MilestonedEvent.MilestoneTitle, ti.MilestonedEvent.CreatedAt.ShortSince())
	case "DemilestonedEvent":
		return fmt.Sprintf("%s removed from milestone '%s' %s ago", formatActor(ti.DemilestonedEvent.Actor), ti.DemilestonedEvent.MilestoneTitle, ti.DemilestonedEvent.CreatedAt.ShortSince())
	}
	return ""
}

func formatClosedEvent(e api.ClosedEvent) string {
	actor := formatActor(e.Actor)
	ago := e.CreatedAt.ShortSince()

	var reason string
	switch e.StateReason {
	case "COMPLETED":
		reason = components.RenderColoredText("completed", "green")
	case "NOT_PLANNED":
		reason = components.RenderColoredText("not planned", "grey")
	default:
		reason = "closed"
	}

	switch e.Closer.Type {
	case "PullRequest":
		pr := e.Closer.PullRequest
		prRef := components.BoldStyle.Render(fmt.Sprintf("#%d %s", pr.Number, pr.Title))
		return fmt.Sprintf("%s closed as %s via %s %s ago", actor, reason, prRef, ago)
	case "Commit":
		commit := e.Closer.Commit
		sha := components.BoldStyle.Render(commit.AbbreviatedOid)
		return fmt.Sprintf("%s closed as %s via commit %s %s ago", actor, reason, sha, ago)
	}
	return fmt.Sprintf("%s closed as %s %s ago", actor, reason, ago)
}

func formatCrossReferencedEvent(e api.CrossReferencedEvent) string {
	ago := e.CreatedAt.ShortSince()
	switch e.Source.Type {
	case "PullRequest":
		pr := e.Source.PullRequest
		prRef := components.BoldStyle.Render(fmt.Sprintf("%s#%d", pr.Repository.NameWithOwner, pr.Number))
		var stateColor string
		switch pr.State {
		case "OPEN":
			stateColor = "green"
		case "MERGED":
			stateColor = "purple"
		default:
			stateColor = "grey"
		}
		state := components.RenderColoredText(pr.State, stateColor)
		if e.WillCloseTarget {
			return fmt.Sprintf("mentioned in PR %s [%s] — will close this issue %s ago", prRef, state, ago)
		}
		return fmt.Sprintf("mentioned in PR %s [%s] %s ago", prRef, state, ago)
	case "Issue":
		issue := e.Source.Issue
		issueRef := components.BoldStyle.Render(fmt.Sprintf("%s#%d", issue.Repository.NameWithOwner, issue.Number))
		return fmt.Sprintf("mentioned in issue %s %s ago", issueRef, ago)
	}
	return ""
}

func formatActor(a api.Actor) string {
	switch a.Type {
	case "User":
		return components.BoldStyle.Render(fmt.Sprintf("%s (%s)", a.User.Login, a.User.Name))
	case "Bot":
		return components.BoldStyle.Render(a.Bot.Login)
	case "Organization":
		return components.BoldStyle.Render(a.Organization.Name)
	case "Mannequin":
		return components.BoldStyle.Render(fmt.Sprintf("%s (%s)", a.Mannequin.Login, a.Mannequin.Claimant.Name))
	}
	return a.Login
}

func formatAssignee(a api.Assignee) string {
	switch a.Type {
	case "User":
		return components.BoldStyle.Render(fmt.Sprintf("%s (%s)", a.User.Login, a.User.Name))
	}
	return ""
}

func renderMarkdown(body string) (string, error) {
	return glamour.Render(body, "dark")
}
