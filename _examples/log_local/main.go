package main

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/src-d/go-git.v4"
	. "gopkg.in/src-d/go-git.v4/_examples"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"
)

// Example of how to log
func main() {
	var path string
	var depth int
	var err error
	switch len(os.Args) {
	case 2:
		break
	case 3:
		depth, err = strconv.Atoi(os.Args[2])
		if err != nil {
			CheckIfError(fmt.Errorf("syntax: %s <path> [depth]\n\n'%s' is not a number", os.Args[0], os.Args[2]))
		}
	default:
		CheckIfError(fmt.Errorf("syntax: %s <path> [depth]", os.Args[0]))
	}

	path = os.Args[1]

	// ... opens the directory
	repo, err := git.PlainOpen(path)

	// ... retrieves the branch pointed by HEAD
	ref, err := repo.Head()
	CheckIfError(err)

	// ... retrieves the commit history
	cIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	CheckIfError(err)

	// ... just iterates over the commits, printing it
	if depth == 0 {
		err = cIter.ForEach(func(c *object.Commit) error {
			fmt.Println(c.Hash.String() /*, c.Message*/)
			return nil
		})
	} else {
		var count int
		err = cIter.ForEach(func(c *object.Commit) error {
			count++
			fmt.Println(c.Hash.String() /*, c.Message*/)
			if count < depth {
				return nil
			}

			return storer.ErrStop
		})
	}
	CheckIfError(err)
}
