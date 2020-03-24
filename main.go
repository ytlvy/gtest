package main

import (
	"fmt"
	"os"

	"github.com/ytlvy/gtest/check"
	"github.com/ytlvy/gtest/gif"
	"github.com/ytlvy/gtest/http"
	"github.com/ytlvy/gtest/rest"
)

func main() {
	// tcase := &mycase.MyCase{}
	// tcase.TestInterfacePara()
	// println("hello " + "world")

	testCheck()
}

func testCheck() {
	check := &check.PreCommitCheck{}
	check.Run()
}

func testFetch() {
	fetcher := &http.FetchManager{}
	fetcher.RunFetch("http://www.baidu.com")
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
