package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ytlvy/gtest/chat"

	"github.com/ytlvy/gtest/count"
	"github.com/ytlvy/gtest/gif"
	"github.com/ytlvy/gtest/http"
	"github.com/ytlvy/gtest/rest"
)

func main() {
	// tcase := &mycase.MyCase{}
	// tcase.TestInterfacePara()
	// println("hello " + "world")
	runChater()
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
