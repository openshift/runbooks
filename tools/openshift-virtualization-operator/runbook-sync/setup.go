package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"k8s.io/klog/v2"
)

func setup(githubToken string) (*git.Repository, *git.Repository) {
	downstreamRepo, err := deleteAndClone(downstreamRepositoryURL, downstreamCloneDir)
	if err != nil {
		klog.Fatal(fmt.Errorf("failed to clone or update repository: %w", err))
	}
	addRemoteWithTokenToLocalRepo(downstreamRepo, githubToken)

	upstreamRepo, err := deleteAndClone(upstreamRepositoryURL, upstreamCloneDir)
	if err != nil {
		klog.Fatal(fmt.Errorf("failed to clone or update repository: %w", err))
	}

	return downstreamRepo, upstreamRepo
}

func deleteAndClone(repoUrl string, dest string) (*git.Repository, error) {
	err := os.RemoveAll(dest)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		klog.Fatal(fmt.Errorf("failed to delete directory: %w", err))
	}

	url := "https://" + repoUrl

	klog.Info("cloning repository ", url, " to ", dest)
	repo, err := git.PlainClone(dest, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to clone repository: %w", err)
	}

	return repo, nil
}

func addRemoteWithTokenToLocalRepo(repo *git.Repository, githubToken string) {
	remote, err := repo.Remote(customRemoteName)
	if err != nil && !errors.Is(err, git.ErrRemoteNotFound) {
		klog.Fatal(fmt.Errorf("failed to get remote: %w", err))
	} else if err == nil {
		klog.Info("remote with token already exists, deleting")
		deleteErr := repo.DeleteRemote(remote.Config().Name)
		if deleteErr != nil {
			klog.Fatal(fmt.Errorf("failed to delete remote: %w", err))
		}
	}

	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: customRemoteName,
		URLs: []string{fmt.Sprintf("https://oauth2:%s@%s", githubToken, downstreamRepositoryURL)},
	})
	if err != nil {
		klog.Fatal(fmt.Errorf("failed to create remote: %w", err))
	}
}
