package main

import (
	"bytes"
	"fmt"
	"log"
	// "strings"
	// "context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	// "net/url"

	"github.com/google/go-github/github"
	// "golang.org/x/oauth2"
)

const (
	address = ""
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

func main() {
	fs := http.FileServer(http.Dir("/Users/jonathancai/go/src/cicd-project/static/"))
	http.Handle("/", fs)
	http.HandleFunc("/api/", APIHandler) // API handles interaction with frontend
	// http.HandleFunc("/webhook/", WebhookHandler) // Webhook handles interaction with github/bitbucket
	fmt.Println("Listening and serving on port :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))

	// ctx := context.Background()
	// ts := oauth2.StaticTokenSource(
	// 	&oauth2.Token{AccessToken: token},
	// )
	// tc := oauth2.NewClient(ctx, ts)

	// client := github.NewClient(tc)

	// // 	// list all repositories for the authenticated user
	// repos, _, err := client.Repositories.List(ctx, "", nil)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// for _, repo := range repos {
	// 	fmt.Println(*repo.Name, ":", *repo.HooksURL)
	// }
}

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
	req.Header.Set("Authorization", " token "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

type GithubConfig struct {
	Url   string
	Token string
}

func APIHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("API RECEIVED REQUEST")
	fmt.Println(r.Method)
	switch r.Method {
	case "GET":
	case "POST":
		fmt.Println("POST REQUEST")
		b, er := ioutil.ReadAll(r.Body)
		if er != nil {
			log.Println("Error reading response")
			return
		}
		fmt.Println(b)
		var config GithubConfig
		err := json.Unmarshal(b, &config)
		if err != nil {
			log.Println("Error unmarshalling payload")
		}
		fmt.Println(config.Url, config.Token)
	case "DELETE":
	default:
	}
}

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
	case "POST":
		b, er := ioutil.ReadAll(r.Body)
		if er != nil {
			log.Println("Error reading response")
			return
		}
		var payload github.WebHookPayload
		err := json.Unmarshal(b, &payload)
		if err != nil {
			log.Println("Error unmarshalling payload")
		}
	case "DELETE":
	default:
	}
}
