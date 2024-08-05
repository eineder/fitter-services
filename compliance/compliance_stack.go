package compliance

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsevents"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseventstargets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type ComplianceServiceStackProps struct {
	awscdk.StackProps
	Stage                string
	SwearwordsLambdaName string
	SwearwordsLambdaArn  string
}

func NewComplianceServiceStack(scope constructs.Construct, id string, props *ComplianceServiceStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	fn := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("handler-on-new-tweet-posted"), &awscdklambdagoalpha.GoFunctionProps{
		Entry:   jsii.String("compliance/lambda/on_tweet_posted.go"),
		Runtime: awslambda.Runtime_PROVIDED_AL2(),
		Environment: &map[string]*string{
			"SWEARWORDS_LAMBDA_NAME": jsii.String(props.SwearwordsLambdaName),
		},
	})
	fn.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("lambda:InvokeFunction"),
		Resources: jsii.Strings(props.SwearwordsLambdaArn),
	}))
	fn.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("dynamodb:UpdateItem"),
		Resources: jsii.Strings("arn:aws:dynamodb:*:*:table/*TweetsTable*"),
	}))

	rule := awsevents.NewRule(stack, jsii.String("rule-on-new-tweet-posted"), &awsevents.RuleProps{
		EventPattern: &awsevents.EventPattern{
			DetailType: jsii.Strings(props.Stage + "_new_tweet_posted"),
		},
	})

	target := awseventstargets.NewLambdaFunction(fn, &awseventstargets.LambdaFunctionProps{})
	rule.AddTarget(target)

	return stack
}
