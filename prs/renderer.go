package prs

import (
	"fmt"
	"xbar/clients"
	"xbar/utils"
	xbar "xbar/xbar_utils"
)

func RenderPRList(prList *clients.PRList, env *Environment) []*xbar.XBarLine {
	lines := []*xbar.XBarLine{}

	if len(prList.MyPRs) > 0 {
		lines = append(lines, RenderMyPRs(prList.MyPRs, env))
	}

	if len(prList.ReviewRequests) > 0 {
		lines = append(lines, RenderReviewRequests(prList.ReviewRequests, env))
	}

	if len(prList.Assigned) > 0 {
		lines = append(lines, RenderAssignedPRs(prList.Assigned, env))
	}
	return lines
}

func RenderMyPRs(prs []*clients.PullRequestMeta, env *Environment) *xbar.XBarLine {
	children := make([]*xbar.XBarLine, len(prs))

	for i, pr := range prs {
		children[i] = RenderPR(pr, env)
	}

	return xbar.NewXBarLine(
		fmt.Sprintf("%s My Open PRs (%d)", IconPR, len(prs)),
		xbar.WithChildren(children...),
	)
}

func RenderReviewRequests(prs []*clients.PullRequestMeta, env *Environment) *xbar.XBarLine {
	children := make([]*xbar.XBarLine, len(prs))

	for i, pr := range prs {
		children[i] = RenderPR(pr, env)
	}

	return xbar.NewXBarLine(
		fmt.Sprintf("%s Review Requests (%d)", IconReviewRequired, len(prs)),
		xbar.WithChildren(children...),
	)
}

func RenderAssignedPRs(prs []*clients.PullRequestMeta, env *Environment) *xbar.XBarLine {
	children := make([]*xbar.XBarLine, len(prs))

	for i, pr := range prs {
		children[i] = RenderPR(pr, env)
	}

	return xbar.NewXBarLine(
		fmt.Sprintf("%s Assigned to Me (%d)", IconAssigned, len(prs)),
		xbar.WithChildren(children...),
	)
}

func RenderPR(pr *clients.PullRequestMeta, env *Environment) *xbar.XBarLine {
	// Get status icons
	icons := GetPRStateIcons(pr)
	iconString := ""
	for _, icon := range icons {
		iconString += icon + " "
	}
	if env.AddIconSpace {
		iconString += " "
	}
	// icons owner/repo #123: Title
	text := fmt.Sprintf(
		"%s%s/%s #%d: %s",
		iconString,
		pr.RepoOwner,
		pr.RepoName,
		pr.Number,
		pr.Title,
	)

	// Add metadata as children (only age info if old)
	children := []*xbar.XBarLine{}

	ageColor := "primary"
	// Age info for old PRs
	if pr.IsOld {
		ageColor = "warning"
		if pr.DaysSinceUpdate >= 14 {
			ageColor = "danger"
		}

		children = append(children, xbar.NewXBarLine(
			fmt.Sprintf("%s %d days since last update", IconStale, pr.DaysSinceUpdate),
			xbar.WithColor(utils.Color(ageColor)),
		))
	}

	return xbar.NewXBarLine(
		text,
		xbar.WithHref(pr.URL),
		xbar.WithMaxLength(80),
		xbar.WithColor(ageColor),
		xbar.WithChildren(children...),
	)
}
