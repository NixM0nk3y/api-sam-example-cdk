package hosting

import (
	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/awsapigatewayv2integrations"
	"github.com/aws/aws-cdk-go/awscdk/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/awslogs"
	"github.com/aws/aws-cdk-go/awscdk/awss3assets"
	"github.com/aws/jsii-runtime-go"

	"github.com/aws/constructs-go/constructs/v3"
)

type LocalBundling struct {
}

// needed to allow sam local testing
func (c LocalBundling) TryBundle(outputDir *string, options *awscdk.BundlingOptions) *bool {
	return jsii.Bool(false)
}

type HostingProps struct {
	Tenant           string                  ``
	Environment      string                  ``
	Appplication     string                  ``
	NestedStackProps awscdk.NestedStackProps ``
}

func HostingStack(scope constructs.Construct, id string, props *HostingProps) awscdk.Construct {

	construct := awscdk.NewConstruct(scope, &id)

	// Go build options
	bundlingOptions := &awscdk.BundlingOptions{
		Image:       awslambda.Runtime_PROVIDED_AL2().BundlingDockerImage(),
		OutputType:  awscdk.BundlingOutput_ARCHIVED,
		Environment: &map[string]*string{},
		Command: jsii.Strings(
			"bash", "-c",
			"make lambda/package",
		),
		Local: &LocalBundling{},
	}

	// webhook lambda
	apiLambda := awslambda.NewFunction(construct, jsii.String("Lambda"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2(),
		Handler: jsii.String("app.handler"),
		Code: awslambda.AssetCode_FromAsset(
			jsii.String("./resources/api"),
			&awss3assets.AssetOptions{
				Bundling: bundlingOptions,
			},
		),
		Tracing:      awslambda.Tracing_ACTIVE,
		LogRetention: awslogs.RetentionDays_ONE_WEEK,
		Architectures: &[]awslambda.Architecture{
			awslambda.Architecture_X86_64(),
		},
		Environment: &map[string]*string{
			"LOG_LEVEL": jsii.String("DEBUG"),
		},
	})

	//
	httpapi := awsapigatewayv2.NewHttpApi(construct, jsii.String("ExampleSamAPI"), &awsapigatewayv2.HttpApiProps{})

	// POST
	apiIntegration := awsapigatewayv2integrations.NewLambdaProxyIntegration(&awsapigatewayv2integrations.LambdaProxyIntegrationProps{
		Handler:              apiLambda,
		PayloadFormatVersion: awsapigatewayv2.PayloadFormatVersion_VERSION_1_0(),
	})

	httpapi.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Integration: apiIntegration,
		Path:        jsii.String("/version"),
		Methods: &[]awsapigatewayv2.HttpMethod{
			awsapigatewayv2.HttpMethod_GET,
		},
	})

	httpapi.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Integration: apiIntegration,
		Path:        jsii.String("/"),
		Methods: &[]awsapigatewayv2.HttpMethod{
			awsapigatewayv2.HttpMethod_GET,
		},
	})

	awscdk.NewCfnOutput(construct, jsii.String("APIEndpoint"), &awscdk.CfnOutputProps{
		Value: awscdk.Fn_Sub(jsii.String("https://${ExampleSamAPI}.execute-api.${AWS::Region}.amazonaws.com"), &map[string]*string{
			"ExampleSamAPI": httpapi.HttpApiId(),
		}),
		Description: jsii.String("API Gateway Endpoint URL"),
	})

	awscdk.NewCfnOutput(construct, jsii.String("FunctionARN"), &awscdk.CfnOutputProps{
		Value:       apiLambda.FunctionArn(),
		Description: jsii.String("Lambda Function ARN"),
	})

	return construct
}
