package pipeline

import (
	"appsyncmasterclass-services/compliance"
	"appsyncmasterclass-services/swearwords"
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/constructs-go/constructs/v10"
)

type MyStageProps struct {
	awscdk.StageProps
	BranchName string
}

func NewDeploymentStage(scope constructs.Construct, id string, props *MyStageProps) awscdk.Stage {
	if props.BranchName == "" {
		panic("BranchName is required")
	}

	stage := awscdk.NewStage(scope, &id, &awscdk.StageProps{})

	swearwordsStackName := getStackName(props.BranchName, id, "swearwordsservice")
	_, swLambdaName := swearwords.NewSwearwordsServiceStack(stage, swearwordsStackName, &swearwords.SwearwordsServiceStackProps{})

	complianceStackName := getStackName(props.BranchName, id, "complianceservice")
	compliance.NewComplianceServiceStack(stage, complianceStackName, &compliance.ComplianceServiceStackProps{
		SwearwordsLambdaName: swLambdaName,
	})

	return stage
}

func getStackName(branchName string, stageId string, serviceName string) string {
	stackName := fmt.Sprintf("appsycmasterclass-%s-%s-%s", branchName, stageId, serviceName)
	return stackName
}
