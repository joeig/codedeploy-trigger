package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/codedeploy"
	"github.com/joeig/codedeploy-trigger/pkg/deploy"
	"log"
	"os"
	"time"
)

const (
	ECSTarget    string = "ECS"
	LambdaTarget string = "Lambda"
)

func failIfWrongTarget(flagName, flagValue string) error {
	if flagValue != ECSTarget && flagValue != LambdaTarget {
		return fmt.Errorf("attribute %q must be either %q or %q", flagName, ECSTarget, LambdaTarget)
	}
	return nil
}

func failIfEmptyFlag(flagName string, flagValue string) error {
	if len(flagValue) == 0 {
		return fmt.Errorf("attribute %q must not be empty", flagName)
	}
	return nil
}

func failIfOutOfPortRangeFlag(flagName string, flagValue int) error {
	if flagValue < 0 || flagValue > 65535 {
		return fmt.Errorf("attribute %q contains an invalid port number", flagName)
	}
	return nil
}

func failIfInsufficientDuration(flagName string, flagValue time.Duration) error {
	if flagValue <= 0 {
		return fmt.Errorf("attribute %q must be greater than zero", flagName)
	}
	return nil
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

	if err := failIfWrongTarget("target", *target); err != nil {
		log.Fatalf(err.Error())
	}
	if err := failIfInsufficientDuration("maxWaitDuration", *maxWaitDuration); err != nil {
		log.Fatalf(err.Error())
	}
	if err := failIfEmptyFlag("applicationName", *applicationName); err != nil {
		log.Fatalf(err.Error())
	}
	if err := failIfEmptyFlag("deploymentGroupName", *deploymentGroupName); err != nil {
		log.Fatalf(err.Error())
	}

	if *target == ECSTarget {
		if err := failIfEmptyFlag("taskDefinitionARN", *taskDefinitionARN); err != nil {
			log.Fatalf(err.Error())
		}
		if err := failIfEmptyFlag("containerName", *containerName); err != nil {
			log.Fatalf(err.Error())
		}
		if err := failIfOutOfPortRangeFlag("containerPort", *containerPort); err != nil {
			log.Fatalf(err.Error())
		}
	}

	if *target == LambdaTarget {
		if err := failIfEmptyFlag("functionName", *functionName); err != nil {
			log.Fatalf(err.Error())
		}
		if err := failIfEmptyFlag("functionAlias", *functionAlias); err != nil {
			log.Fatalf(err.Error())
		}
		if err := failIfEmptyFlag("currentVersion", *currentVersion); err != nil {
			log.Fatalf(err.Error())
		}
		if err := failIfEmptyFlag("targetVersion", *targetVersion); err != nil {
			log.Fatalf(err.Error())
		}
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

	if *target == ECSTarget {
		appSpec = deploy.NewECS(*taskDefinitionARN, *containerName, *containerPort)
	}

	if *target == LambdaTarget {
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
