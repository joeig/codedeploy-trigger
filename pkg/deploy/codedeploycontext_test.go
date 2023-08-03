package deploy

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codedeploy"
	"github.com/aws/aws-sdk-go-v2/service/codedeploy/types"
	"testing"
	"time"
)

func Test_assembleCreateDeploymentInput(t *testing.T) {
	input := assembleCreateDeploymentInput("app", "group", []byte("{}"))
	want := codedeploy.CreateDeploymentInput{
		ApplicationName:     aws.String("app"),
		DeploymentGroupName: aws.String("group"),
		Revision: &types.RevisionLocation{
			AppSpecContent: &types.AppSpecContent{
				Content: aws.String("{}"),
				Sha256:  aws.String("44136fa355b3678a1146ad16f7e8649e94fb4fc21fe77e8310c060f61caaff8a"),
			},
			RevisionType: types.RevisionLocationTypeAppSpecContent,
		},
	}

	inputJSON, _ := json.Marshal(input)
	wantJSON, _ := json.Marshal(want)

	if !bytes.Equal(inputJSON, wantJSON) {
		t.Error("wrong result")
	}
}

type mockCodeDeployClient struct {
	CreateDeploymentOutput *codedeploy.CreateDeploymentOutput
	CreateDeploymentErr    error
	GetDeploymentOutput    *codedeploy.GetDeploymentOutput
	GetDeploymentErr       error
}

func (m *mockCodeDeployClient) CreateDeployment(_ context.Context, _ *codedeploy.CreateDeploymentInput, _ ...func(*codedeploy.Options)) (*codedeploy.CreateDeploymentOutput, error) {
	return m.CreateDeploymentOutput, m.CreateDeploymentErr
}

func (m *mockCodeDeployClient) GetDeployment(_ context.Context, _ *codedeploy.GetDeploymentInput, _ ...func(*codedeploy.Options)) (*codedeploy.GetDeploymentOutput, error) {
	return m.GetDeploymentOutput, m.GetDeploymentErr
}

func TestCodeDeployContext_CreateDeployment(t *testing.T) {
	deploymentID := "mock"
	codeDeployContext, _ := NewCodeDeployContext(&mockCodeDeployClient{CreateDeploymentOutput: &codedeploy.CreateDeploymentOutput{DeploymentId: aws.String(deploymentID)}}, nil, nil).WithAppSpec(&AppSpec{})
	output, err := codeDeployContext.CreateDeployment(context.Background(), "a", "d")

	if output != deploymentID {
		t.Error("wrong deployment ID")
	}

	if err != nil {
		t.Error("unexpected error")
	}
}

func TestCodeDeployContext_CreateDeployment_error(t *testing.T) {
	codeDeployContext, _ := NewCodeDeployContext(&mockCodeDeployClient{CreateDeploymentErr: errors.New("mock")}, nil, nil).WithAppSpec(&AppSpec{})
	output, err := codeDeployContext.CreateDeployment(context.Background(), "a", "d")

	if output != "" {
		t.Error("unexpected deployment ID")
	}

	if err == nil {
		t.Error("no error")
	}
}

func NewMockDeploymentSuccessfulWaiter(waitErr error) DeploymentSuccessfulWaiter {
	return func(ctx context.Context, params *codedeploy.GetDeploymentInput, maxWaitDur time.Duration, optFns ...func(*codedeploy.DeploymentSuccessfulWaiterOptions)) error {
		return waitErr
	}
}

func TestCodeDeployContext_WaitForSuccessfulDeployment(t *testing.T) {
	codeDeployContext := CodeDeployContext{DeploymentSuccessfulWaiter: NewMockDeploymentSuccessfulWaiter(nil)}

	if err := codeDeployContext.WaitForSuccessfulDeployment(context.Background(), "mock", 1); err != nil {
		t.Error("unexpected err")
	}
}

func TestCodeDeployContext_WaitForSuccessfulDeployment_error(t *testing.T) {
	codeDeployContext := CodeDeployContext{Client: &mockCodeDeployClient{}, DeploymentSuccessfulWaiter: NewMockDeploymentSuccessfulWaiter(errors.New("mock"))}

	if err := codeDeployContext.WaitForSuccessfulDeployment(context.Background(), "mock", 1); err == nil {
		t.Error("no err")
	}
}

func TestCodeDeployContext_WaitForSuccessfulDeployment_error_details(t *testing.T) {
	codeDeployContext := CodeDeployContext{
		Client: &mockCodeDeployClient{
			GetDeploymentOutput: &codedeploy.GetDeploymentOutput{
				DeploymentInfo: &types.DeploymentInfo{
					ErrorInformation: &types.ErrorInformation{
						Message: aws.String("info"),
					},
				},
			},
		},
		DeploymentSuccessfulWaiter: NewMockDeploymentSuccessfulWaiter(errors.New("mock")),
	}

	err := codeDeployContext.WaitForSuccessfulDeployment(context.Background(), "mock", 1)
	if err.Error() != "mock" {
		t.Error("unexpected err")
	}
}

func TestCodeDeployContext_WithAppSpec(t *testing.T) {
	codeDeployContext := CodeDeployContext{}
	appSpec := &AppSpec{}

	if _, err := codeDeployContext.WithAppSpec(appSpec); err != nil {
		t.Error("unexpected err")
	}

	if !bytes.Equal(codeDeployContext.appSpecJson, []byte("{\"version\":\"\",\"Resources\":null}")) {
		t.Error("invalid app spec JSON")
	}
}

func NewMockFileReader(want []byte, wantErr error) FileReader {
	return func(fileName string) ([]byte, error) {
		return want, wantErr
	}
}

func TestCodeDeployContext_WithAppSpecFile(t *testing.T) {
	want := []byte("foo")
	codeDeployContext := CodeDeployContext{FileReader: NewMockFileReader(want, nil)}

	if _, err := codeDeployContext.WithAppSpecFile("foo"); err != nil {
		t.Error("unexpected err")
	}

	if !bytes.Equal(codeDeployContext.appSpecJson, []byte("foo")) {
		t.Error("invalid app spec JSON")
	}
}

func TestCodeDeployContext_WithAppSpecFile_Error(t *testing.T) {
	codeDeployContext := CodeDeployContext{FileReader: NewMockFileReader([]byte(""), errors.New("mock"))}

	if _, err := codeDeployContext.WithAppSpecFile("foo"); err == nil {
		t.Error("no err")
	}
}
