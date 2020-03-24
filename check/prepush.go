package check

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

type PrePushCheck struct {
}

func (c *PrePushCheck) Run() {

	cmd := exec.Command("git", "diff", "--cached", "--diff-filter=d", "--", "'*.go'")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("in all caps: %q\n", out.String())

}
