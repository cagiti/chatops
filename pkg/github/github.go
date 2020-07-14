package github

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v32/github"
	"github.com/jenkins-x/jx-logging/pkg/log"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

var (
	githubRepository = os.Getenv("GITHUB_REPOSITORY")
)

const (
	githubEventPath = "GITHUB_EVENT_PATH"
	forwardSlash    = "/"
	timestampFormat = "2006-01-02T15:04:05Z"
)

// Issue entity.
type Issue struct {
	Number    int    `json:"number"`
	UpdatedAt string `json:"updated_at"`
}

// PullRequest entity.
type PullRequest struct {
	Number    int    `json:"number"`
	UpdatedAt string `json:"updated_at"`
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

// GetGitHubEventDetails retrieves the events details
func GetGitHubEventDetails() (int, time.Time, error) {
	log.Logger().Info("getting github event details")
	event, err := getGitHubEvent()
	if err != nil {
		return 0, time.Time{}, errors.Wrap(err, "when getting github event")
	}
	log.Logger().Infof("Event: %v", event)
	number, err := getIssuePRNumber(event)
	if err != nil {
		return 0, time.Time{}, errors.Wrap(err, "when getting github issue/pr number")
	}

	updatedAt, err := getIssuePRUpdatedAtTime(event)
	if err != nil {
		return 0, time.Time{}, errors.Wrap(err, "when getting github issue/pr updatedAt time")
	}

	return number, updatedAt, nil
}
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
func getIssuePRNumber(event *Event) (int, error) {
	number := event.Issue.Number
	if event.PullRequest.Number != 0 {
		number = event.PullRequest.Number
	}
	if number == 0 {
		return number, errors.New("unable to determine issue/pr number from event")
	}
	return number, nil
}

func getIssuePRUpdatedAtTime(event *Event) (time.Time, error) {
	updatedAt := event.Issue.UpdatedAt
	if event.PullRequest.UpdatedAt != "" {
		updatedAt = event.PullRequest.UpdatedAt
	}
	if updatedAt == "" {
		return time.Time{}, errors.New("unable to determine issue/pr updatedAt time from event")
	}
	updatedAtTime, err := getTimestampFromString(updatedAt)
	if err != nil {
		return time.Time{}, errors.Wrapf(err, "when parsing time")
	}
	return updatedAtTime, nil
}

func getTimestampFromString(stringTime string) (time.Time, error) {

	timestamp, err := time.Parse(timestampFormat, stringTime)
	if err != nil {
		return time.Time{}, errors.Wrapf(err, "parsing the updatedAt time")
	}
	return timestamp, nil
}

//GetLastCommentForEvent attempts to get the last comment for the triggered event
func GetLastCommentForEvent(ctx context.Context, client *github.Client, number int, updatedAt time.Time) (string, error) {
	var comment string
	repo := GetRepositoryInfo()

	comments, _, err := client.Issues.ListComments(ctx, repo.Owner, repo.Name, number, nil)
	if err != nil {
		return "", errors.Wrapf(err, "when retrieving comments")
	}
	// get latest comment for the current event - looping through latest first
	for idx := len(comments) - 1; idx >= 0; idx-- {
		if comments[idx].GetCreatedAt().Before(updatedAt) ||
			comments[idx].GetCreatedAt().Equal(updatedAt) {
			comment = comments[idx].GetBody()
			break
		}
	}
	return comment, nil
}
