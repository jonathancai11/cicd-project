package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type EventType string

var (
	Push        EventType = "push"
	PullRequest EventType = "pull_request"
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
*/

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

func GetGitHooks(token, username, repo string) map[string][]string {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	hooks, resp, err := client.Repositories.ListHooks(ctx, username, repo, nil)
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

func CreateGitHook(username, reponame, trgURL, token string, events []string) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	config := make(map[string]interface{})
	config["url"] = trgURL
	name := "web"
	hook := github.Hook{
		Name:   &name,
		Events: events,
		Config: config,
	}
	_, resp, err := client.Repositories.CreateHook(ctx, username, reponame, &hook)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Webhook create response:", resp)
}

func DeleteGitHook(token, username, repo string, hookid int64) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	resp, err := client.Repositories.DeleteHook(ctx, username, repo, hookid)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Delete hook response:", resp)
}

func EditGitHook(token, username, repo, trgURL string, events []string, hookid int64) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	config := make(map[string]interface{})
	config["url"] = trgURL
	name := "web"
	hook := github.Hook{
		Name:   &name,
		Events: events,
		Config: config,
	}
	_, resp, err := client.Repositories.EditHook(ctx, username, repo, hookid, &hook)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Edit hook response:", resp)
}
