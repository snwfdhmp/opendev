package main

import (
	"errors"
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

func gitHead() (string, error) {
	gitRepo, err := git.PlainOpen("./")
	if err != nil {
		return "", err
	}

	head, err := gitRepo.Head()
	if err != nil {
		return "", err
	}

	return head.Hash().String(), nil
}

func (h *History) Run(tasks ...Task) error {
	commit, err := gitHead()
	if err != nil {
		return err
	}

	if _, ok := h.States[commit]; ok {
		return errors.New(fmt.Sprintf("commit '%s' has already been tested", commit))
	}

	for _, t := range tasks {
		if err := exec.Command("/bin/zsh", "-c", t.Test).Run(); err != nil {
			h.Add(commit, t.Name, false)
			continue
		}
		h.Add(commit, t.Name, true)
	}

	h.Tip = commit

	return nil
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

func OpenHistory() (*History, error) {
	var h History

	file, err := afero.ReadFile(fs, "./.opendev/history.yaml")
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(file, &h); err != nil {
		return nil, err
	}

	return &h, nil
}

func ParseTasks(path string) (tasks []Task, err error) {
	file, err := afero.ReadFile(fs, path)
	if err != nil {
		return
	}

	if err = yaml.Unmarshal(file, &tasks); err != nil {
		return
	}

	return
}

func main() {

	repo, err := OpenHistory()
	tasks := ParseTasks("task.yaml")

	if err := repo.Run(tasks...); err != nil {
		fmt.Println("fatal:", err)
		return
	}

	if err := repo.Save(); err != nil {
		fmt.Println("fatal:", err)
		return
	}
}
