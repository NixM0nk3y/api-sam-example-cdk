package hosting

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/aws-cdk-go/awscdkapigatewayv2alpha/v2"
	"github.com/aws/aws-cdk-go/awscdkapigatewayv2integrationsalpha/v2"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type CommandHooks struct {
}

// needed to allow sam local testing
func (c CommandHooks) AfterBundling(inputDir *string, outputDir *string) *[]*string {
	return jsii.Strings(fmt.Sprintf("cp ../../test/sam.Makefile %s/Makefile", *outputDir))
}

func (c CommandHooks) BeforeBundling(inputDir *string, outputDir *string) *[]*string {
	return &[]*string{}
}

type HostingProps struct {
	Tenant           string                  ``
	Environment      string                  ``
	Appplication     string                  ``
	NestedStackProps awscdk.NestedStackProps ``
}

func HostingStack(scope constructs.Construct, id string, props *HostingProps) constructs.Construct {

	construct := constructs.NewConstruct(scope, &id)

	buildNumber, ok := os.LookupEnv("CODEBUILD_BUILD_NUMBER")
	if !ok {
		// default version
		buildNumber = "0"
	}

	sourceVersion, ok := os.LookupEnv("CODEBUILD_RESOLVED_SOURCE_VERSION")
	if !ok {
		sourceVersion = "unknown"
	}

	buildDate, ok := os.LookupEnv("BUILD_DATE")
	if !ok {
		t := time.Now()
		buildDate = t.Format("20060102")
	}

	// Go build options
	bundlingOptions := &awscdklambdagoalpha.BundlingOptions{
		GoBuildFlags: &[]*string{jsii.String(fmt.Sprintf(`-ldflags "-s -w
			-X api/pkg/version.Version=1.0.%s
			-X api/pkg/version.BuildHash=%s
			-X api/pkg/version.BuildDate=%s
			"`,
			buildNumber,
			sourceVersion,
			buildDate,
		)),
		},
		Environment: &map[string]*string{
			"GOARCH":      jsii.String("arm64"),
			"GO111MODULE": jsii.String("on"),
			"GOOS":        jsii.String("linux"),
		},
		CommandHooks: &CommandHooks{},
	}

	// webhook lambda
	apiLambda := awscdklambdagoalpha.NewGoFunction(construct, jsii.String("Lambda"), &awscdklambdagoalpha.GoFunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Entry:        jsii.String("resources/api/cmd/api"),
		Bundling:     bundlingOptions,
		Tracing:      awslambda.Tracing_ACTIVE,
		LogRetention: awslogs.RetentionDays_ONE_WEEK,
		Architecture: awslambda.Architecture_ARM_64(),
		Environment: &map[string]*string{
			"LOG_LEVEL": jsii.String("DEBUG"),
		},
		ModuleDir: jsii.String("resources/api/go.mod"),
	})

	//
	httpapi := awscdkapigatewayv2alpha.NewHttpApi(construct, jsii.String("ExampleSamAPI"), &awscdkapigatewayv2alpha.HttpApiProps{})

	// POST
	apiIntegration := awscdkapigatewayv2integrationsalpha.NewHttpLambdaIntegration(jsii.String("helloworld"), apiLambda, &awscdkapigatewayv2integrationsalpha.HttpLambdaIntegrationProps{
		PayloadFormatVersion: awscdkapigatewayv2alpha.PayloadFormatVersion_VERSION_1_0(),
	})

	httpapi.AddRoutes(&awscdkapigatewayv2alpha.AddRoutesOptions{
		Integration: apiIntegration,
		Path:        jsii.String("/version"),
		Methods: &[]awscdkapigatewayv2alpha.HttpMethod{
			awscdkapigatewayv2alpha.HttpMethod_GET,
		},
	})

	httpapi.AddRoutes(&awscdkapigatewayv2alpha.AddRoutesOptions{
		Integration: apiIntegration,
		Path:        jsii.String("/hello"),
		Methods: &[]awscdkapigatewayv2alpha.HttpMethod{
			awscdkapigatewayv2alpha.HttpMethod_GET,
		},
	})

	return construct
}
