package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type EventType string

// type Hook struct {
// 	URL    string
// 	Events []*EventType
// }

var (
	Push        EventType = "push"
	PullRequest EventType = "pull_request"
)

type GithubWebhookCreator struct {
	Name   string
	Active bool
	Events []EventType
	Config GithubWebhookCreatorConfig
}

type GithubWebhookCreatorConfig struct {
	Url         string
	ContentType string
	Secret      string
	InsecureSSL string
}

type GithubConfig struct {
	Token    string
	Username string
	Repo     string
	// CreateWebhookConfig CreateWebhookConfig
}

type CreateWebhookConfig struct {
	HookName   string
	PushConfig string
	PullConfig string
}

/* Notes:
- Should be using token instead of username/password (allows for more flexible authentication)
- In order to manage webhooks, personal tokens must be configured to allow (repo_hook)
- (404 Error if no repo_hook authorization)
- Add authentication (token) in header of requests or as parameter ("Authorization: token TOKEN")
- Never add backslash to end of API call for url
- You can create up to 20 hooks for a single
*/

// CreateGithubWebhook creates a github webhook.
// srcUrl should be https://api.github.com/repos/:owner/:repo/hooks
// name just gives the webhook a name
// trgUrl is the target url that we want the payload to be sent to
// events are the type of events we want the webhook to post for (defaults to "push")
func CreateGithubWebhook(hookname string, username string, reponame string, trgUrl string, events []EventType, token string) {
	fmt.Println("Creating webhook...")
	srcUrl := "https://api.github.com/repos/" + username + "/" + reponame + "/hooks?access_token=" + Token
	fmt.Println("url:", srcUrl)
	config := GithubWebhookCreatorConfig{
		Url:         srcUrl,
		ContentType: "json",
	}
	payload := GithubWebhookCreator{
		Name:   hookname,
		Active: true,
		Events: events,
		Config: config,
	}
	fmt.Println("Created payload")
	b, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", trgUrl, bytes.NewBuffer(b))
	req.Header.Set("Authorization", " token "+token)
	fmt.Println("Created request")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response:", err)
	}
	fmt.Printf("%s\n", string(contents))
}

// GetRepos returns a list of repos authorized to view under given token
func GetRepos(token string) []string {
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

func GetHooksClient(token, username, repo string) [][]string {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	hooks, _, err := client.Repositories.ListHooks(ctx, username, repo, nil)
	if err != nil {
		log.Println(err)
	}
	ret := make([][]string, 0)
	fmt.Println("Hooks:")
	for idx, hook := range hooks {
		a := make([]string, 0)
		a = append(a, *hook.URL)
		for _, event := range hook.Events {
			a = append(a, event)
		}
		fmt.Println("Hook #", idx+1, "- Created:", *hook.CreatedAt, "ID:", *hook.ID, "Config:", hook.Config, "Events:", hook.Events)
		ret = append(ret, a)
	}
	return ret
}
