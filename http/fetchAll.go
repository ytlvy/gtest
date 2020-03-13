// go run fetchAll.go https://www.baidu.com http://www.kuwo.cn
//

package http

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type FetchAllManage struct {
}

func (c *FetchAllManage) BeginFetch() {
	start := time.Now()
	ch := make(chan string)

	for _, url := range os.Args[1:] {
		if len(url) < 1 {
			fmt.Fprintf(os.Stderr, "please input url")
			os.Exit(1)
		}

		go c.fetch(url, ch)
	}
	for range os.Args[1:] {
		fmt.Println(<-ch)
	}

	fmt.Printf("%0.2fs elapsed\n", time.Since(start).Seconds())
}

func (c *FetchAllManage) fetch(url string, ch chan<- string) {
	start := time.Now()

	if !strings.HasPrefix(url, "http") {
		url = fmt.Sprintf("http://%s", url)
	}
	resp, err := http.Get(url)
	if err != nil {
		// fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
		// os.Exit(1)
		ch <- fmt.Sprint(err)
		return
	}

	nbytes, err := io.Copy(ioutil.Discard, resp.Body) //ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		ch <- fmt.Sprintf("while reading %s %v", url, err)
		return
	}

	secs := time.Since(start)
	ch <- fmt.Sprintf("%0.2fs  %7d %s", secs, nbytes, url)

}
