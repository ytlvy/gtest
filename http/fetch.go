//

package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type FetchManager struct {
}

func (f *FetchManager) RunFetch(url string) {

	if len(url) < 1 {
		fmt.Fprintf(os.Stderr, "please input url")
		os.Exit(1)
	}

	if !strings.HasPrefix(url, "http") {
		url = fmt.Sprintf("http://%s", url)
	}
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
		os.Exit(1)
	}

	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
		os.Exit(1)
	}

	// fmt.Printf("%s", b)
	doc, err := html.Parse(bytes.NewReader(b))
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
		os.Exit(1)
	}
	for _, link := range f.outline(nil, "", doc) {
		fmt.Println(link)
	}
}

func (f *FetchManager) outline(stack []string, prefix string, n *html.Node) []string {
	if n.Type == html.ElementNode {
		stack = append(stack, prefix+"  "+n.Data)
		prefix += " " + n.Data
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		stack = f.outline(stack, prefix, c)
	}

	return stack
}

func (f *FetchManager) visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = f.visit(links, c)
	}

	return links
}
