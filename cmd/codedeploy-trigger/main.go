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

func checkTarget(flagName, flagValue string) error {
	if flagValue != ECSTarget && flagValue != LambdaTarget {
		return fmt.Errorf("attribute %q must be either %q or %q", flagName, ECSTarget, LambdaTarget)
	}
	return nil
}

func checkNotEmpty(flagName string, flagValue string) error {
	if len(flagValue) == 0 {
		return fmt.Errorf("attribute %q must not be empty", flagName)
	}
	return nil
}

func checkPortRange(flagName string, flagValue int) error {
	if flagValue < 0 || flagValue > 65535 {
		return fmt.Errorf("attribute %q contains an invalid port number", flagName)
	}
	return nil
}

func checkDuration(flagName string, flagValue time.Duration) error {
	if flagValue <= 0 {
		return fmt.Errorf("attribute %q must be greater than zero", flagName)
	}
	return nil
}

type FlagContext struct {
	FlagSet *flag.FlagSet

	maxWaitDuration     *time.Duration
	applicationName     *string
	deploymentGroupName *string
	appSpecFileName     *string
	target              *string
	taskDefinitionARN   *string
	containerName       *string
	containerPort       *int
	functionName        *string
	functionAlias       *string
	currentVersion      *string
	targetVersion       *string
}

func (f *FlagContext) Parse(arguments []string) error {
	f.maxWaitDuration = f.FlagSet.Duration("maxWaitDuration", 30*time.Minute, "Max wait duration for a deployment to finish")
	f.applicationName = f.FlagSet.String("applicationName", "", "CodeDeploy application name")
	f.deploymentGroupName = f.FlagSet.String("deploymentGroupName", "", "CodeDeploy deployment group name")
	f.target = f.FlagSet.String("target", "", "Deployment target (\"ECS\" or \"Lambda\"; if appSpecFileName is unset)")
	f.appSpecFileName = f.FlagSet.String("appSpecFileName", "", "Custom AppSpec file name")
	f.taskDefinitionARN = f.FlagSet.String("taskDefinitionARN", "", "ECS task definition ARN (if appSpecFileName is unset)")
	f.containerName = f.FlagSet.String("containerName", "", "ECS container name (if appSpecFileName is unset)")
	f.containerPort = f.FlagSet.Int("containerPort", 0, "ECS container port (if appSpecFileName is unset)")
	f.functionName = f.FlagSet.String("functionName", "", "Lambda function name (if appSpecFileName is unset)")
	f.functionAlias = f.FlagSet.String("functionAlias", "", "Lambda function alias (if appSpecFileName is unset)")
	f.currentVersion = f.FlagSet.String("currentVersion", "", "Current Lambda function version (if appSpecFileName is unset)")
	f.targetVersion = f.FlagSet.String("targetVersion", "", "Target Lambda function version (if appSpecFileName is unset)")

	if err := f.FlagSet.Parse(arguments); err != nil {
		return err
	}

	return f.validate()
}

func (f *FlagContext) validate() error {
	if err := checkDuration("maxWaitDuration", *f.maxWaitDuration); err != nil {
		return err
	}
	if err := checkNotEmpty("applicationName", *f.applicationName); err != nil {
		return err
	}
	if err := checkNotEmpty("deploymentGroupName", *f.deploymentGroupName); err != nil {
		return err
	}

	if f.appSpecFileName != nil {
		if err := checkNotEmpty("appSpecFileName", *f.appSpecFileName); err != nil {
			return err
		}
	} else {
		if err := checkTarget("target", *f.target); err != nil {
			return err
		}

		if *f.target == ECSTarget {
			if err := checkNotEmpty("taskDefinitionARN", *f.taskDefinitionARN); err != nil {
				return err
			}
			if err := checkNotEmpty("containerName", *f.containerName); err != nil {
				return err
			}
			if err := checkPortRange("containerPort", *f.containerPort); err != nil {
				return err
			}
		}

		if *f.target == LambdaTarget {
			if err := checkNotEmpty("functionName", *f.functionName); err != nil {
				return err
			}
			if err := checkNotEmpty("functionAlias", *f.functionAlias); err != nil {
				return err
			}
			if err := checkNotEmpty("currentVersion", *f.currentVersion); err != nil {
				return err
			}
			if err := checkNotEmpty("targetVersion", *f.targetVersion); err != nil {
				return err
			}
		}
	}

	return nil
}

func main() {
	flagContext := &FlagContext{FlagSet: flag.CommandLine}
	if err := flagContext.Parse(os.Args[1:]); err != nil {
		log.Fatalln(err)
	}

	awsConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("cannot load AWS configuration: %s", err)
	}

	codeDeployClient := codedeploy.NewFromConfig(awsConfig)
	codeDeployContext := &deploy.CodeDeployContext{
		Client:                     codeDeployClient,
		DeploymentSuccessfulWaiter: codedeploy.NewDeploymentSuccessfulWaiter(codeDeployClient).Wait,
		FileReader:                 os.ReadFile,
	}

	log.Printf("creating deployment for application %q (group %q)", *flagContext.applicationName, *flagContext.deploymentGroupName)

	if flagContext.appSpecFileName != nil {
		var appSpec *deploy.AppSpec

		if *flagContext.target == ECSTarget {
			appSpec = deploy.NewECS(*flagContext.taskDefinitionARN, *flagContext.containerName, *flagContext.containerPort)
		}

		if *flagContext.target == LambdaTarget {
			appSpec = deploy.NewLambda(*flagContext.functionName, *flagContext.functionAlias, *flagContext.currentVersion, *flagContext.targetVersion)
		}

		if _, err := codeDeployContext.WithAppSpec(appSpec); err != nil {
			log.Fatalf(err.Error())
		}
	} else {
		if _, err := codeDeployContext.WithAppSpecFile(*flagContext.appSpecFileName); err != nil {
			log.Fatalf(err.Error())
		}
	}

	deploymentID, err := codeDeployContext.CreateDeployment(context.Background(), *flagContext.applicationName, *flagContext.deploymentGroupName)
	if err != nil {
		log.Printf("cannot create deployment: %s", err)
		os.Exit(1)
	}

	log.Printf("deployment ID %q created", deploymentID)
	log.Printf("waiting for deployment ID %q to finish", deploymentID)

	if err := codeDeployContext.WaitForSuccessfulDeployment(context.Background(), deploymentID, *flagContext.maxWaitDuration); err != nil {
		log.Printf("deployment failed: %s", err)
		os.Exit(1)
	}

	log.Print("deployment finished successfully")
}
