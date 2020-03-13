package main

import (
	"os"

	"github.com/ytlvy/gtest/gif"
	"github.com/ytlvy/gtest/http"
)

func main() {
	runServer()
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
