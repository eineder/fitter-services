package main

import (
	"appsyncmasterclass-services/pipeline"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"
)

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	pipeline.NewPipelineStack(app, "appsyncmasterclass-pipelinestack", &pipeline.PipelineStackProps{})

	app.Synth(nil)
}
