package swearwords

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3deployment"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

const SWEARWORDS_TABLE_NAME string = "SWEARWORDS_TABLE_NAME"
const BUCKET_NAME string = "BUCKET_NAME"
const BUCKET_KEY string = "BUCKET_KEY"

type SwearwordsServiceStackProps struct {
	awscdk.StackProps
	SwearwordsFileName string
}

func NewSwearwordsServiceStack(scope constructs.Construct, id string, props *SwearwordsServiceStackProps) (stack awscdk.Stack, lambdaName string) {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack = awscdk.NewStack(scope, &id, &sprops)

	bucket := awss3.NewBucket(stack, jsii.String("swearwords-bucket"), &awss3.BucketProps{})
	cdkAsset := awss3deployment.Source_Asset(jsii.String("swearwords/assets"), &awss3assets.AssetOptions{})
	awss3deployment.NewBucketDeployment(
		stack,
		jsii.String("swearwords-bucket-deployment"),
		&awss3deployment.BucketDeploymentProps{
			Sources:           &[]awss3deployment.ISource{cdkAsset},
			DestinationBucket: bucket,
		})

	table := awsdynamodb.NewTable(stack, jsii.String("swearwords-table"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("word"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		BillingMode: awsdynamodb.BillingMode_PAY_PER_REQUEST,
	})

	fn := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("contains-swearwords-lambda"), &awscdklambdagoalpha.GoFunctionProps{
		Entry:   jsii.String("swearwords/contains_swearwords/contains_swearwords.go"),
		Runtime: awslambda.Runtime_PROVIDED_AL2(),
		Environment: &map[string]*string{
			SWEARWORDS_TABLE_NAME: table.TableName()},
	})
	fn.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Effect: awsiam.Effect_ALLOW,
		Actions: jsii.Strings(
			"dynamodb:GetItem"),
		Resources: &[]*string{table.TableArn()},
	}))

	primeFn := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("prime-swearwords-lambda"), &awscdklambdagoalpha.GoFunctionProps{
		Entry:   jsii.String("swearwords/prime_swearwords/prime_swearwords.go"),
		Runtime: awslambda.Runtime_PROVIDED_AL2(),
		Environment: &map[string]*string{
			SWEARWORDS_TABLE_NAME: table.TableName(),
			BUCKET_NAME:           bucket.BucketName(),
			BUCKET_KEY:            &props.SwearwordsFileName},
	})
	primeFn.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Effect: awsiam.Effect_ALLOW,
		Actions: jsii.Strings(
			"dynamodb:BatchWriteItem",
			"dynamodb:GetItem",
			"dynamodb:PutItem",
			"dynamodb:ImportTable",
			"dynamodb:DescribeImport",
			"dynamodb:ListImports"),
		Resources: &[]*string{table.TableArn()},
	}))
	primeFn.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Effect: awsiam.Effect_ALLOW,
		Actions: jsii.Strings("s3:GetObject",
			"s3:ListBucket"),
		Resources: jsii.Strings(*bucket.BucketArn(), *bucket.BucketArn()+"/*"),
	}))
	primeFn.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Effect: awsiam.Effect_ALLOW,
		Actions: jsii.Strings("logs:CreateLogGroup",
			"logs:CreateLogStream",
			"logs:DescribeLogGroups",
			"logs:DescribeLogStreams",
			"logs:PutLogEvents",
			"logs:PutRetentionPolicy"),
		Resources: jsii.Strings("*"),
	}))

	awscdk.NewCfnOutput(stack, jsii.String("BucketNameOutput"), &awscdk.CfnOutputProps{
		Key:   jsii.String("BucketName"),
		Value: bucket.BucketName(),
	})
	awscdk.NewCfnOutput(stack, jsii.String("BucketKeyOutput"), &awscdk.CfnOutputProps{
		Key:   jsii.String("BucketKey"),
		Value: &props.SwearwordsFileName,
	})
	awscdk.NewCfnOutput(stack, jsii.String("SwearwordsTableNameOutput"), &awscdk.CfnOutputProps{
		Key:   jsii.String("SwearwordsTableName"),
		Value: table.TableName(),
	})
	awscdk.NewCfnOutput(stack, jsii.String("PrimeSwearwordsFunctionNameOutput"), &awscdk.CfnOutputProps{
		Key:   jsii.String("PrimeSwearwordsFunctionName"),
		Value: primeFn.FunctionName(),
	})

	return stack, *fn.FunctionName()
}
