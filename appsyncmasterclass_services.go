package main

import (
	"appsyncmasterclass-services/pipeline"
	"fmt"
	"os/exec"
	"strings"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"
)

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)
	branchName, err := getCurrentGitBranch()
	if err != nil {
		panic(err)
	}

	pipelineStackName := fmt.Sprintf("appsyncmasterclass-%s-pipelinestack", branchName)
	pipeline.NewPipelineStack(app, pipelineStackName, &pipeline.PipelineStackProps{
		BranchName: branchName,
	})

	app.Synth(nil)
}

func getCurrentGitBranch() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
