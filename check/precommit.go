package check

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"regexp"
)

type PreCommitCheck struct {
}

func (c *PreCommitCheck) Run() {

	// path, err := os.Getwd()
	// if err != nil {
	// 	log.Println(err)
	// }

	cmd := exec.Command("git", "diff", "--cached", "--diff-filter=d") //, "--", "'*.go'")
	// cmd.Dir = path
	var out bytes.Buffer
	cmd.Stdout = &out
	err1 := cmd.Run()
	if err1 != nil {
		log.Fatal(err1)
	}
	// fmt.Printf("in all caps: %q\n", out.String())

	matched, err := regexp.Match(`^([^/;\n+]*)(https:|http:)?([^/;\n+]*)test(\.|-)([^/;\n+]*)kuwo`, out.Bytes())
	fmt.Println(matched, err)
}

// func main() {
// 	check := &PreCommitCheck{}
// 	check.Run()
// }
