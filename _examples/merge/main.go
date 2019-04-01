package main

import (
	"fmt"
	"os"

	"gopkg.in/src-d/go-git.v4"
	utils "gopkg.in/src-d/go-git.v4/_examples"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type exitCode int

const (
	exitCodeNoFastForward         utils.ExitCode = 1
	exitCodeFeatureNotImplemented utils.ExitCode = 9
)

const (
	cmdDesc = "Performs the merge between two commits, and moves HEAD into the merge commit:"

	helpShortMsg = `
  usage: %_COMMAND_NAME_% <path> <baseCommitRev> <commitRev> [-m <msg>] [--ff-only]
                           [--no-ff] [--no-commit] [--allow-unrelated-histories]
     or: %_COMMAND_NAME_% --help

 params:
    <path>           Path to the git repository
    <baseCommitRev>  Git revision of the commit that will be the base of the merge
    <commitRev>      Git revision of the commit that will be merged over the base

options:
    (no options)   Performs the regular merge, as fast-forward if possible
    --ff-only      If the merge is not a fast-forward the process exits with 1
    --help         Show the full help message of %_COMMAND_NAME_%

[NOT IMPLEMENTED]
    -m <msg>       Uses the passed <msg> for the merge commit message
    --no-ff        Create a merge commit even when it could be a fast-forward
    --no-commit    Performs the merge in the worktree and do not create a commit
    --allow-unrelated-histories
                   Lets the merge to operate with commits not sharing its history
`
)

func main() {
	if len(os.Args) == 1 {
		utils.HelpAndExit(cmdDesc, helpShortMsg)
	}

	if os.Args[1] == "--help" || os.Args[1] == "-h" {
		utils.HelpAndExit(cmdDesc, helpLongMsg)
	}

	if len(os.Args) < 4 {
		utils.WrongSyntaxAndExit(helpShortMsg)
	}

	path := os.Args[1]
	commitRevs := os.Args[2:4]

	mergeOptions := &git.MergeOptions{}
	if len(os.Args) > 4 {
		args := os.Args[4:]
		skipNext := false
		for i, v := range args {
			if skipNext {
				skipNext = false
				continue
			}
			switch v {
			case "--no-ff":
				mergeOptions.NoFF = true
			case "--ff-only":
				mergeOptions.FFOnly = true
			case "--no-commit":
				mergeOptions.NoCommit = true
			case "--allow-unrelated-histories":
				mergeOptions.AllowUnrelated = true
			case "-m":
				if len(args) < i+2 {
					utils.WrongSyntaxAndExit("message not defined, should be -m <msg>")
				}

				mergeOptions.Message = args[i+1]
				skipNext = true
			default:
				utils.WrongSyntaxAndExit(helpShortMsg)
			}
		}
	}

	// Open a git repository from current directory
	repo, err := git.PlainOpen(path)
	utils.ExitIfError(err, utils.ExitCodeCouldNotOpenRepository)

	// Get the hashes of the passed revisions
	var hashes []*plumbing.Hash
	for _, rev := range commitRevs {
		hash, err := repo.ResolveRevision(plumbing.Revision(rev))
		utils.ExitIfError(err, utils.ExitCodeCouldNotParseRevision, rev)
		hashes = append(hashes, hash)
	}

	// Get the commits identified by the passed hashes
	var commits []*object.Commit
	for _, hash := range hashes {
		commit, err := repo.CommitObject(*hash)
		utils.ExitIfError(err, utils.ExitCodeWrongCommitHash, hash.String())
		commits = append(commits, commit)
	}

	commit, err := git.Merge(repo, commits[0], commits[1], mergeOptions)
	switch err {
	case git.ErrSameCommit:
		fmt.Println("Already up to date. Both are the same commit.")
		return
	case git.ErrAlreadyUpToDate:
		fmt.Println("Already up to date.")
		return
	case git.ErrNotImplementedUnrelated,
		git.ErrNotImplementedNoCommit,
		git.ErrNotImplementedNoFF,
		git.ErrNotImplementedMessage:
		utils.ExitIfError(err, exitCodeFeatureNotImplemented, "unimplemented feature")
	}

	utils.ExitIfError(err, utils.ExitCodeExpected, "merge failed")

	utils.PrintCommits(commit)
}
