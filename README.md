# CDK SAM Demo

Demo using sam to drive a CDK serverless api

# Stack Setup

```bash
go build -v ./cmd/application
🛠️  cmd/application done
✓  stacks/build done
✓  build done
cdk deploy --app ./application
2021/11/03 21:56:33 Starting Application Build
Bundling asset OpenenterpriseProductionSamAppStack/Hosting/Lambda/Code/Stage...
OpenenterpriseProductionSamAppStack: deploying...
[0%] start: Publishing e9deedc87e534faa84cd401ae14a2e47c1e235ff4cbab5bbadd3d810bf38c6ec:current
[50%] success: Published e9deedc87e534faa84cd401ae14a2e47c1e235ff4cbab5bbadd3d810bf38c6ec:current
[50%] start: Publishing 54aa9368cb0af77728b139df544966eb0ee9754a76a567e67d1db78a861cdb72:current
[100%] success: Published 54aa9368cb0af77728b139df544966eb0ee9754a76a567e67d1db78a861cdb72:current
OpenenterpriseProductionSamAppStack: creating CloudFormation changeset...

 ✅  OpenenterpriseProductionSamAppStack

Stack ARN:
arn:aws:cloudformation:eu-west-1:074705540277:stack/OpenenterpriseProductionSamAppStack/18c8c170-3ce5-11ec-8c2b-06e6d7c1014f
🛠️  deploy/application done
✓  deploy done
```

# License

MIT

## Useful commands

 * `make deploy`            deploy this stack to your default AWS account/region
 * `make sam/test/api`      start the api
 * `make sam/test/hello`    call the hello endpoint
 * `make sam/test/version`  call the version endpoint

## Notes

Had to set the following due to https://github.com/aws/aws-sam-cli/issues/2849

```bash
"@aws-cdk/core:newStyleStackSynthesis": false,
```


