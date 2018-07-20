package main

import (
	// "bytes"
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

func main() {
	fs := http.FileServer(http.Dir("/Users/jonathancai/go/src/cicd-project/static/"))
	http.Handle("/", fs)
	http.HandleFunc("/api/", APIHandler) // API handles interaction with frontend
	// http.HandleFunc("/webhook/", WebhookHandler) // Webhook handles interaction with github/bitbucket
	fmt.Println("Listening and serving on port :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func APIHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("APIHandler called")
	fmt.Println(r.Method + " request")
	switch r.Method {
	case "GET":
	case "POST":
		b, er := ioutil.ReadAll(r.Body)
		if er != nil {
			log.Println("Error reading response")
			return
		}
		var config GithubConfig
		err := json.Unmarshal(b, &config)
		if err != nil {
			log.Println("Error unmarshalling payload")
		}
		fmt.Println("Resulting github config: ", config.Url, config.Token)
		PrintRepos(config.Token)
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
