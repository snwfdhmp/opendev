package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/afero"
	"gopkg.in/src-d/go-git.v4"
	// "gopkg.in/src-d/go-git.v4/plumbing/object"
	//"gopkg.in/src-d/go-git.v4/storage/memory"
	"gopkg.in/yaml.v2"
)

var (
	fs = afero.NewOsFs()
)

type Task struct {
	Name string
	Test string
}

type TaskState struct {
	Task *Task
	Pass bool
}

type RepoState struct {
	TaskStates []TaskState
	Commit     string
}

func main() {
	file, err := afero.ReadFile(fs, "task.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}

	var tasks []Task
	var repo RepoState

	if err = yaml.Unmarshal(file, &tasks); err != nil {
		fmt.Println(err)
		return
	}

	for _, t := range tasks {
		fmt.Print("task['" + t.Name + "'] : ")
		if err := exec.Command("/bin/zsh", "-c", t.Test).Run(); err != nil {
			fmt.Println("FAIL")
			repo.TaskStates = append(repo.TaskStates, TaskState{
				Task: &t,
				Pass: false,
			})
			continue
		}
		fmt.Println("PASS")
		repo.TaskStates = append(repo.TaskStates, TaskState{
			Task: &t,
			Pass: true,
		})
	}

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

	commit, err := gitRepo.CommitObject(head.Hash())
	if err != nil {
		fmt.Println("fatal:", err)
		return
	}

	repo.Commit = commit.String()

	if err := fs.MkdirAll("./.opendev", 0700); err != nil {
		fmt.Println(err)
		return
	}

	out, err := fs.OpenFile("./.opendev/history.yaml", os.O_CREATE|os.O_WRONLY, 0700)
	if err != nil {
		fmt.Println(err)
		return
	}

	b, err := yaml.Marshal(repo)
	if err != nil {
		fmt.Println("fatal:", err)
		return
	}

	if _, err := out.Write(b); err != nil {
		fmt.Println("fatal:", err)
		return
	}
}
