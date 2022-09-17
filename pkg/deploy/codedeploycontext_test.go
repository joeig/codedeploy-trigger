package deploy

import (
	"bytes"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codedeploy"
	"github.com/aws/aws-sdk-go-v2/service/codedeploy/types"
	"testing"
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
