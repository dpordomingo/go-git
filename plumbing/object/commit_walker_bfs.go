package object

import (
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type bfsCommitIterator struct {
	seenExternal map[plumbing.Hash]bool
	seen         map[plumbing.Hash]bool
	queue        []*Commit
}

// NewCommitIterBSF returns a CommitIter that walks the commit history,
// starting at the given commit and visiting its parents in pre-order.
// The given callback will be called for each visited commit. Each commit will
// be visited only once. If the callback returns an error, walking will stop
// and will return the error. Other errors might be returned if the history
// cannot be traversed (e.g. missing objects). Ignore allows to skip some
// commits from being iterated.
func NewCommitIterBSF(
	c *Commit,
	seenExternal map[plumbing.Hash]bool,
	ignore []plumbing.Hash,
) CommitIter {
	seen := make(map[plumbing.Hash]bool)
	for _, h := range ignore {
		seen[h] = true
	}

	var hasBeenSeen CommitFilter = func(c *Commit) bool {
		return seenExternal[c.Hash] || seen[c.Hash]
	}

	var notSeen CommitFilter = func(c *Commit) bool {
		return !hasBeenSeen(c)
	}

	return NewFilterCommitIter(c, &notSeen, &hasBeenSeen)
}
