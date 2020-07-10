package runner

import (
	"context"
	"strings"

	"github.com/pkg/errors"

	gh "github.com/google/go-github/v32/github"
	"github.com/jenkins-x/jx-logging/pkg/log"
	"github.com/plumming/chatops-actions/pkg/github"
)

// Runner struct containing context and the github client.
type Runner struct {
	ctx      *context.Context
	ghClient *gh.Client
}

// New creates a new Runner.
func New(ctx *context.Context, ghClient *gh.Client) Runner {
	return Runner{
		ctx:      ctx,
		ghClient: ghClient,
	}
}

func (runner *Runner) Run() error {

	number, err := github.GetIssueOrPRNumber()
	if err != nil {
		log.Logger().Error(err, "error whilst retrieving issue/pr number")
	}

	repo := github.GetRepositoryInfo()
	comments, _, err := runner.ghClient.Issues.ListComments(*runner.ctx, repo.Owner, repo.Name, number, nil)
	if err != nil {
		return err
	}

	for _, comment := range comments {
		if strings.Contains(strings.ToLower(comment.GetBody()), "/woof") {
			//TODO: send dog image
			prComment := &gh.PullRequestComment{Body: gh.String("woof woof woof")}
			_, resp, err := runner.ghClient.PullRequests.CreateComment(*runner.ctx, repo.Owner, repo.Name, number, prComment)
			defer resp.Body.Close()
			if err != nil {
				return errors.Wrapf(err, "when posting a comment to #%d", number)
			}
		}
	}
	return nil
}
