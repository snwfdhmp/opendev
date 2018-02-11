package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/src-d/go-git.v4"
)

func main() {
	fmt.Println("PRE RECEIVE HOOK START")
	abs, _ := filepath.Abs("./")
	fmt.Println("wd:", abs)
	gitRepo, err := git.PlainOpen("./..")
	if err != nil {
		fmt.Println("fatal:", err)
		return
	}

	head, err := gitRepo.Head()
	if err != nil {
		fmt.Println("fatal:", err)
		return
	}

	fmt.Println("head:", head.String())

	fmt.Println("PRE RECEIVE HOOK END")
}
