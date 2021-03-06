package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// type EventType string

var (
	// Push        EventType = "push"
	// PullRequest EventType = "pull_request"
	PullRequest = "pull_request"
	Push        = "push"
)

type GithubConfig struct {
	Token    string
	Username string
	Repo     string
}

type CreateGitHookConfig struct {
	PushConfig bool
	PullConfig bool
	HookURL    string
}

/* Notes:
- Should be using token instead of username/password (allows for more flexible authentication)
- In order to manage webhooks, personal tokens must be configured to allow (repo_hook)
- (404 Error if no repo_hook authorization)
- Add authentication (token) in header of requests or as parameter ("Authorization: token TOKEN")
- Never add backslash to end of API call for url
- You can create up to 20 hooks for a single repo???

- Hook name should be "web" (a few other supported names i suppose are unnecessary)
- There are many (around 20 events that you can listen to), but it defaults to "push" for commits
- Source URL should be https://api.github.com/repos/:owner/:repo/hooks
- Target URL just has to be valid IP/url address, will send initial packaged to confirm
- Wildcard event trigger = *, listens to ALL events
*/

// GetGitRepos takes in a Github personal token and returns a string of all repos authorized under the token
func GetGitRepos(token string) []string {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	repos, _, err := client.Repositories.List(ctx, "", nil)
	if err != nil {
		log.Println(err)
	}
	result := make([]string, 0)
	for _, repo := range repos {
		a := *repo.Name
		result = append(result, a)
	}
	return result
}

// GetGitHooks returns hooks represented by a map of targetURL to event triggers
func GetGitHooks(config GithubConfig) map[string][]string {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	hooks, resp, err := client.Repositories.ListHooks(ctx, config.Username, config.Repo, nil)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("List webhooks response:", resp)
	ret := make(map[string][]string)
	fmt.Println("List of hooks:")
	for _, hook := range hooks {
		ret[hook.Config["url"].(string)] = hook.Events
		fmt.Println("Hook - "+*hook.Name+" Created:", *hook.CreatedAt, "ID:", *hook.ID, "Config:", hook.Config, "Events:", hook.Events)
	}
	return ret
}

// CreateGitHook creates a hook listening to events and sending payloads to trgURL under given githubconfig
func CreateGitHook(config GithubConfig, trgURL string, events []string) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	hookconfig := make(map[string]interface{})
	hookconfig["url"] = trgURL
	name := "web"
	hook := github.Hook{
		Name:   &name,
		Events: events,
		Config: hookconfig,
	}
	_, resp, err := client.Repositories.CreateHook(ctx, config.Username, config.Repo, &hook)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Webhook create response:", resp)
}

// DeleteGitHook deletes hook corresponding to hookid under the given githubconfig
func DeleteGitHook(config GithubConfig, hookid int64) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	resp, err := client.Repositories.DeleteHook(ctx, config.Username, config.Repo, hookid)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Delete hook response:", resp)
}

// DeleteGitHook edits hook corresponding to hookid under the given githubconfig, replaces hook configurations with trgURL and events
func EditGitHook(config GithubConfig, trgURL string, events []string, hookid int64) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	hookconfig := make(map[string]interface{})
	hookconfig["url"] = trgURL
	name := "web"
	hook := github.Hook{
		Name:   &name,
		Events: events,
		Config: hookconfig,
	}
	_, resp, err := client.Repositories.EditHook(ctx, config.Username, config.Repo, hookid, &hook)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Edit hook response:", resp)
}
