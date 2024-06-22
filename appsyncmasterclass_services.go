package main

import (
	"appsyncmasterclass-services/compliance"
	"appsyncmasterclass-services/swearwords"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"
)

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	_, swearwordsLambdaName := swearwords.NewSwearwordsServiceStack(app, "SwearwordsServiceStack", &swearwords.SwearwordsServiceStackProps{})
	compliance.NewComplianceServiceStack(app, "ComplianceServiceStack", &compliance.ComplianceServiceStackProps{
		SwearwordsLambdaName: swearwordsLambdaName,
	})

	app.Synth(nil)
}
