package stacks

import (
	"api-sam-example-cdk/pkg/hosting"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/constructs-go/constructs/v10"
)

type ApplicationProps struct {
	Tenant      string            `envconfig:"TENANT" default:"openenterprise"`
	Environment string            `envconfig:"ENVIRONMENT" default:"staging"`
	StackProps  awscdk.StackProps ``
}

func ApplicationStack(scope constructs.Construct, id string, props *ApplicationProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}

	stack := awscdk.NewStack(scope, &id, &sprops)

	hosting.HostingStack(stack, "Hosting", &hosting.HostingProps{
		Tenant:      props.Tenant,
		Environment: props.Environment,
	})

	return stack
}
