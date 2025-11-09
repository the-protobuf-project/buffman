// Package remote provides functionality for interacting with remote Git repositories,
// including operations like cloning and checking out specific versions.
package remote

import (
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
)

// PullOptions would be defined here or in a shared package, containing
// parameters for a pull operation. For context, it would look like this:
// type PullOptions struct {
//     Url    string
//     Out    string
//     Commit *string
// }

// github implements repository operations for standard Git repositories.
// Although named github, it uses generic Git commands and is not specific
// to GitHub's proprietary APIs.
type github struct{}

// Pull clones a Git repository from a URL into a specified output directory.
// If a commit hash is provided via opts.Commit, it performs a checkout to
// that specific commit after the clone is complete. If opts.Commit is nil,
// the repository will be left on the default branch (e.g., 'main' or 'master').
//
// It returns an error if the cloning or checkout operations fail.
func (g *github) Pull(opts PullOptions) error {
	// The PlainClone function returns the repository object on success.
	repo, err := git.PlainClone(opts.Out, &git.CloneOptions{
		URL: opts.Url,
	})
	if err != nil {
		return err
	}

	// If no specific commit is requested, the operation is complete.
	if opts.Commit == nil {
		return nil
	}

	// Get the worktree for the repository to perform a checkout.
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	// Checkout the specific commit hash provided in the options.
	return worktree.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(*opts.Commit),
	})
}
