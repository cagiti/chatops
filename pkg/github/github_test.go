package github

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetGitHubEventDetailsWithIssue(t *testing.T) {
	os.Setenv(githubEventPath, "testdata/issue_event.json")
	wantNumber := 8
	wantTime, err := getTimestampFromString("2020-05-01T12:04:11Z")
	assert.NoError(t, err)

	gotNumber, gotTime, err := GetGitHubEventDetails()
	assert.NoError(t, err)
	assert.Equal(t, wantTime, gotTime)
	assert.Equal(t, wantNumber, gotNumber)
}

func TestGetGitHubEventDetailsWithPR(t *testing.T) {
	os.Setenv(githubEventPath, "testdata/pr_event.json")
	wantNumber := 8
	wantTime, err := getTimestampFromString("2020-05-01T12:04:11Z")
	assert.NoError(t, err)

	gotNumber, gotTime, err := GetGitHubEventDetails()
	assert.NoError(t, err)
	assert.Equal(t, wantTime, gotTime)
	assert.Equal(t, wantNumber, gotNumber)
}

func TestGetGitHubEvent(t *testing.T) {
	wantEvent := &Event{
		Issue{
			Number:    8,
			UpdatedAt: "2020-05-01T12:04:11Z",
		},
		PullRequest{},
	}
	os.Setenv(githubEventPath, "testdata/issue_event.json")
	gotEvent, err := getGitHubEvent()

	assert.NoError(t, err)
	assert.Equal(t, wantEvent, gotEvent)
}

func TestGetIssuePRNumberWithIssue(t *testing.T) {
	event := Event{
		Issue: Issue{
			Number:    8,
			UpdatedAt: "2020-04-02T17:00:00Z",
		},
		PullRequest: PullRequest{},
	}
	number, err := getIssuePRNumber(&event)

	assert.NoError(t, err)
	assert.Equal(t, 8, number)
}

func TestGetIssuePRNumberWithEmptyIssueAndPR(t *testing.T) {
	event := Event{
		Issue:       Issue{},
		PullRequest: PullRequest{},
	}
	number, err := getIssuePRNumber(&event)

	assert.Error(t, err, "unable to determine issue/pr number from event")
	assert.Equal(t, 0, number)
}

func TestGetIssuePRNumberWithPR(t *testing.T) {
	event := Event{
		Issue: Issue{},
		PullRequest: PullRequest{
			Number:    8,
			UpdatedAt: "2020-04-02T17:00:00Z",
		},
	}
	number, err := getIssuePRNumber(&event)

	assert.NoError(t, err)
	assert.Equal(t, 8, number)
}

func TestGetIssuePRUpdatedAtWithIssue(t *testing.T) {
	updatedAt := "2020-04-02T17:00:00Z"
	wantUpdatedAt, err := getTimestampFromString(updatedAt)
	assert.NoError(t, err)
	event := Event{
		Issue: Issue{
			Number:    8,
			UpdatedAt: updatedAt,
		},
		PullRequest: PullRequest{},
	}
	gotUpdatedAt, err := getIssuePRUpdatedAtTime(&event)

	assert.NoError(t, err)
	assert.Equal(t, wantUpdatedAt, gotUpdatedAt)
}

func TestGetIssuePRUpdatedAtWithPR(t *testing.T) {
	updatedAt := "2020-04-02T17:00:00Z"
	wantUpdatedAt, err := getTimestampFromString(updatedAt)
	assert.NoError(t, err)
	event := Event{
		Issue: Issue{},
		PullRequest: PullRequest{
			Number:    8,
			UpdatedAt: updatedAt,
		},
	}
	gotUpdatedAt, err := getIssuePRUpdatedAtTime(&event)

	assert.NoError(t, err)
	assert.Equal(t, wantUpdatedAt, gotUpdatedAt)
}

func TestGetIssuePRUpdatedAtEmptyPRAndIssue(t *testing.T) {
	updatedAt := ""
	wantUpdatedAt, err := getTimestampFromString(updatedAt)
	assert.Error(t, err, "parsing the updatedAt time")
	event := Event{
		Issue: Issue{
			Number:    8,
			UpdatedAt: updatedAt,
		},
		PullRequest: PullRequest{},
	}
	gotUpdatedAt, err := getIssuePRUpdatedAtTime(&event)

	assert.Error(t, err, "parsing the updatedAt time")
	assert.Equal(t, wantUpdatedAt, gotUpdatedAt)
}

func TestGetTimestampFromString(t *testing.T) {
	updatedAt := "2020-04-02T17:00:00Z"
	gotUpdatedAt, err := getTimestampFromString(updatedAt)
	assert.NoError(t, err)
	assert.NotNil(t, gotUpdatedAt)
}

func TestGetTimestampFromStringUsingEmptyString(t *testing.T) {
	updatedAt := ""
	wantUpdatedAt := time.Time{}
	gotUpdatedAt, err := getTimestampFromString(updatedAt)
	assert.Error(t, err, "parsing the updatedAt time")
	assert.Equal(t, wantUpdatedAt, gotUpdatedAt)
}
