package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type EventType string

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
		panic(err)
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
	}
	fmt.Printf("%s\n", string(contents))
}

// GetHooks retrieves all hooks active/inactive for a username/reponame
// Only prints out hooks right now...
func GetHooks(username string, reponame string, token string) {

	url := "https://api.github.com/repos/" + username + "/" + reponame + "/hooks"
	fmt.Println("url:", url)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", " token "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
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
		fmt.Println(err)
	}
	result := make([]string, 0)
	for _, repo := range repos {
		a := *repo.Name
		result = append(result, a)
	}
	return result
}
