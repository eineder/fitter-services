package pipeline

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	codebuild "github.com/aws/aws-cdk-go/awscdk/v2/awscodebuild"
	pipeline "github.com/aws/aws-cdk-go/awscdk/v2/pipelines"
	"github.com/aws/jsii-runtime-go"

	"github.com/aws/constructs-go/constructs/v10"
)

type PipelineStackProps struct {
	awscdk.StackProps
	SwearwordsLambdaName string
}

func NewPipelineStack(scope constructs.Construct, id string, props *PipelineStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, jsii.String(id), &sprops)

	// GitHub repo with owner and repository name
	githubRepo := pipeline.CodePipelineSource_GitHub(jsii.String("eineder/appsyncmasterclass-services"), jsii.String("pipelines"), &pipeline.GitHubSourceOptions{
		Authentication: awscdk.SecretValue_SecretsManager(jsii.String("github-token"), nil),
	})

	// self mutating pipeline
	myPipeline := pipeline.NewCodePipeline(stack, jsii.String("cdkPipeline"), &pipeline.CodePipelineProps{
		PipelineName: jsii.String("CdkPipeline"),
		// self mutation true - pipeline changes itself before application deployment
		SelfMutation: jsii.Bool(true),
		CodeBuildDefaults: &pipeline.CodeBuildOptions{
			BuildEnvironment: &codebuild.BuildEnvironment{
				// image version 6.0 recommended for newer go version
				BuildImage: codebuild.LinuxBuildImage_FromCodeBuildImageId(jsii.String("aws/codebuild/standard:6.0")),
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

	myPipeline.AddStage(NewDeploymentStage(stack, "TEST", &MyStageProps{}), &pipeline.AddStageOpts{})
	myPipeline.AddStage(NewDeploymentStage(stack, "PROD", &MyStageProps{}), &pipeline.AddStageOpts{})

	return stack
}
