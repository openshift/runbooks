package main

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func getDirCurrentTree(repo *git.Repository, dir string) (*object.Tree, error) {
	ref, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD: %w", err)
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get worktree: %w", err)
	}

	commitTree, err := commit.Tree()
	if err != nil {
		return nil, fmt.Errorf("failed to get commit tree: %w", err)
	}

	runbooksTree, err := commitTree.Tree(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to get runbooks tree: %w", err)
	}

	return runbooksTree, nil
}

func newBranchFromMain(repo *git.Repository, name string) (*git.Worktree, error) {
	w, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("failed to get worktree: %w", err)
	}

	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(downstreamMainBranch),
		Create: false,
		Keep:   false,
		Force:  true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to checkout to main: %w", err)
	}

	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(name),
		Create: true,
		Keep:   false,
		Force:  false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create branch: %w", err)
	}

	return w, nil
}
