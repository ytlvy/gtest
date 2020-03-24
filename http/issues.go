package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const IssuesURL = "https://api.github.com/search/issues"

type IssuesSearchResult struct {
	TotalCount int `json:"total_count"`
	Items      []*Issue
}

type Issue struct {
	Number   int
	HTMLURL  string `json:"html_url"`
	Title    string
	State    string
	User     *User
	CreateAt time.Time `json:"create_at"`
	Body     string    //Markdown
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

type IssuesGeter struct {
}

func (i *IssuesGeter) Run() {
	result, err := i.SearchIssues([]string{"repo:golang/go", "is:open", "json", "decoder"})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d issues:\n", result.TotalCount)
	for _, item := range result.Items {
		fmt.Printf("#%-5d %9.9s %.55s\n", item.Number, item.User.Login, item.Title)
	}
}

func (i *IssuesGeter) SearchIssues(terms []string) (*IssuesSearchResult, error) {
	q := url.QueryEscape(strings.Join(terms, " "))
	urlpath := IssuesURL + "?q=" + q
	resp, err := http.Get(urlpath)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("url: %s resp: %s\n", urlpath, string(bytes))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search query faild: %s", resp.Status)
	}

	var result IssuesSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
