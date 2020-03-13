package main

import (
	"github.com/ytlvy/gtest/gif"
	"github.com/ytlvy/gtest/http"
)

func main() {
	runGif()
}

func runGif() {
	gif := new(gif.Maker)
	gif.MakeGif()
}

func runFetch() {
	fetAll := new(http.FetchAllManage)
	fetAll.BeginFetch()
}

func runFetchAll() {
	fetAll := new(http.FetchAllManage)
	fetAll.BeginFetch()
}
