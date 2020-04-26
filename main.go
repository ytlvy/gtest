package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ytlvy/gtest/chat"
	"github.com/ytlvy/gtest/jsonkit"

	"github.com/ytlvy/gtest/count"
	"github.com/ytlvy/gtest/gif"
	"github.com/ytlvy/gtest/http"
	"github.com/ytlvy/gtest/rest"
)

func main() {
	// tcase := &mycase.MyCase{}
	// tcase.TestInterfacePara()
	// println("hello " + "world")
	runJson()
}

func runJson() {
	type Movie struct {
		Title, Subtitle string
		Year            int
		Actor           map[string]string
		Oscars          []string
		Sequel          *string
	}
	strangelove := Movie{
		Title:    "Dr. Strangelove",
		Subtitle: "How I Learned to Stop Worrying and Love the Bomb",
		Year:     1964,
		Actor: map[string]string{
			"Dr. Strangelove":            "Peter Sellers",
			"Grp. Capt. Lionel Mandrake": "Peter Sellers",
			"Pres. Merkin Muffley":       "Peter Sellers",
			"Gen. Buck Turgidson":        "George C. Scott",
			"Brig. Gen. Jack D. Ripper":  "Sterling Hayden",
			`Maj. T.J. "King" Kong`:      "Slim Pickens",
		},
		Oscars: []string{
			"Best Actor (Nomin.)",
			"Best Adapted Screenplay (Nomin.)",
			"Best Director (Nomin.)",
			"Best Picture (Nomin.)",
		},
	}
	encode, _ := jsonkit.Encode(strangelove)
	fmt.Println(string(encode))

	var xiao Movie
	jsonkit.Decode(encode, &xiao)
	// sexpr.Unmarshal(encode, &xiao)
	fmt.Println(xiao)

}

func runChater() {
	chater := &chat.SimpleChater{}
	chater.Start()
}

func runCoundown() {
	count := &count.CountdownManager{}
	count.RunDu1()
}

func testDbServer() {
	dbSer := &http.DBServer{}
	dbSer.RunDbServer()
}

func testCraw() {
	fetcher := http.NewFetch()
	fetcher.RunCrawWithChan()
}

func testFetch1() {
	fetcher := &http.FetchManager{}
	fetcher.RunCrawl("http://www.baidu.com")
}

func testFetch() {
	fetcher := &http.FetchManager{}
	links, _ := fetcher.ExtractUrls("http://www.baidu.com")
	fmt.Println(strings.Join(links, "\n"))
}

func testIssues() {
	getter := &http.IssuesGeter{}
	getter.Run()
}

func testBit() {
	var t int = 1
	fmt.Println(^t)
	fmt.Println(0 ^ t)
}

func runRest() {
	server := rest.NewServer()
	server.Run()
}

func runServer() {
	server := http.NewLServer()
	server.Run()
}

func runGif() {
	gif := new(gif.Maker)
	gif.MakeGif(os.Stdout)
}

func runFetch() {
	fetAll := new(http.FetchAllManage)
	fetAll.BeginFetch()
}

func runFetchAll() {
	fetAll := new(http.FetchAllManage)
	fetAll.BeginFetch()
}
