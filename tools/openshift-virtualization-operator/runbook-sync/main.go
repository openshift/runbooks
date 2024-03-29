package main

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/google/go-github/v60/github"
	"k8s.io/klog/v2"
)

const (
	githubUsername = "openshift-merge-bot[bot]"
	githubEmail    = "148852131+openshift-merge-bot[bot]@users.noreply.github.com"

	upstreamCloneDir      = "/tmp/kubevirt-monitoring"
	upstreamRepositoryURL = "github.com/kubevirt/monitoring"
	upstreamRunbooksDir   = "docs/runbooks"

	downstreamMainBranch      = "master"
	downstreamCloneDir        = "/tmp/runbooks"
	downstreamRepositoryOwner = "openshift"
	downstreamRepositoryName  = "runbooks"
	downstreamRunbooksDir     = "alerts/openshift-virtualization-operator"

	customRemoteName = "tokenized"
)

var (
	downstreamRepositoryURL = fmt.Sprintf("github.com/%s/%s", downstreamRepositoryOwner, downstreamRepositoryName)
)

type runbookSyncArgs struct {
	githubToken string
	dryRun      bool
}

type runbookSync struct {
	ghClient       *github.Client
	downstreamRepo *git.Repository
	dryRun         bool
}

func main() {
	rbSyncArgs := getRunbookSyncArgs()

	downstreamRepo, upstreamRepo := setup(rbSyncArgs.githubToken)
	runbooksToUpdate, runbooksToDelete := listRunbooksThatNeedUpdate(downstreamRepo, upstreamRepo)

	for _, r := range runbooksToUpdate {
		klog.Infof("runbook %s needs update. Last update: %s, upstream last update: %s", r.name, r.lastLocalUpdate, r.upstreamLastUpdated)
	}

	rbSync := &runbookSync{
		ghClient:       github.NewClient(nil).WithAuthToken(rbSyncArgs.githubToken),
		downstreamRepo: downstreamRepo,
		dryRun:         rbSyncArgs.dryRun,
	}

	rbSync.createRunbooksBranches(runbooksToUpdate, runbooksToDelete)
}

func getRunbookSyncArgs() runbookSyncArgs {
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		klog.Fatal("GITHUB_TOKEN environment variable is required")
	}

	dryRun := os.Getenv("DRY_RUN")
	if dryRun == "" {
		dryRun = "true"
	}
	if dryRun != "true" && dryRun != "false" {
		klog.Fatal("DRY_RUN environment variable must be 'true' or 'false'")
	}
	klog.Infof("dry run: %s", dryRun)

	return runbookSyncArgs{
		githubToken: githubToken,
		dryRun:      dryRun != "false",
	}
}

func (rbSync *runbookSync) createRunbooksBranches(runbooksToUpdate []runbook, runbooksToDelete []runbook) {
	if len(runbooksToUpdate) == 0 {
		klog.Info("no runbooks to update")
	}

	for _, rb := range runbooksToUpdate {
		klog.Infof("---")
		_ = rbSync.updateRunbook(rb)
	}

	// TODO: @machadovilaca - reenable
	// Disable delete/deprecations for now
	// if len(runbooksToDelete) == 0 {
	// 	klog.Info("no runbooks to delete")
	// }

	// for _, rb := range runbooksToDelete {
	// 	klog.Infof("---")
	// 	_ = rbSync.deleteRunbook(rb)
	//
	// }
}

func (rbSync *runbookSync) updateRunbook(rb runbook) string {
	lastUpdateDate := rb.upstreamLastUpdated.Format("20060102150405")
	runbookName := strings.Replace(rb.name, ".md", "", -1)
	branchName := fmt.Sprintf("cnv-runbook-sync-%s/%s", lastUpdateDate, runbookName)

	prExists, pr, err := rbSync.prForBranchPreviouslyCreated(branchName)
	if err != nil {
		klog.Fatalf("failed to check if branch exists: %v", err)
	}

	if prExists {
		klog.Infof("PR for '%s' was previously created: %s", branchName, pr.GetHTMLURL())
		return branchName
	}

	worktree, err := newBranchFromMain(rbSync.downstreamRepo, branchName)
	if err != nil {
		klog.Fatalf("failed to create branch: %v", err)
	}

	err = copyRunbook(rb.name)
	if err != nil {
		klog.Fatalf("failed to copy file: %v", err)
	}

	commitMessage := fmt.Sprintf("Sync CNV runbook %s (Updated at %s)", rb.name, rb.upstreamLastUpdated)

	err = rbSync.commitAndPush(worktree, commitMessage)
	if err != nil {
		klog.Fatalf("failed to push changes: %v", err)
	}

	body := fmt.Sprintf(
		"This is an automated PR by 'tools/openshift-virtualization-operator/runbook-sync'.\n\n"+
			"CNV runbook '%s' was updated in upstream https://%s at %s.\n"+
			"This PR syncs the runbook in this repository to contain all new added changes.\n\n"+
			"/cc @machadovilaca",
		rb.name, upstreamRepositoryURL, rb.upstreamLastUpdated,
	)

	err = rbSync.createPR(branchName, commitMessage, body)
	if err != nil {
		klog.Fatalf("failed to create PR: %v", err)
	}

	return branchName
}

