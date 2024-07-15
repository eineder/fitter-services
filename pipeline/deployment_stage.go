package pipeline

import (
	"appsyncmasterclass-services/compliance"
	"appsyncmasterclass-services/swearwords"
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

	swearwordsStackName := getStackName("swearwordsservice")
	_, swLambdaName := swearwords.NewSwearwordsServiceStack(stage, swearwordsStackName, &swearwords.SwearwordsServiceStackProps{
		SwearwordsFileName: props.SwearwordsFileName,
	})

	complianceStackName := getStackName("complianceservice")
	compliance.NewComplianceServiceStack(stage, complianceStackName, &compliance.ComplianceServiceStackProps{
		SwearwordsLambdaName: swLambdaName,
	})

	return stage
}

func getStackName(serviceName string) string {
	stackName := fmt.Sprintf("appsycmasterclass-%s", serviceName)
	return stackName
}
