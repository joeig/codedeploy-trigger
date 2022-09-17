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
	applicationName := flag.String("applicationName", "", "Client application name")
	deploymentGroupName := flag.String("deploymentGroupName", "", "Client deployment group name")
	containerName := flag.String("containerName", "", "ECS container name")
	containerPort := flag.Int("containerPort", 0, "ECS container port")
	taskDefinitionARN := flag.String("taskDefinitionARN", "", "ECS task definition ARN")
	maxWaitDuration := flag.Duration("maxWaitDuration", 30*time.Minute, "Max wait duration for a deployment to finish")

	flag.Parse()

	failIfEmptyFlag("applicationName", *applicationName)
	failIfEmptyFlag("deploymentGroupName", *deploymentGroupName)
	failIfEmptyFlag("containerName", *containerName)
	failIfOutOfPortRangeFlag("containerPort", *containerPort)
	failIfEmptyFlag("taskDefinitionARN", *taskDefinitionARN)
	failIfInsufficientDuration("maxWaitDuration", *maxWaitDuration)

	awsConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("cannot load AWS configuration: %s", err)
	}

	codeDeployContext := &deploy.CodeDeployContext{
		Client: codedeploy.NewFromConfig(awsConfig),
	}

	log.Printf("creating deployment for application %q (group %q)", *applicationName, *deploymentGroupName)

	appSpec := deploy.NewECS(*taskDefinitionARN, *containerName, *containerPort)

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
