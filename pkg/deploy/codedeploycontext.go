package deploy

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codedeploy"
	"github.com/aws/aws-sdk-go-v2/service/codedeploy/types"
	"time"
)

func assembleCreateDeploymentInput(applicationName, deploymentGroupName string, appSpecJson []byte) *codedeploy.CreateDeploymentInput {
	appSpecJsonHash := sha256.Sum256(appSpecJson)

	return &codedeploy.CreateDeploymentInput{
		ApplicationName:     aws.String(applicationName),
		DeploymentGroupName: aws.String(deploymentGroupName),
		Revision: &types.RevisionLocation{
			AppSpecContent: &types.AppSpecContent{
				Content: aws.String(string(appSpecJson)),
				Sha256:  aws.String(fmt.Sprintf("%x", appSpecJsonHash)),
			},
			RevisionType: types.RevisionLocationTypeAppSpecContent,
		},
	}
}

type CodeDeployClient interface {
	CreateDeployment(ctx context.Context, params *codedeploy.CreateDeploymentInput, optFns ...func(*codedeploy.Options)) (*codedeploy.CreateDeploymentOutput, error)
	GetDeployment(ctx context.Context, params *codedeploy.GetDeploymentInput, optFns ...func(*codedeploy.Options)) (*codedeploy.GetDeploymentOutput, error)
}

type DeploymentSuccessfulWaiter interface {
	Wait(ctx context.Context, params *codedeploy.GetDeploymentInput, maxWaitDur time.Duration, optFns ...func(*codedeploy.DeploymentSuccessfulWaiterOptions)) error
}

type CodeDeployContext struct {
	Client                     CodeDeployClient
	DeploymentSuccessfulWaiter DeploymentSuccessfulWaiter
}

func (c *CodeDeployContext) CreateDeployment(ctx context.Context, applicationName, deploymentGroupName string, appSpec json.Marshaler) (string, error) {
	appSpecJson, err := appSpec.MarshalJSON()
	if err != nil {
		return "", fmt.Errorf("cannot marshal JSON: %w", err)
	}

	deployment, err := c.Client.CreateDeployment(ctx, assembleCreateDeploymentInput(applicationName, deploymentGroupName, appSpecJson))
	if err != nil {
		return "", fmt.Errorf("cannot create deployment: %w", err)
	}

	return *deployment.DeploymentId, nil
}

func (c *CodeDeployContext) WaitForSuccessfulDeployment(ctx context.Context, deploymentID string, maxWaitDur time.Duration) error {
	deployment := &codedeploy.GetDeploymentInput{DeploymentId: aws.String(deploymentID)}

	if err := c.DeploymentSuccessfulWaiter.Wait(ctx, deployment, maxWaitDur); err != nil {
		if output, getDeploymentErr := c.Client.GetDeployment(ctx, deployment); output != nil && output.DeploymentInfo != nil && output.DeploymentInfo.ErrorInformation != nil && output.DeploymentInfo.ErrorInformation.Message != nil && getDeploymentErr != nil {
			return fmt.Errorf("%s (%w)", *output.DeploymentInfo.ErrorInformation.Message, err)
		}

		return err
	}

	return nil
}
