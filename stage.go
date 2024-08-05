package main

import (
	"fitter-services/compliance"
	"fitter-services/swearwords"
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/constructs-go/constructs/v10"
)

type DeploymentStageProps struct {
	awscdk.StageProps
	SwearwordsFileName string
}

func NewDeploymentStage(scope constructs.Construct, id string, props *DeploymentStageProps) awscdk.Stage {

	stage := awscdk.NewStage(scope, &id, &awscdk.StageProps{})

	swearwordsStackName := getStackName("swearwords")
	_, swLambdaName, swLambdaArn := swearwords.NewSwearwordsServiceStack(stage, swearwordsStackName, &swearwords.SwearwordsServiceStackProps{
		SwearwordsFileName: props.SwearwordsFileName,
	})

	complianceStackName := getStackName("compliance")
	compliance.NewComplianceServiceStack(stage, complianceStackName, &compliance.ComplianceServiceStackProps{
		Stage:                id,
		SwearwordsLambdaName: swLambdaName,
		SwearwordsLambdaArn:  swLambdaArn,
	})

	return stage
}

func getStackName(serviceName string) string {
	stackName := fmt.Sprintf("fitter-services-%s", serviceName)
	return stackName
}
