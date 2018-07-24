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

const (
	projectaddr = "/Users/jonathancai/go/src/cicd-project/"
	// projectaddr = "/root/go/src/cicd-project/"
)

var (
	config GithubConfig
	token  string
)

func main() {
	fs := http.FileServer(http.Dir(projectaddr + "static/"))
	http.Handle("/", fs)
	http.HandleFunc("/api/", APIHandler)     // API handles interaction with frontend
	http.HandleFunc("/webhook/", APIHandler) // Webhook handles the actual webhook payloads
	fmt.Println("Listening and serving on port :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))

	// hookname := "test"
	// githubapi := "https://api.github.com/"
	// username := "jonathancai11"
	// reponame := "cicd-project"
	// trgurl := "http://192.168.1.192:8080"
	// var events []EventType
	// events = append(events, Push)
	// CreateGithubWebhook(hookname, username, reponame, trgurl, events, Token)
}

// APIHandler deals with JS
func APIHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("APIHandler called:", r.Method+" request")
	switch r.Method {
	case "GET":
		switch r.URL.Path {
		case "/api/token": // Get saved token
			fmt.Println("Getting saved token...")
			b, err := ioutil.ReadFile(projectaddr + "config.txt")
			if err != nil {
				log.Println("Error reading file")
				return
			}
			s := string(b[:])
			fmt.Println("Saved token:", s)
			token = s
			w.Write([]byte(b))
		case "/api/repos": // Get repositories
			fmt.Println("Getting repos...")
			repos := GetRepos(token) // (must save token before calling get repos)
			var buffer bytes.Buffer
			for _, repo := range repos {
				buffer.WriteString(repo + ", ")
				// fmt.Println(repo)
			}
			result := buffer.String()
			fmt.Println("Repos:", result)
			w.Write([]byte(result))
		}
	case "POST":
		switch r.URL.Path {
		case "/api/token": // Save token
			fmt.Println("Saving token")
			b, er := ioutil.ReadAll(r.Body)
			if er != nil {
				log.Println("Error reading response")
				return
			}
			var config GithubConfig
			err := json.Unmarshal(b, &config)
			if err != nil {
				fmt.Println("Error unmarshalling token")
			}
			err = ioutil.WriteFile(projectaddr+"config.txt", []byte(config.Token), 0644)
			if err != nil {
				fmt.Println("Error writing to config file")
				return
			}
			token = config.Token // save token on server
			fmt.Println("Saved token: " + config.Token)
		case "/api/list": // Get list of hooks
			b, er := ioutil.ReadAll(r.Body)
			if er != nil {
				log.Println("Error reading response")
				return
			}
			var config GithubConfig
			err := json.Unmarshal(b, &config)
			if err != nil {
				log.Println("Error unmarshalling github config")
			}
			// ctx := context.Background()
			// ts := oauth2.StaticTokenSource(
			// 	&oauth2.Token{AccessToken: token},
			// )
			// tc := oauth2.NewClient(ctx, ts)
			// client := github.NewClient(tc)
			// rs := client.Repositories
			// lsopts := github.ListOptions{}
			// rs.ListHooks(ctx, config.Username, config.Repo, lsopts)
			GetHooks(config.Username, config.Repo, config.Token)
		case "/api/create": // Create a hook
			fmt.Println("Creating hook")
			b, er := ioutil.ReadAll(r.Body)
			if er != nil {
				log.Println("Error reading response")
				return
			}
			return
			// NEED TO FIGURE OUT THIS PART:
			var config GithubConfig
			err := json.Unmarshal(b, &config)
			if err != nil {
				log.Println("Error unmarshalling token")
			}
			err = ioutil.WriteFile(projectaddr+"config.txt", []byte(config.Token), 0644)
			if err != nil {
				log.Println("Error writing to config file")
				return
			}
		}
	case "DELETE":
	default:
	}
}

// WebhookHandler handles the actual payloads from Github
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
		fmt.Println(payload)
	case "DELETE":
	default:
	}
}
