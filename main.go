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
	Name   string
	Test   string
	Reward int
}

type History struct {
	States map[string]map[string]bool
	Tip    string
}

type RunReport struct {
	States  map[string]bool
	Rewards map[string]int
}

func NewRunReport() *RunReport {
	return &RunReport{make(map[string]bool), make(map[string]int)}
}

func (r *RunReport) Add(testName string, state bool, reward int) {
	if r.States == nil {
		r.States = make(map[string]bool)
	}
	if r.Rewards == nil {
		r.States = make(map[string]bool)
	}

	r.States[testName] = state
	r.Rewards[testName] = reward
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

// func (h *History) GetStates(commit string) (map[string]bool, bool) {
// 	s, ok := h.States[commit]
// 	return s, ok
// }

// func (h *History) GetLastState(testName string) (bool, bool) {
// 	s, ok := h.States[h.Tip][testName]
// 	return s, ok
// }

func (h *History) Run(tasks ...Task) (report RunReport, err error) {
	commit, err := gitHead()
	if err != nil {
		return
	}

	if _, ok := h.States[commit]; ok {
		err = errors.New(fmt.Sprintf("commit '%.6s' has already been tested", commit))
		return
	}

	for _, t := range tasks {
		var state bool
		if err := exec.Command("/bin/zsh", "-c", t.Test).Run(); err != nil {
			state = false
		} else {
			state = true
		}
		h.Add(commit, t.Name, state)
		lastState, ok := h.States[h.Tip][t.Name]
		if !ok {
			err = errors.New(fmt.Sprintf("cannot access %.6s/%s in history"))
			return
		}
		if lastState == state {
			continue
		}
		reward := t.Reward
		if lastState == true && state == false {
			reward = (-2) * t.Reward
		}
		report.Add(t.Name, state, reward)
	}

	h.Tip = commit

	return
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
	if err != nil {
		fmt.Println("fatal:", err)
		return
	}

	tasks, err := ParseTasks("task.yaml")
	if err != nil {
		fmt.Println("fatal:", err)
		return
	}

	if err := repo.Run(tasks...); err != nil {
		fmt.Println("fatal:", err)
		return
	}

	if err := repo.Save(); err != nil {
		fmt.Println("fatal:", err)
		return
	}
}
