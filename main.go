package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/google/go-github/v39/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	// Set up authentication using a personal access token
	// For local create .env file, for GitHub actions see workflow yaml file
	token := os.Getenv("GH_TOKEN")
	if token == "" {
		log.Fatal("Please set the GH_TOKEN environment variable")
	}
	ctx := context.Background()
	client := github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})))

	// Check if we can get authenticated user
	_, _, err := client.Users.Get(ctx, "")
	if err != nil {
		log.Fatalf("Failed to get user: %v", err)
	}
	allRepos := getAllUserRepository(&ctx, client)
	langMap := getLanguageUsedWithProjectCount(ctx, client, allRepos)
	sortedLanguage := getSortedLanguageByUsage(langMap)

	for i, lang := range sortedLanguage[:5] {
		fmt.Printf("%d. %s: (%d)\n", i+1, lang, langMap[lang])
	}

}

func getLanguageUsedWithProjectCount(ctx context.Context, client *github.Client, allRepos []*github.Repository) map[string]int {
	languages := make(map[string]int)
	for _, repo := range allRepos {
		langs, _, err := client.Repositories.ListLanguages(ctx, *repo.Owner.Login, *repo.Name)
		if err != nil {
			log.Printf("Failed to list languages for repository %s: %v", *repo.FullName, err)
			continue
		}
		for lang := range langs {
			languages[lang]++
		}
	}
	return languages
}

func getSortedLanguageByUsage(langMap map[string]int) []string {
	var topLanguages []string
	for lang := range langMap {
		topLanguages = append(topLanguages, lang)
	}
	sort.Slice(topLanguages, func(i, j int) bool {
		return langMap[topLanguages[i]] > langMap[topLanguages[j]]
	})
	return topLanguages
}

func getAllUserRepository(ctx *context.Context, client *github.Client) []*github.Repository {
	var allRepos []*github.Repository

	// Get the user's repositories
	opt := &github.RepositoryListOptions{
		Affiliation: "owner,collaborator",
		ListOptions: github.ListOptions{PerPage: 100},
	}
	for {
		repos, resp, err := client.Repositories.List(*ctx, "", opt)
		if err != nil {
			log.Fatalf("Failed to list repositories: %v", err)
		}
		for _, repo := range repos {
			if hasCommitThisYear(repo.PushedAt) && !includesInCsv(os.Getenv("EXCLUDE_REPOS"), *repo.Name) {
				allRepos = append(allRepos, repo)
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	if os.Getenv("PRINT_PROJECT_NAME") == "TRUE" {
		for _, repo := range allRepos {
			log.Printf("Found repository: %s", *repo.FullName)
		}
	}

	return allRepos
}

func hasCommitThisYear(pushedAt *github.Timestamp) bool {
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
