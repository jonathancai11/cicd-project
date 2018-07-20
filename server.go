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
	// "bytes"

	"github.com/google/go-github/github"
	// "golang.org/x/oauth2"
)

const (
	address = ""
)

var (
	config GithubConfig
)

func main() {
	// p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
	// p1.save()
	// p2, _ := loadPage("TestPage")
	// fmt.Println(string(p2.Body))

	// fs := http.FileServer(http.Dir("/Users/jonathancai/go/src/cicd-project/static/"))
	// http.Handle("/", fs)
	// http.HandleFunc("/api/", APIHandler) // API handles interaction with frontend
	// http.HandleFunc("/view/", viewHandler)
	// fmt.Println("Listening and serving on port :8080...")
	// log.Fatal(http.ListenAndServe(":8080", nil))
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil {
		fmt.Println("Error loading page")
		p = &Page{Title: title}
	}
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func APIHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("APIHandler called")
	fmt.Println(r.Method + " request")
	switch r.Method {
	case "GET":
		fmt.Println(config.Token)
		repos := GetRepos(config.Token)
		// var buffer bytes.Buffer
		for _, repo := range repos {
			// buffer.WriteString(repo)
			fmt.Println(repo)
		}
		// result := buffer.String()
		// fmt.Println(result)
		// w.Write([]bytes(result))
	case "POST":
		b, er := ioutil.ReadAll(r.Body)
		if er != nil {
			log.Println("Error reading response")
			return
		}
		err := json.Unmarshal(b, &config)
		if err != nil {
			log.Println("Error unmarshalling payload")
		}
		fmt.Println("Resulting github config: ", config.Token)
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
