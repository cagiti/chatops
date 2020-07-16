package chatops

import (
	"context"

	"github.com/cagiti/chatops/pkg/loghelper"
	"github.com/cagiti/chatops/pkg/runner"
	"github.com/sethvargo/go-githubactions"

	"github.com/cagiti/chatops/pkg/github"
	"github.com/jenkins-x/jx-logging/pkg/log"
)

const (
	tokenValue = "token"
)

// Run the action
func Run() {
	// initialise logging
	loghelper.InitLogrus()
	ctx := context.Background()
	githubToken := githubactions.GetInput(tokenValue)
	if githubToken == "" {
		githubactions.Fatalf("missing input '%s'", tokenValue)
	}

	githubClient := github.NewClient(&ctx, githubToken)
	runner := runner.New(&ctx, githubClient)
	err := runner.Run()
	if err != nil {
		log.Logger().Error(err, "when running the runner!")
		githubactions.Fatalf("when running the runner!")
	}
}
