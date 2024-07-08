package pipeline

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	codebuild "github.com/aws/aws-cdk-go/awscdk/v2/awscodebuild"
	pipeline "github.com/aws/aws-cdk-go/awscdk/v2/pipelines"
	"github.com/aws/jsii-runtime-go"

	"github.com/aws/constructs-go/constructs/v10"
)

type PipelineStackProps struct {
	awscdk.StackProps
	SwearwordsLambdaName string
	BranchName           string
}

func NewPipelineStack(scope constructs.Construct, id string, props *PipelineStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, jsii.String(id), &sprops)

	pipelineName := fmt.Sprintf("appsyncmasterclass_%s_pipeline", props.BranchName)
	githubRepo := pipeline.CodePipelineSource_GitHub(jsii.String("eineder/appsyncmasterclass-services"), &props.BranchName, &pipeline.GitHubSourceOptions{
		Authentication: awscdk.SecretValue_SecretsManager(jsii.String("github-token"), nil),
	})

	// self mutating pipeline
	myPipeline := pipeline.NewCodePipeline(stack, &pipelineName, &pipeline.CodePipelineProps{
		PipelineName: &pipelineName,
		// self mutation true - pipeline changes itself before application deployment
		SelfMutation: jsii.Bool(true),
		CodeBuildDefaults: &pipeline.CodeBuildOptions{
			BuildEnvironment: &codebuild.BuildEnvironment{
				// image version 6.0 recommended for newer go version
				BuildImage: codebuild.LinuxBuildImage_FromCodeBuildImageId(jsii.String("aws/codebuild/standard:7.0")),
			},
		},
		Synth: pipeline.NewCodeBuildStep(jsii.String("Synth"), &pipeline.CodeBuildStepProps{
			Input: githubRepo,
			Commands: &[]*string{
				jsii.String("npm install -g aws-cdk"),
				jsii.String("cdk synth"),
			},
		}),
	})

	testStage := myPipeline.AddStage(NewDeploymentStage(stack, "TEST", &MyStageProps{
		BranchName:         props.BranchName,
		SwearwordsFileName: "swearwords_test.txt",
	}), &pipeline.AddStageOpts{})
	testStage.AddPost(pipeline.NewCodeBuildStep(jsii.String("Test"), &pipeline.CodeBuildStepProps{
		Commands: &[]*string{
			jsii.String("go test ./..."),
		},
	}))

	if props.BranchName == "main" {
		myPipeline.AddStage(NewDeploymentStage(stack, "PROD", &MyStageProps{
			BranchName:         props.BranchName,
			SwearwordsFileName: "swearwords_prod.txt",
		}), &pipeline.AddStageOpts{})
	}

	return stack
}
