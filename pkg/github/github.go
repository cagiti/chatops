package github

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/google/go-github/v32/github"
	"github.com/jenkins-x/jx-logging/pkg/log"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

var (
	githubRef        = os.Getenv("GITHUB_REF")
	githubEventName  = os.Getenv("GITHUB_EVENT_NAME")
	githubRepository = os.Getenv("GITHUB_REPOSITORY")
)

const (
	githubEventPath = "GITHUB_EVENT_PATH"
	forwardSlash    = "/"
)

// Issue entity.
type Issue struct {
	Number int `json:"number"`
}

// PullRequest entity.
type PullRequest struct {
	Number int `json:"number"`
}

// Event entity.
type Event struct {
	Issue       Issue       `json:"issue"`
	PullRequest PullRequest `json:"pull_request"`
}

// Repository entity
type Repository struct {
	Owner string
	Name  string
}

// NewClient returns a new github client.
func NewClient(ctx *context.Context, token string) *github.Client {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(*ctx, tokenSource)
	client := github.NewClient(httpClient)
	return client
}

// GetRepositoryInfo returns information about the repository
func GetRepositoryInfo() Repository {
	if strings.Count(githubRepository, forwardSlash) == 1 {
		r := strings.Split(githubRepository, forwardSlash)
		return Repository{Owner: r[0], Name: r[1]}
	}
	log.Logger().Warn("github repository string not formatted as expected")
	return Repository{}
}

// GetGitHubEvent using the GITHUB_EVENT_PATH environment variable.
func getGitHubEvent() (*Event, error) {
	log.Logger().Info("reading github event")
	eventPath := os.Getenv(githubEventPath)
	file, err := os.Open(eventPath)
	if err != nil {
		return nil, errors.Wrapf(err, "when attempting to read the event at: '%s'", eventPath)
	}
	defer file.Close()

	event := &Event{}
	err = json.NewDecoder(file).Decode(event)
	if err != nil {
		return nil, errors.Wrap(err, "when decoding json content")
	}
	return event, nil
}
func GetIssueOrPRNumber() (int, error) {
	number := 0
	event, err := getGitHubEvent()
	if err != nil {
		return number, errors.Wrap(err, "when getting getting github event")
	}

	number = event.Issue.Number
	if event.PullRequest.Number != 0 {
		number = event.PullRequest.Number
	}
	if number == 0 {
		return number, errors.New("unable to determine issue/pr number from event")
	}

	return number, nil
}
