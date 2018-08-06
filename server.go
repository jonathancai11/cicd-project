package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/go-github/github"
)

// Currently, config.txt saves the personal token, saved on server as "token" when page is loaded and calls GET request.
// I use "config" as a global var to save token, username, and repo.
// projectaddr is used to find static files to serve on server.

// Server is not completely finished:
//   -  edit/delete webhooks is not hooked up to frontend yet
//   -  bitbucket is not done
//   -  not all event triggers are configured (only can create hooks listenening to push or pull request)

const (
	projectaddr = "/Users/jonathancai/go/src/cicd-webhooks/"
	// projectaddr = "/root/go/src/cicd-webhooks/"
)

var (
	config GithubConfig
	token  string
)

func main() {
	fs := http.FileServer(http.Dir(projectaddr + "static/"))
	http.Handle("/", fs)
	http.HandleFunc("/api/", APIHandler)               // API handles interaction with frontend
	http.HandleFunc("/webhook/", GithubWebhookHandler) // Webhook handles the actual webhook payloads
	fmt.Println("Listening and serving on port :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// APIHandler deals with JS
func APIHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("APIHandler called:", r.Method+" request")
	fmt.Println("Current config:", config)
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
			repos := GetGitRepos(token) // (must save token before calling get repos)
			var buffer bytes.Buffer
			for _, repo := range repos {
				buffer.WriteString(repo + ", ")
				// fmt.Println(repo)
			}
			result := buffer.String()
			fmt.Println("Repos:", result)
			w.Write([]byte(result))
		case "/api/list": // Get hooks
			fmt.Println("Getting hooks...")
			hooks := GetGitHooks(config)
			res, err := json.Marshal(hooks)
			if err != nil {
				fmt.Println("Error marshalling hooks:", err)
			}
			w.Write(res)
		}
	case "POST":
		switch r.URL.Path {
		case "/api/config": // Save config
			fmt.Println("Saving config")
			b, er := ioutil.ReadAll(r.Body)
			if er != nil {
				log.Println("Error reading response")
				return
			}
			err := json.Unmarshal(b, &config)
			if err != nil {
				fmt.Println("Error unmarshalling token")
			}
			err = ioutil.WriteFile(projectaddr+"config.txt", []byte(config.Token), 0644)
			if err != nil {
				fmt.Println("Error writing to config file")
				return
			}
			fmt.Println("Saved config:", config)
		case "/api/create": // Create a hook
			fmt.Println("Creating hook")
			b, er := ioutil.ReadAll(r.Body)
			if er != nil {
				log.Println("Error reading response")
				return
			}
			var create CreateGitHookConfig
			err := json.Unmarshal(b, &create)
			if err != nil {
				log.Println("Error unmarshalling token")
			}
			fmt.Println("NEW HOOK URL:", create.HookURL, "PUSH:", create.PushConfig, "PULL", create.PullConfig)
			events := make([]string, 0)
			// Here, parse through event triggers (right now, only push and pull request are configured)
			if create.PushConfig {
				events = append(events, Push)
			}
			if create.PullConfig {
				events = append(events, PullRequest)
			}
			CreateGitHook(config, create.HookURL, events)
		}
	case "DELETE": // Yet to be configured, but already wrote DeleteGitHook function
	default:
	}
}

/* Notes:
- When creating a webhook, Github will initially send a request (which does not contain commits).
- Further payloads (for ex. push) will contain all commits or commithead (latest commit).
- Webhooks can be configured for specific events (default is push).
*/

// GithubWebhookHandler handles the actual payloads from Github
func GithubWebhookHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Webhook Handler received " + r.Method + " request")
	switch r.Method {
	case "POST": // Webhook handler only needs to handle
		fmt.Println("URL path:", r.URL.Path)
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
		headcom := payload.HeadCommit
		if headcom != nil { // Payload will have headcommit when reporting push
			fmt.Println("Head Commit:", "- Author:", *headcom.Author.Name, "ID:", *headcom.ID, "Time:", *headcom.Timestamp)
			// IF payload is push
			commits := payload.Commits
			fmt.Println("List of commits:")
			for idx, com := range commits {
				fmt.Println("Commit #", idx, "- Author:", *com.Author.Name, "ID:", *com.ID, "Time:", com.Timestamp, "Modified:", com.Modified, "Removed:", com.Removed)
			}
		} else { // Contains no commits (either initial payload, or reporting non-push like merge request...)
			fmt.Println(payload)
		}
	default:
		fmt.Println("Webhook handler recieved non-POST request...")
	}
}
