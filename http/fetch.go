//

package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

//FetchManager url 抓取类
type FetchManager struct {
	tokens chan struct{}
}

//NewFetch return a Fetchmanager
func NewFetch() *FetchManager {
	tokens := make(chan struct{}, 1)
	fetch := &FetchManager{tokens}
	return fetch
}

//RunFetch 获取 url 连接
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

func (f *FetchManager) forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		f.forEachNode(c, pre, post)
	}

	if post != nil {
		post(n)
	}
}

// ExtractUrls 提取 url 连接
func (f *FetchManager) ExtractUrls(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("getting %s error: %s", url, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML : %v", url, err)
	}

	var links []string

	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := resp.Request.URL.Parse(a.Val)
				if err == nil {
					links = append(links, link.String())
				}
			}
		}
	}
	f.forEachNode(doc, visitNode, nil)
	return links, nil
}

//RunCrawl 抓取连接以及相关子连接
func (f *FetchManager) RunCrawl(url string) {
	worklist := []string{url}
	f.breadthFirst(f.crawl, worklist)
}

//RunCrawWithChan 爬虫
func (f *FetchManager) RunCrawWithChan() {
	f.breadthFirstChan(f.crawl)
}

func (f *FetchManager) breadthFirstChan(craw func(item string) []string) {
	workchon := make(chan []string)
	numOfWork := 0
	numOfWork++

	go func() {
		workchon <- os.Args[1:]
	}()

	seen := make(map[string]bool)

	for ; numOfWork > 0; numOfWork-- {
		list := <-workchon
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				numOfWork++
				go func(link string) {
					fmt.Printf("======>>begin fetch: %s\n", link)
					workchon <- craw(link)
				}(link)
			}
		}
	}
}

func (f *FetchManager) breadthFirst(fun func(item string) []string, worklist []string) {
	seen := make(map[string]bool)
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				//append的参数“f(item)...”，会将f返回的一组元素一个个添加到worklist中
				fmt.Printf("======>>begin fetch: %s\n", item)
				worklist = append(worklist, fun(item)...)
			}
		}
	}
}

func (f *FetchManager) crawl(url string) []string {
	// fmt.Println(url)

	f.tokens <- struct{}{}
	list, err := f.ExtractUrls(url)
	<-f.tokens
	if err != nil {
		log.Print(err)
	}
	return list
}
