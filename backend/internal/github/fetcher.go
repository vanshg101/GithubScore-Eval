package github

import (
	"log"
	"sync"
	"time"

	"github.com/Madhur/GithubScoreEval/backend/internal/model"
)

type RepoData struct {
	Repo     Repository
	Commits  []Commit
	PRs      []PullRequest
	Issues   []Issue
	Reviews  []ReviewComment
}

func (c *Client) FetchDeveloperData(username string) (*model.Developer, error) {
	profile, err := c.FetchUserProfile(username)
	if err != nil {
		return nil, err
	}

	repos, err := c.FetchRepos(username)
	if err != nil {
		return nil, err
	}

	sem := make(chan struct{}, maxConcurrent)
	var mu sync.Mutex
	var wg sync.WaitGroup
	var repoDataList []RepoData

	for _, repo := range repos {
		wg.Add(1)
		go func(r Repository) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			rd := RepoData{Repo: r}

			commits, err := c.FetchCommits(r.Owner.Login, r.Name, username)
			if err != nil {
				log.Printf("error fetching commits for %s/%s: %v", r.Owner.Login, r.Name, err)
			} else {
				rd.Commits = commits
			}

			prs, err := c.FetchPullRequests(r.Owner.Login, r.Name)
			if err != nil {
				log.Printf("error fetching PRs for %s/%s: %v", r.Owner.Login, r.Name, err)
			} else {
				rd.PRs = prs
			}

			issues, err := c.FetchIssues(r.Owner.Login, r.Name)
			if err != nil {
				log.Printf("error fetching issues for %s/%s: %v", r.Owner.Login, r.Name, err)
			} else {
				rd.Issues = issues
			}

			reviews, err := c.FetchReviewComments(r.Owner.Login, r.Name)
			if err != nil {
				log.Printf("error fetching reviews for %s/%s: %v", r.Owner.Login, r.Name, err)
			} else {
				rd.Reviews = reviews
			}

			mu.Lock()
			repoDataList = append(repoDataList, rd)
			mu.Unlock()
		}(repo)
	}

	wg.Wait()

	developer := aggregateData(username, profile, repos, repoDataList)
	return developer, nil
}

func aggregateData(username string, profile *UserProfile, repos []Repository, repoData []RepoData) *model.Developer {
	var totalCommits, totalPRs, mergedPRs int
	var totalIssuesOpened, totalIssuesClosed int
	var totalReviewComments int
	var totalStars, totalForks int
	var totalPRLines int
	var prCount int
	var totalResponseHours float64
	var issueResponseCount int

	commitWeeks := make(map[string]bool)
	languageSet := make(map[string]bool)
	reposContributed := 0

	for _, r := range repos {
		totalStars += r.StargazersCount
		totalForks += r.ForksCount
		if r.Language != "" {
			languageSet[r.Language] = true
		}
	}

	for _, rd := range repoData {
		hasContribution := false

		for _, commit := range rd.Commits {
			totalCommits++
			hasContribution = true
			week := commit.Commit.Author.Date.Format("2006-W01")
			commitWeeks[week] = true
		}

		for _, pr := range rd.PRs {
			if pr.User.Login == username {
				totalPRs++
				hasContribution = true
				if pr.MergedAt != nil {
					mergedPRs++
				}
				linesChanged := pr.Additions + pr.Deletions
				totalPRLines += linesChanged
				prCount++
			}
		}

		for _, issue := range rd.Issues {
			if issue.User.Login == username {
				totalIssuesOpened++
				hasContribution = true
				if issue.State == "closed" && issue.ClosedAt != nil {
					totalIssuesClosed++
					responseTime := issue.ClosedAt.Sub(issue.CreatedAt).Hours()
					totalResponseHours += responseTime
					issueResponseCount++
				}
			}
		}

		for _, review := range rd.Reviews {
			if review.User.Login == username {
				totalReviewComments++
				hasContribution = true
			}
		}

		if hasContribution {
			reposContributed++
		}
	}

	var avgPRLines float64
	if prCount > 0 {
		avgPRLines = float64(totalPRLines) / float64(prCount)
	}

	var avgResponseHours float64
	if issueResponseCount > 0 {
		avgResponseHours = totalResponseHours / float64(issueResponseCount)
	}

	commitTrend := "stable"

	languages := make([]string, 0, len(languageSet))
	for lang := range languageSet {
		languages = append(languages, lang)
	}

	return &model.Developer{
		Username: username,
		Profile: model.DeveloperProfile{
			Name:        profile.Name,
			Bio:         profile.Bio,
			PublicRepos: profile.PublicRepos,
			Followers:   profile.Followers,
			AvatarURL:   profile.AvatarURL,
		},
		Metrics: model.DeveloperMetrics{
			TotalCommits:          totalCommits,
			TotalPRs:              totalPRs,
			MergedPRs:             mergedPRs,
			TotalIssuesOpened:     totalIssuesOpened,
			TotalIssuesClosed:     totalIssuesClosed,
			ReviewComments:        totalReviewComments,
			ActiveWeeks:           len(commitWeeks),
			ReposContributed:      reposContributed,
			TotalStars:            totalStars,
			TotalForks:            totalForks,
			AvgPRLinesChanged:     avgPRLines,
			AvgIssueResponseHours: avgResponseHours,
			CommitTrend:           commitTrend,
			Languages:             languages,
		},
		FetchedAt: time.Now(),
	}
}
