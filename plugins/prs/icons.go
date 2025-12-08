package prs

import (
	"xbar/pkg/clients"
)

const (
	IconPlugin           = "\uf408" // nf-cod-github_alt
	IconPR               = "\uf407" // nf-oct-git_pull_request
	IconComment          = "\uf442" // nf-oct-comment_discussion
	IconApproved         = "\uf42e" // nf-oct-check
	IconChangesRequested = "\uf467" // nf-oct-x
	IconReviewRequired   = "\uf441" // nf-oct-eye
	IconStale            = "\uf43a" // nf-oct-clock
	IconDraft            = "\uf4dd" // nf-oct-git_pull_request_draft
	IconAssigned         = "\uf415" // nf-oct-person
	IconWarning          = "⚠️"
)

func GetPRStatusIcon(pr *clients.PullRequestMeta) string {
	if pr.IsDraft {
		return IconDraft
	}
	switch pr.ReviewDecision {
	case "APPROVED":
		return IconApproved
	case "CHANGES_REQUESTED":
		return IconChangesRequested
	case "REVIEW_REQUIRED":
		return IconReviewRequired
	default:
		return IconPR
	}
}

func GetPRStateIcons(pr *clients.PullRequestMeta) []string {
	icons := []string{}
	icons = append(icons, GetPRStatusIcon(pr))
	if pr.HasUnresolvedComments {
		icons = append(icons, IconComment)
	}
	if pr.IsOld {
		icons = append(icons, IconStale)
	}
	return icons
}
