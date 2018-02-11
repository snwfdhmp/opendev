package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/afero"
	"gopkg.in/src-d/go-git.v4"
	// "gopkg.in/src-d/go-git.v4/plumbing/object"
	//"gopkg.in/src-d/go-git.v4/storage/memory"
	"gopkg.in/yaml.v2"
)

var (
	fs = afero.NewOsFs()

	appDir = "./.opendev"
)

type Task struct {
	Name string
	Test string
}

type History struct {
	States map[string]map[string]bool
	Tip    string
}

// commit:
// 	test: true

func (h *History) Run(tasks ...Task) {
	for _, t := range tasks {
		if err := exec.Command("/bin/zsh", "-c", t.Test).Run(); err != nil {
			h.Add(commit, t.Name, false)
			continue
		}
		h.Add(commit, t.Name, true)
	}
}

func (h *History) Add(commit string, testName string, testValue bool) {
	if h.States == nil {
		h.States = make(map[string]map[string]bool)
	}
	if _, ok := h.States[commit]; !ok {
		h.States[commit] = make(map[string]bool)
	}

	h.States[commit][testName] = testValue
}

func (h *History) Save() error {
	if err := fs.MkdirAll(appDir, 0700); err != nil {
		return err
	}

	out, err := fs.OpenFile(filepath.Join(appDir, "history.yaml"), os.O_CREATE|os.O_WRONLY, 0700)
	if err != nil {
		return err
	}

	b, err := yaml.Marshal(h)
	if err != nil {
		return err
	}

	if _, err := out.Write(b); err != nil {
		return err
	}

	return nil
}

func main() {
	file, err := afero.ReadFile(fs, "task.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}

	historyFile, err := afero.ReadFile(fs, "./.opendev/history.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}

	var repo History
	var tasks []Task

	if err = yaml.Unmarshal(file, &tasks); err != nil {
		fmt.Println(err)
		return
	}

	if err = yaml.Unmarshal(historyFile, &repo); err != nil {
		fmt.Println(err)
		return
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

	commit := head.Hash().String()

	if _, ok := repo.States[commit]; ok {
		fmt.Printf("Commit '%s' has already been tested.\n", commit)
		return
	}

	repo.Run(tasks...)

	repo.Tip = commit

	if err := repo.Save(); err != nil {
		fmt.Println("fatal:", err)
		return
	}
}
