package github_client

import (
	"strings"
	"time"

	"github.com/google/go-github/github"
)

func hasCommitInLast12Months(pushedAt *github.Timestamp) bool {
	now := time.Now()
	lastTouchDate := pushedAt.Time
	diff := now.Sub(lastTouchDate)
	return diff.Seconds()/60/60/24 < 365
}

func includesInCsv(csv string, repoName string) bool {
	arr := strings.Split(csv, ",")
	for _, repo := range arr {
		if repo == repoName {
			return true
		}
	}
	return false
}