func (rbSync *runbookSync) deleteRunbook(rb runbook) string {
	runbookName := strings.Replace(rb.name, ".md", "", -1)
	branchName := fmt.Sprintf("cnv-runbook-delete-%s", runbookName)

	prExists, pr, err := rbSync.prForBranchPreviouslyCreated(branchName)
	if err != nil {
		klog.Fatalf("failed to check if branch exists: %v", err)
	}

	if prExists {
		klog.Infof("PR for '%s' was previously created: %s", branchName, pr.GetHTMLURL())
		return branchName
	}

	worktree, err := newBranchFromMain(rbSync.downstreamRepo, branchName)
	if err != nil {
		klog.Fatalf("failed to create branch: %v", err)
	}

	klog.Infof("deleting file %s", rb.name)
	err = os.Remove(path.Join(downstreamCloneDir, downstreamRunbooksDir, rb.name))
	if err != nil {
		klog.Fatalf("failed to delete file: %v", err)
	}

	commitMessage := fmt.Sprintf("Delete CNV runbook %s", rb.name)

	err = rbSync.commitAndPush(worktree, commitMessage)
	if err != nil {
		klog.Fatalf("failed to push changes: %v", err)
	}

	body := fmt.Sprintf(
		"This is an automated PR by 'tools/openshift-virtualization-operator/runbook-sync'.\n\n"+
			"CNV runbook '%s' was deleted in upstream https://%s.\n"+
			"This PR deletes the runbook in this repository.\n\n"+
			"/cc @machadovilaca",
		rb.name, upstreamRepositoryURL,
	)

	err = rbSync.createPR(branchName, commitMessage, body)
	if err != nil {
		klog.Fatalf("failed to create PR: %v", err)
	}

	return branchName
}

func (rbSync *runbookSync) commitAndPush(worktree *git.Worktree, msg string) error {
	_, err := worktree.Add(downstreamRunbooksDir)
	if err != nil {
		return fmt.Errorf("failed to add changes: %w", err)
	}

	_, err = worktree.Commit(msg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  githubUsername,
			Email: githubEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to commit changes: %w", err)
	}
	klog.Infof("successfully committed: %s", msg)

	if rbSync.dryRun {
		klog.Warning("[DRY RUN] skipping push")
		return nil
	}

	err = rbSync.downstreamRepo.Push(&git.PushOptions{
		RemoteName: customRemoteName,
	})
	if err != nil {
		return fmt.Errorf("failed to push changes: %w", err)
	}
	klog.Info("successfully pushed changes")

	return nil
}

func (rbSync *runbookSync) prForBranchPreviouslyCreated(branchName string) (bool, *github.PullRequest, error) {
	prs, _, err := rbSync.ghClient.PullRequests.List(context.Background(), downstreamRepositoryOwner, downstreamRepositoryName, &github.PullRequestListOptions{
		State: "all",
		Head:  fmt.Sprintf("%s:%s", downstreamRepositoryOwner, branchName),
	})
	if err != nil {
		return false, nil, err
	}

	if len(prs) == 0 {
		return false, nil, nil
	}

	return true, prs[0], nil
}

func (rbSync *runbookSync) createPR(branchName string, title string, body string) error {
	headBranch := fmt.Sprintf("%s:%s", downstreamRepositoryOwner, branchName)
	baseBranch := downstreamMainBranch

	prOpts := &github.NewPullRequest{
		Title: &title,
		Head:  &headBranch,
		Base:  &baseBranch,
		Body:  &body,
	}

	if rbSync.dryRun {
		klog.Warningf("[DRY RUN] skipping PR creation '%s', %s => %s", *prOpts.Title, *prOpts.Head, *prOpts.Base)
		return nil
	}

	pr, _, err := rbSync.ghClient.PullRequests.Create(context.Background(), downstreamRepositoryOwner, downstreamRepositoryName, prOpts)
	if err != nil {
		return err
	}

	klog.Infof("PR created: %s", pr.GetHTMLURL())

	return nil
}
