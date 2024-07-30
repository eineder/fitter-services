package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"
)

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewDeploymentStage(app, "dev", &DeploymentStageProps{
		SwearwordsFileName: "swearwords_dev.txt",
	})

	NewDeploymentStage(app, "prod", &DeploymentStageProps{
		SwearwordsFileName: "swearwords_prod.txt",
	})

	app.Synth(nil)
}
