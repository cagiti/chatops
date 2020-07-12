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

	number, updatedAt, err := github.GetGitHubEventDetails()
	if err != nil {
		log.Logger().Error(err, "error whilst retrieving event details")
		return err
	}

	comment, err := github.GetLastCommentForEvent(*runner.ctx, runner.ghClient, number, updatedAt)
	if err != nil {
		return err
	}

	if strings.Contains(strings.ToLower(comment), "/bark") {
		repo := github.GetRepositoryInfo()
		log.Logger().Info("contains a bark")
		issueComment := &gh.IssueComment{Body: gh.String("woof woof")}
		_, resp, err := runner.ghClient.Issues.CreateComment(*runner.ctx, repo.Owner, repo.Name, number, issueComment)
		defer resp.Body.Close()
		if err != nil {
			return errors.Wrapf(err, "when posting a comment to #%d", number)
		}
	}

	return nil
}
