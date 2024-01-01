# CodeDeploy Trigger

[![Tests](https://github.com/joeig/codedeploy-trigger/actions/workflows/tests.yml/badge.svg)](https://github.com/joeig/codedeploy-trigger/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/joeig/codedeploy-trigger)](https://goreportcard.com/report/github.com/joeig/codedeploy-trigger)

This handy tool creates a new CodeDeploy deployment for ECS services and Lambda functions and waits until it completes.
In case the deployment fails, it tells you why.

It's usually used in conjunction with CI/CD or as a Terraform provisioner.

## Usage

```shell
$ codedeploy-trigger -help
Usage of codedeploy-trigger:
  -appSpecFileName string
        Custom AppSpec file name
  -applicationName string
        CodeDeploy application name
  -containerName string
        ECS container name (if appSpecFileName is unset)
  -containerPort int
        ECS container port (if appSpecFileName is unset)
  -currentVersion string
        Current Lambda function version (if appSpecFileName is unset)
  -deploymentGroupName string
        CodeDeploy deployment group name
  -functionAlias string
        Lambda function alias (if appSpecFileName is unset)
  -functionName string
        Lambda function name (if appSpecFileName is unset)
  -maxWaitDuration duration
        Max wait duration for a deployment to finish (default 30m0s)
  -target string
        Deployment target ("ECS" or "Lambda"; if appSpecFileName is unset)
  -targetVersion string
        Target Lambda function version (if appSpecFileName is unset)
  -taskDefinitionARN string
        ECS task definition ARN (if appSpecFileName is unset)
```

## Install from source

The following command builds and installs `codedeploy-trigger` into your `GOBIN` directory (usually `~/go/bin`):

```shell
go install github.com/joeig/codedeploy-trigger/cmd/codedeploy-trigger@latest
```

## AWS client configuration

`codedeploy-trigger` assumes the [environment](https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/#environment-variables) provides a working AWS client configuration.

```shell
export AWS_REGION="eu-central-1"

export AWS_ACCESS_KEY_ID="access key ID"
export AWS_SECRET_ACCESS_KEY="secret access key"
export AWS_SESSION_TOKEN="token"  # Optional

# Alternatively:
export AWS_PROFILE="your profile"
```

## Terraform snippets

### ECS service

You can easily integrate it in your Terraform project using a [`null_resource`](https://registry.terraform.io/providers/hashicorp/null/latest/docs/resources/resource) which is triggered by a task definition change:

```terraform
resource "null_resource" "deploy" {
  triggers = {
    task_definition_arn = aws_ecs_task_definition.example.arn
  }

  provisioner "local-exec" {
    command = <<EOT
      codedeploy-trigger \
        -target "ECS" \
        -applicationName "codedeploy-application-name" \
        -deploymentGroupName "codedeploy-deployment-group-name" \
        -taskDefinitionARN "${aws_ecs_task_definition.api.arn}" \
        -containerName "container-name" \
        -containerPort "1337"
    EOT
  }
}
```

### Lambda function

You can easily integrate it in your Terraform project using a [`null_resource`](https://registry.terraform.io/providers/hashicorp/null/latest/docs/resources/resource) which is triggered by a function version change:

```terraform
resource "null_resource" "deploy" {
  triggers = {
    target_version = aws_lambda_function.example.version
  }

  provisioner "local-exec" {
    command = <<EOT
      codedeploy-trigger \
        -target "Lambda" \
        -applicationName "codedeploy-application-name" \
        -deploymentGroupName "codedeploy-deployment-group-name" \
        -functionName "function-name" \
        -functionAlias "function-alias" \
        -currentVersion "42" \
        -targetVersion "43"
    EOT
  }
}
```
