package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type PreCommitCheck struct {
}

func (c *PreCommitCheck) Run() {

	workpath, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	cmd := exec.Command("git", "diff", "--cached", "--name-only", "--diff-filter=d", "--", "*.mm", "*.m", "*.go", "*.h")
	// cmd.Dir = path
	var out bytes.Buffer
	cmd.Stdout = &out
	err1 := cmd.Run()
	if err1 != nil {
		log.Fatal(err1)
	}
	// fmt.Printf("in all caps: %q\n", out.String())
	paths := strings.Split(out.String(), "\n")
	for _, fpath := range paths {
		if len(fpath) < 1 {
			continue
		}
		fpath = workpath + string(os.PathSeparator) + fpath
		// fmt.Println(fpath)

		data, err := ioutil.ReadFile(fpath)
		check(err)

		// fmt.Println(string(data))
		r := regexp.MustCompile(`.*?(https:|http:)?([^/;\n+]*)test(\.|-)([^/;\n+]*)kuwo`)
		mdata := r.Find(data)
		if len(mdata) < 1 {
			continue
		}

		// fmt.Printf("%s\n", string(mdata))
		if strings.HasPrefix(string(mdata), "//") {
			continue
		}

		fmt.Printf("fpath: %s contain test url \n", fpath)
		fmt.Printf("content: %s \n", string(mdata))
		os.Exit(1)
	}

	fmt.Println("this commit is ok")

}

func check(e error) {
	if e != nil {
		fmt.Println(e)
		panic(e)
	}
}

func main() {
	check := &PreCommitCheck{}
	check.Run()
}
