

# Output helpers
# --------------

TASK_DONE = echo "✓  $@ done"
TASK_BUILD = echo "🛠️  $@ done"

# Deno Helpers
# --------------
DENO_VERSION := 1.15.3
DENO_INSTALL := .
include deno.mk

# various app related vars
export BUILD_NUMBER ?= 0
export RESOLVED_SOURCE_VERSION ?=$(shell git rev-list -1 HEAD --abbrev-commit)
export RESOLVED_BRANCH ?=$(shell git rev-parse --abbrev-ref HEAD)
export BUILD_DATE=$(shell date -u '+%Y%m%d')
export ARTIFACTS_DIR ?= /asset-output

all: lambda/build

deps: denobin
	@DENO_DIR=.deno_dir ./bin/deno cache app.ts 
	@cp -R .deno_dir/gen/file/${PWD}/ .deno_dir/LAMBDA_TASK_ROOT
	@$(TASK_DONE)

denobin:
	@test -f ./bin/deno || make ./bin/deno
	@$(TASK_DONE)

version:
	@printf "{\n    \"gitHash\": \"${RESOLVED_SOURCE_VERSION}\",\n    \"buildBranch\": \"${RESOLVED_BRANCH}\",\n    \"buildNumber\": \"1.0.${BUILD_NUMBER}\",\n    \"buildDate\": \"%s\"\n}\n" ${BUILD_DATE} > version.json
	@$(TASK_DONE)

clean:
	@rm -rf .deno_dir lambda.zip version.json ./bin/deno
	@$(TASK_DONE)

lambda/build: deps version
	zip lambda.zip -x '.deno_dir/gen/file/*' -r .deno_dir bootstrap bin app.ts version.json
	@$(TASK_DONE)

lambda/package: lambda/build
	cp ./lambda.zip ${ARTIFACTS_DIR}/
	@$(TASK_BUILD)
