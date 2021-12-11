package main

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"os/exec"
)

// libgit's clone is very capricious so fuck it
func cloneRepo(url string) (func(), error) {
	if url == ""{
		return nil, errors.New("no repository provided")
	}
	_, err := exec.Command("git", "clone", url, "tmp").Output()
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
