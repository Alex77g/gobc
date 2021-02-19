package scm

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func Commit(msg string) {
	path, _ := GitRoot()
	r, _ := git.PlainOpen(path)

	w, _ := r.Worktree()

	status, _ := w.Status()

	log.Println(status)

	// extract user and mail from gitconf global
	u, _ := Username()
	e, _ := Email()

	commit, _ := w.Commit("example go-git commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  u,
			Email: e,
			When:  time.Now(),
		},
	})

	obj, _ := r.CommitObject(commit)

	fmt.Println(obj)
}

type ErrNotFound struct {
	Key string
}

func GitRoot() (string, error) {
	return execGitConfig("rev-parse", "--show-toplevel")
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("the key `%s` is not found", e.Key)
}

// Entire extracts configuration value from `$HOME/.gitconfig` file ,
// `$GIT_CONFIG`, /etc/gitconfig or include.path files.
func Entire(key string) (string, error) {
	return execGitConfig(key)
}

// Global extracts configuration value from `$HOME/.gitconfig` file or `$GIT_CONFIG`.
func Global(key string) (string, error) {
	return execGitConfig("--global", key)
}

// Local extracts configuration value from current project repository.
func Local(key string) (string, error) {
	return execGitConfig("--local", key)
}

// GithubUser extracts github.user name from `Entire gitconfig`
// This is same as Entire("github.user")
func GithubUser() (string, error) {
	return Entire("github.user")
}

// Username extracts git user name from `Entire gitconfig`.
// This is same as Entire("user.name")
func Username() (string, error) {
	return Entire("user.name")
}

// Email extracts git user email from `$HOME/.gitconfig` file or `$GIT_CONFIG`.
// This is same as Global("user.email")
func Email() (string, error) {
	return Entire("user.email")
}

// OriginURL extract remote origin url from current project repository.
// This is same as Local("remote.origin.url")
func OriginURL() (string, error) {
	return Local("remote.origin.url")
}

// Repository extract repository name of current project repository.
func Repository() (string, error) {
	url, err := OriginURL()
	if err != nil {
		return "", err
	}

	repo := retrieveRepoName(url)
	return repo, nil
}

// Github extracts github token from `Entire gitconfig`.
// This is same as Entire("github.token")
func GithubToken() (string, error) {
	return Entire("github.token")
}

func execGitConfig(args ...string) (string, error) {
	gitArgs := append([]string{"config", "--get", "--null"}, args...)
	var stdout bytes.Buffer
	cmd := exec.Command("git", gitArgs...)
	cmd.Stdout = &stdout
	cmd.Stderr = ioutil.Discard

	err := cmd.Run()
	if exitError, ok := err.(*exec.ExitError); ok {
		if waitStatus, ok := exitError.Sys().(syscall.WaitStatus); ok {
			if waitStatus.ExitStatus() == 1 {
				return "", &ErrNotFound{Key: args[len(args)-1]}
			}
		}
		return "", err
	}

	return strings.TrimRight(stdout.String(), "\000"), nil
}

var RepoNameRegexp = regexp.MustCompile(`.+/([^/]+)(\.git)?$`)

func retrieveRepoName(url string) string {
	matched := RepoNameRegexp.FindStringSubmatch(url)
	return strings.TrimSuffix(matched[1], ".git")
}
