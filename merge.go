package git

import (
	"errors"
	"fmt"

	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type strategy int

const (
	// Fail will cause the merge fail if a conflict is found
	Fail strategy = iota
	// Ours will use the changes from candidate if a conflict is found
	Ours
	// Theirs will use the changes from the base if a conflict is found
	Theirs
)

var (
	// ErrSameCommit returned when passed commits are the same
	ErrSameCommit = errors.New("passed commits are the same")
	// ErrAlreadyUpToDate returned when the target is behind base
	ErrAlreadyUpToDate = errors.New("target is behind base")
	// ErrHasConflicts returned when conflicts found
	ErrHasConflicts = errors.New("conflicts found")
	// ErrNoCommonHistory returned when no shared history
	ErrNoCommonHistory = errors.New("no shared history")
	// ErrNonFastForwardUpdate returned when no fast forward was possible
	//  defined at worktree.go
	// ErrWorktreeNotClean returned when no clean state in worktree
	//  defined at worktree.go

	// ExitCodeUnexpected returned when commit merge is required
	ErrNotImplementedNoFF = errors.New("no fast-forward merge is not implemented")
	// ErrNotImplementedUnrelated returned
	ErrNotImplementedUnrelated = errors.New("unrelated merge is not implemented")
	// ErrNotImplementedMessage returned
	ErrNotImplementedMessage = errors.New("custom message is not implemented")
)

// MergeOptions describes how a merge should be performed.
type MergeOptions struct {
	NoFF           bool   // NoFF when set to true, Merge will always create a merge commit
	FFOnly         bool   // FFOnly causes the Merge fail if it is not a fast forward
	AllowUnrelated bool   // AllowUnrelated performs the merge even with unrelated histories
	Message        string // Message text to be used for the message
}

// Merge merges the target commit over the base one, and moves `HEAD` to the merge.
// It returns the merge commit, and an error if the HEAD was not moved or
// when the merge operation could not be done.
func Merge(
	repo *Repository,
	base *object.Commit,
	target *object.Commit,
	options *MergeOptions,
) (*object.Commit, error) {
	if options == nil {
		options = &MergeOptions{}
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return nil, err
	}

	status, err := worktree.Status()
	if err != nil {
		return nil, err
	}

	for range status {
		return nil, ErrWorktreeNotClean
	}

	if base.Hash == target.Hash {
		return nil, ErrSameCommit
	}

	ancestors, err := MergeBase(base, target)
	if err != nil {
		return nil, err
	}

	if len(ancestors) == 0 {
		if options.AllowUnrelated {
			return merge(base, target, nil, options.Message)
		}

		return nil, ErrNoCommonHistory
	}

	for _, ancestor := range ancestors {
		if ancestor.Hash == base.Hash {
			if options.NoFF {
				// TODO(dpordomingo): there is a special case;
				// if asked with `--no-ff` it should be created an empty merge-commit.
				return nil, ErrNotImplementedNoFF
			}

			return target, nil
		}

		if ancestor.Hash == target.Hash {
			return nil, ErrAlreadyUpToDate
		}
	}

	mergeBase := ancestors[0]

	if options.FFOnly {
		return nil, ErrNonFastForwardUpdate
	}

	return merge(base, target, mergeBase, options.Message)
}

func merge(base, target, mergeBase *object.Commit, msg string) (*object.Commit, error) {

	if mergeBase == nil {
		// TODO(dpordomingo): handle --no-commit flag
		return nil, ErrNotImplementedUnrelated
	}

	var trees []*object.Tree
	for _, commit := range []*object.Commit{base, target} {
		tree, err := commit.Tree()
		if err != nil {
			return nil, err
		}

		trees = append(trees, tree)
	}

	changes, err := object.DiffTree(trees[0], trees[1])
	if err != nil {
		return nil, err
	}
	fmt.Println(changes)

	if msg != "" {
		// TODO(dpordomingo): handle -m option
		return nil, ErrNotImplementedMessage
	}

	// TODO(dpordomingo)
	return nil, ErrNotImplementedNoFF
}
