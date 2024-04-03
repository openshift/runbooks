package main

import (
	"fmt"
	"path"
	"time"

	"github.com/go-git/go-git/v5"
)

func findLastCommitDate(repo *git.Repository, dirname, filename string, since *time.Time) (time.Time, error) {
	filepath := path.Join(dirname, filename)
	logOptions := &git.LogOptions{
		FileName: &filepath,
	}
	if since != nil {
		logOptions.Since = since
	}

	commitIter, err := repo.Log(logOptions)
	if err != nil {
		return time.Now(), fmt.Errorf("failed to get log: %w", err)
	}

	commit, err := commitIter.Next()
	if err != nil {
		return time.Now(), fmt.Errorf("failed to get next commit: %w", err)
	}

	lastCommitDate := commit.Committer.When.In(time.UTC)

	return lastCommitDate, nil
}
