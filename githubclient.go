package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
	Token string
}

// CreateGithubWebhook creates a github webhook.
// srcUrl should be https://api.github.com/repos/:owner/:repo/hooks
// name just gives the webhook a name
// trgUrl is the target url that we want the payload to be sent to
// events are the type of events we want the webhook to post for (defaults to "push")
func CreateGithubWebhook(name string, srcUrl string, trgUrl string, events []EventType) {
	config := GithubWebhookCreatorConfig{
		Url:         srcUrl,
		ContentType: "json",
	}
	payload := GithubWebhookCreator{
		Name:   name,
		Active: true,
		Events: events,
		Config: config,
	}
	b, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", trgUrl, bytes.NewBuffer(b))
	req.Header.Set("Authorization", " token "+Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func GetRepos(token string) []string {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	// fmt.Println("Created client")
	// list all repositories for the authenticated user
	repos, _, err := client.Repositories.List(ctx, "", nil)
	if err != nil {
		fmt.Println(err)
	}
	result := make([]string, 0)
	for _, repo := range repos {
		// fmt.Println(*repo.Name, ":", *repo.HooksURL)
		a := *repo.Name
		result = append(result, a)
	}
	return result
}
