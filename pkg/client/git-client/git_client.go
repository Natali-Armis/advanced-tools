package git_client

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type GitClient struct{}

func GetGitClient() *GitClient {
	return &GitClient{}
}

func (g *GitClient) runGitCommand(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out.String())
	}
	return out.String(), nil
}

func (g *GitClient) GitPull() (string, error) {
	return g.runGitCommand("pull")
}

func (g *GitClient) GitAdd(files string) (string, error) {
	return g.runGitCommand("add", files)
}

func (g *GitClient) GitCommit(message string) (string, error) {
	return g.runGitCommand("commit", "-m", message)
}

func (g *GitClient) GitCheckoutNewBranch(branchName string) (string, error) {
	return g.runGitCommand("checkout", "-b", branchName)
}

func (g *GitClient) GitPush(branch string) (string, error) {
	return g.runGitCommand("push", "origin", branch)
}

func (g *GitClient) CreatePullRequest(owner, repo, title, head, base, token string) (*github.PullRequest, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	newPR := &github.NewPullRequest{
		Title: github.String(title),
		Head:  github.String(head),
		Base:  github.String(base),
	}

	pr, _, err := client.PullRequests.Create(ctx, owner, repo, newPR)
	if err != nil {
		return nil, err
	}
	return pr, nil
}
