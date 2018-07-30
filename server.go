package main

import (
	"bytes"
	// "context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/go-github/github"
	// "golang.org/x/oauth2"
)

const (
	projectaddr = "/Users/jonathancai/go/src/cicd-project/"
	// projectaddr = "/root/go/src/cicd-project/"
)

var (
	config GithubConfig
	token  string
)

type HooksReturn struct {
	Hooks [][]string
}

func main() {
	fs := http.FileServer(http.Dir(projectaddr + "static/"))
	http.Handle("/", fs)
	http.HandleFunc("/api/", APIHandler)         // API handles interaction with frontend
	http.HandleFunc("/webhook/", WebhookHandler) // Webhook handles the actual webhook payloads
	fmt.Println("Listening and serving on port :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
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
			GetHooksClient(config.Token, config.Username, config.Repo)
			// result := GetHooksClient(config.Token, config.Username, config.Repo)
			// a, err := json.Marshal(result)
			// w.Write(a)
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

/* Notes:
- When creating a webhook, Github will initially send a request (which does not contain commits).
- Further payloads (push) will contain all commits or commithead (latest commit).
- Webhooks can be configured for specific events (default is push).
*/

// WebhookHandler handles the actual payloads from Github
func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Webhook Handler received " + r.Method + " request")
	switch r.Method {
	case "GET":
	case "POST":
		fmt.Println("URL:", r.URL.Path)
		fmt.Println("Webhook handler received post request")
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
		if headcom != nil {
			fmt.Println("Head Commit:", "- Author:", *headcom.Author.Name, "ID:", *headcom.ID, "Time:", *headcom.Timestamp)
			// IF payload is push
			commits := payload.Commits
			fmt.Println("List of commits:")
			for idx, com := range commits {
				fmt.Println("Commit #", idx, "- Author:", com.Author.Name, "ID:", com.ID, "Time:", com.Timestamp, "Modified:", com.Modified, "Removed:", com.Removed)
			}
		}
		// fmt.Println(payload)
	case "DELETE":
	default:
	}
}
