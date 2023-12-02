package github_client

import (
	"context"
	"log"
	"os"
	"sort"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func New() *gitHubClient {
	gc := gitHubClient{}
	return gc.new()
}

type gitHubClient struct {
	Token      string
	AllRepos   []*github.Repository
	LangMap    map[string]int
	SortedLang []string
	ctx        context.Context
	client     *github.Client
}

func (gc *gitHubClient) new() *gitHubClient {
	token := os.Getenv("GH_TOKEN")
	if token == "" {
		log.Fatalln("Please set the GH_TOKEN environment variable")
	}
	gc.Token = token
	gc.ctx = context.Background()
	gc.client = github.NewClient(oauth2.NewClient(gc.ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: gc.Token})))
	return gc
}

func (gc *gitHubClient) FetchAllRepository() *gitHubClient {
	if gc.ctx == nil || gc.client == nil {
		log.Fatalln("gitHubClient New method not called")
	}
	var allRepos []*github.Repository
	opt := &github.RepositoryListOptions{
		Affiliation: "owner,collaborator",
		ListOptions: github.ListOptions{PerPage: 100},
	}
	for {
		repos, resp, err := gc.client.Repositories.List(gc.ctx, "", opt)
		if err != nil {
			log.Fatalf("Failed to list repositories: %v", err)
		}
		for _, repo := range repos {
			if hasCommitInLast12Months(repo.PushedAt) && !includesInCsv(os.Getenv("EXCLUDE_REPOS"), *repo.Name) {
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
	gc.AllRepos = allRepos

	gc.prepareLanguageMapByContentByte()
	gc.prepareSortedLanguageByUsage()

	return gc
}

func (gc *gitHubClient) prepareLanguageMapByContentByte() {
	languages := make(map[string]int)
	for _, repo := range gc.AllRepos {
		langs, _, err := gc.client.Repositories.ListLanguages(gc.ctx, *repo.Owner.Login, *repo.Name)
		if err != nil {
			log.Printf("Failed to list languages for repository %s: %v", *repo.FullName, err)
			continue
		}
		for lang, contentByte := range langs {
			if !includesInCsv(os.Getenv("EXCLUDE_LANGS"), lang) {
				languages[lang] += contentByte
			}
		}
	}
	gc.LangMap = languages
}

func (gc *gitHubClient) prepareSortedLanguageByUsage() {
	var topLanguages []string
	for lang := range gc.LangMap {
		topLanguages = append(topLanguages, lang)
	}
	sort.Slice(topLanguages, func(i, j int) bool {
		return gc.LangMap[topLanguages[i]] > gc.LangMap[topLanguages[j]]
	})
	gc.SortedLang = topLanguages
}
