# CodeDeploy Trigger

This handy tool creates a new CodeDeploy deployment for an ECS task definition.

## Usage

```shell
$ codedeploy-trigger -help
Usage of codedeploy-trigger:
  -applicationName string
        CodeDeploy application name
  -containerName string
        ECS container name
  -containerPort int
        ECS container port
  -deploymentGroupName string
        CodeDeploy deployment group name
  -maxWaitDuration duration
        Max wait duration for a deployment to finish (default 30m0s)
  -taskDefinitionARN string
        ECS task definition ARN
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

## Terraform snippet

You can easily integrate it in your Terraform project using a [`null_resource`](https://registry.terraform.io/providers/hashicorp/null/latest/docs/resources/resource) which is triggered by a task definition change:

```terraform
resource "null_resource" "deploy" {
  triggers = {
    task_definition_arn = aws_ecs_task_definition.api.arn
  }

  provisioner "local-exec" {
    command = <<EOT
      codedeploy-trigger \
        -applicationName "codedeploy-application-name" \
        -deploymentGroupName "codedeploy-deployment-group-name" \
        -containerName "container-name" \
        -containerPort "1337" \
        -taskDefinitionARN "${aws_ecs_task_definition.api.arn}"
    EOT
  }
}
```
