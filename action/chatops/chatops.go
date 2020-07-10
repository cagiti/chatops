package main

import (
	"context"

	"github.com/plumming/chatops-actions/pkg/loghelper"
	"github.com/plumming/chatops-actions/pkg/runner"
	"github.com/sethvargo/go-githubactions"

	"github.com/jenkins-x/jx-logging/pkg/log"
	"github.com/plumming/chatops-actions/pkg/github"
)

const (
	tokenValue = "token"
)

func main() {
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

	//githubactions.AddMask(fruit)
}
