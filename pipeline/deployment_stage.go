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
	SwearwordsFileName string
}

func NewDeploymentStage(scope constructs.Construct, id string, props *MyStageProps) awscdk.Stage {

	stage := awscdk.NewStage(scope, &id, &awscdk.StageProps{})

	swearwordsStackName := getStackName(id, "swearwordsservice")
	_, swLambdaName := swearwords.NewSwearwordsServiceStack(stage, swearwordsStackName, &swearwords.SwearwordsServiceStackProps{
		SwearwordsFileName: props.SwearwordsFileName,
	})

	complianceStackName := getStackName(id, "complianceservice")
	compliance.NewComplianceServiceStack(stage, complianceStackName, &compliance.ComplianceServiceStackProps{
		SwearwordsLambdaName: swLambdaName,
	})

	return stage
}

func getStackName(stageId string, serviceName string) string {
	stackName := fmt.Sprintf("appsycmasterclass-%s-%s", stageId, serviceName)
	return stackName
}
