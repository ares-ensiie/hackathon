package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func request(w http.ResponseWriter, r *http.Request) {

	t, _ := template.ParseFiles("index.html")
	t.Execute(w, nil)
}

type repos_hash struct {
	url       string
	hashtable map[string]int
}

func gather(repo *github.Repository, client *github.Client, out chan repos_hash) {
	var repos_tmp repos_hash
	//filter them
	fmt.Println("repo_URL : ", *repo.URL)
	fmt.Println("repo_Owner : ", *repo.Owner.Login)
	fmt.Println("repo_Owner : ", *repo.Name)

	res, _, err := client.Repositories.ListLanguages(*repo.Owner.Login, *repo.Name)
	if err != nil {
		fmt.Println(err)
	}
	repos_tmp.hashtable = res
	repos_tmp.url = *repo.URL

	out <- repos_tmp
}

func search(w http.ResponseWriter, r *http.Request) {
	/*Get the form content*/
	/*r.ParseForm()
		// logic part of log in
		fmt.Println("language: ", r.Form["language"])
	/*	language := strings.Join(r.Form["language"], "")*/

	//start the query on github
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "b1abf4e31f153ad21e19cf70dabce2310a731b1c"},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	//get the 100 last repositories
	opt := &github.RepositoryListAllOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	// get all pages of results
	var allRepos []*github.Repository
	repos, _, err := client.Repositories.ListAll(opt)
	if err != nil {
		fmt.Println(err)
	}
	allRepos = append(allRepos, repos...)
	reposChan := make(chan repos_hash, len(allRepos))
	for _, repo := range allRepos {
		go gather(repo, client, reposChan)
	}

	result := make([]repos_hash, len(allRepos))
	for i := 0; i < len(allRepos); i++ {
		result[i] = <-reposChan
		fmt.Println(result[i])
	}

	// for _, elmt := range allRepos {
	// 	fmt.Printf("elmt :%s\n", elmt.url)
	// }

	t, _ := template.ParseFiles("result.html")
	t.Execute(w, nil)
}

func main() {

	//On initialise la valeur du port
	port := flag.String("port", "3000", "server port number")

	http.HandleFunc("/", request)
	http.HandleFunc("/search", search)

	log.Println("Listening on :" + *port)

	//manage request
	log.Println("Listenning...")
	err := http.ListenAndServe(":"+*port, nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
