package main

import (
	"strings"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"os/exec"
)

var (
	ErrRepoNotFound = fmt.Errorf("repo not found")
)

// libgit's clone is very capricious so fuck it
func cloneRepo(url string) (func(), error) {
	if url == ""{
		return nil, errors.New("no repository provided")
	}
	out, err := exec.Command("git", "clone", url, "tmp").CombinedOutput()
	if strings.Contains(string(out), "Repository not found") {
		return nil, errors.Wrapf(ErrRepoNotFound, "url is %s", url)
	}
	if err != nil {
		return nil, err
	}
	return func() {
		err := os.RemoveAll("./tmp")
		if err != nil {
			fmt.Printf("Could not remove ./tmp: %+v\n", err)
		}
	}, nil
}
