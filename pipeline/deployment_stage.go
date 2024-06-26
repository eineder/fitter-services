package pipeline

import (
	"appsyncmasterclass-services/compliance"
	"appsyncmasterclass-services/swearwords"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/constructs-go/constructs/v10"
)

type MyStageProps struct {
	awscdk.StageProps
}

func NewDeploymentStage(scope constructs.Construct, id string, props *MyStageProps) awscdk.Stage {

	stage := awscdk.NewStage(scope, &id, &awscdk.StageProps{})

	_, swLambdaName := swearwords.NewSwearwordsServiceStack(stage, "SwearwordsServiceStack", &swearwords.SwearwordsServiceStackProps{})
	compliance.NewComplianceServiceStack(stage, "ComplianceServiceStack", &compliance.ComplianceServiceStackProps{
		SwearwordsLambdaName: swLambdaName,
	})

	return stage
}
