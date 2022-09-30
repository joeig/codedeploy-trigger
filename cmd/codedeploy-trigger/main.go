package main

import (
	"context"
	"flag"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/codedeploy"
	"github.com/joeig/codedeploy-trigger/pkg/deploy"
	"log"
	"os"
	"time"
)

const (
	ecsTarget    string = "ECS"
	lambdaTarget string = "Lambda"
)

func failIfWrongTarget(flagName, flagValue string) {
	if flagValue != ecsTarget && flagValue != lambdaTarget {
		log.Fatalf("attribute %q must be either %q or %q", flagName, ecsTarget, lambdaTarget)
	}
}

func failIfEmptyFlag(flagName string, flagValue string) {
	if len(flagValue) == 0 {
		log.Fatalf("attribute %q must not be empty", flagName)
	}
}

func failIfOutOfPortRangeFlag(flagName string, flagValue int) {
	if flagValue < 0 || flagValue > 65535 {
		log.Fatalf("attribute %q contains an invalid port number", flagName)
	}
}

func failIfInsufficientDuration(flagName string, flagValue time.Duration) {
	if flagValue <= 0 {
		log.Fatalf("attribute %q must be greater than zero", flagName)
	}
}

func main() {
	target := flag.String("target", "", "Deployment target (\"ECS\" or \"Lambda\")")
	maxWaitDuration := flag.Duration("maxWaitDuration", 30*time.Minute, "Max wait duration for a deployment to finish")
	applicationName := flag.String("applicationName", "", "CodeDeploy application name")
	deploymentGroupName := flag.String("deploymentGroupName", "", "CodeDeploy deployment group name")

	taskDefinitionARN := flag.String("taskDefinitionARN", "", "ECS task definition ARN")
	containerName := flag.String("containerName", "", "ECS container name")
	containerPort := flag.Int("containerPort", 0, "ECS container port")

	functionName := flag.String("functionName", "", "Lambda function name")
	functionAlias := flag.String("functionAlias", "", "Lambda function alias")
	currentVersion := flag.String("currentVersion", "", "Current Lambda function version")
	targetVersion := flag.String("targetVersion", "", "Target Lambda function version")

	flag.Parse()

	failIfWrongTarget("target", *target)
	failIfInsufficientDuration("maxWaitDuration", *maxWaitDuration)
	failIfEmptyFlag("applicationName", *applicationName)
	failIfEmptyFlag("deploymentGroupName", *deploymentGroupName)

	if *target == ecsTarget {
		failIfEmptyFlag("taskDefinitionARN", *taskDefinitionARN)
		failIfEmptyFlag("containerName", *containerName)
		failIfOutOfPortRangeFlag("containerPort", *containerPort)
	}

	if *target == lambdaTarget {
		failIfEmptyFlag("functionName", *functionName)
		failIfEmptyFlag("functionAlias", *functionAlias)
		failIfEmptyFlag("currentVersion", *currentVersion)
		failIfEmptyFlag("targetVersion", *targetVersion)
	}

	awsConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("cannot load AWS configuration: %s", err)
	}

	codeDeployContext := &deploy.CodeDeployContext{
		Client: codedeploy.NewFromConfig(awsConfig),
	}

	log.Printf("creating deployment for application %q (group %q)", *applicationName, *deploymentGroupName)

	var appSpec *deploy.AppSpec

	if *target == ecsTarget {
		appSpec = deploy.NewECS(*taskDefinitionARN, *containerName, *containerPort)
	}

	if *target == lambdaTarget {
		appSpec = deploy.NewLambda(*functionName, *functionAlias, *currentVersion, *targetVersion)
	}

	deploymentID, err := codeDeployContext.CreateDeployment(context.Background(), *applicationName, *deploymentGroupName, appSpec)
	if err != nil {
		log.Printf("cannot create deployment: %s", err)
		os.Exit(1)
	}

	log.Printf("deployment ID %q created", deploymentID)
	log.Printf("waiting for deployment ID %q to finish", deploymentID)

	if err := codeDeployContext.WaitForSuccessfulDeployment(context.Background(), deploymentID, *maxWaitDuration); err != nil {
		log.Printf("deployment failed: %s", err)
		os.Exit(1)
	}

	log.Print("deployment finished successfully")
}
