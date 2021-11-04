# Meta tasks
# ----------

# Useful variables
export SAM_CLI_TELEMETRY ?= 0

# deployment environment
export ENVIRONMENT ?= production

# region
export AWS_REGION ?= eu-west-1

#
export AWS_ACCOUNT ?= 074705540277

export CODEBUILD_BUILD_NUMBER ?= 0
export CODEBUILD_RESOLVED_SOURCE_VERSION ?=$(shell git rev-list -1 HEAD --abbrev-commit)
export DATE=$(shell date -u '+%Y%m%d')

# Output helpers
# --------------

TASK_DONE = echo "✓  $@ done"
TASK_BUILD = echo "🛠️  $@ done"

# ----------------
STACKS = $(shell find ./cmd/ -mindepth 1 -maxdepth 1 -type d)

.DEFAULT_GOAL := build

clean:
	@rm -rf cdk.out .aws-sam application
	@$(TASK_DONE)

test:
	go test -v -p 1 ./...
	@$(TASK_BUILD)

bootstrap:
	CDK_NEW_BOOTSTRAP=1 cdk bootstrap aws://$(AWS_ACCOUNT)/$(AWS_REGION) --require-approval never --cloudformation-execution-policies=arn:aws:iam::aws:policy/AdministratorAccess --show-template
	@$(TASK_BUILD)

diff: diff/application
	@$(TASK_DONE)

synth: synth/application
	@$(TASK_DONE)

deploy: deploy/application
	@$(TASK_DONE)

synth/application: build
	cdk synth --app ./application
	@$(TASK_BUILD)

diff/application: build
	cdk diff --app ./application
	@$(TASK_BUILD)

deploy/application: build
	cdk deploy --app ./application
	@$(TASK_BUILD)

ci/deploy/application: build
	cdk deploy --app ./application --ci true --require-approval never 
	@$(TASK_BUILD)

build: stacks/build
	@$(TASK_DONE)

sam/build:
	@sam-beta-cdk build --debug
	@$(TASK_BUILD)

sam/test/version: sam/build
	@sam-beta-cdk local invoke "OpenenterpriseProductionSamAppStack/Lambda" --env-vars ./test/testenvironment.json --event ./test/events/version.json

sam/test/hello: sam/build
	@sam-beta-cdk local invoke "OpenenterpriseProductionSamAppStack/Lambda" --env-vars ./test/testenvironment.json --event ./test/events/hello.json

sam/test/api: sam/build
	@sam-beta-cdk local start-api --env-vars ./test/testenvironment.json

.PHONY: stacks/build $(STACKS)

stacks/build: $(STACKS)
	@$(TASK_DONE)

$(STACKS):
	go build -v ./$@
	@$(TASK_BUILD)    
	
init: 
	go mod download
	@$(TASK_BUILD)

