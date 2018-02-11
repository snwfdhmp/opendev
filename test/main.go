package main

import (
	"fmt"
	"os"

	"gopkg.in/src-d/go-git.v4"
)

func main() {
	fmt.Println("PRE RECEIVE HOOK START")
	gitRepo, err := git.PlainOpen("./")
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

	fmt.Println("env('GITHUB_USER_LOGIN'):", os.Getenv("GITHUB_USER_LOGIN"))
	fmt.Println("env('GIT_DIR'):", os.Getenv("GIT_DIR"))
	fmt.Println("env('GITHUB_USER_IP'):", os.Getenv("GITHUB_USER_IP"))
	fmt.Println("env('GITHUB_REPO_NAME'):", os.Getenv("GITHUB_REPO_NAME"))
	fmt.Println("env('GITHUB_PULL_REQUEST_AUTHOR_LOGIN'):", os.Getenv("GITHUB_PULL_REQUEST_AUTHOR_LOGIN"))
	fmt.Println("env('GITHUB_REPO_PUBLIC'):", os.Getenv("GITHUB_REPO_PUBLIC"))
	fmt.Println("env('GITHUB_PUBLIC_KEY_FINGERPRINT'):", os.Getenv("GITHUB_PUBLIC_KEY_FINGERPRINT"))
	fmt.Println("env('GITHUB_PULL_REQUEST_HEAD'):", os.Getenv("GITHUB_PULL_REQUEST_HEAD"))
	fmt.Println("env('GITHUB_PULL_REQUEST_BASE'):", os.Getenv("GITHUB_PULL_REQUEST_BASE"))
	fmt.Println("env('GITHUB_VIA'):", os.Getenv("GITHUB_VIA"))

	fmt.Println("PRE RECEIVE HOOK END")
}
