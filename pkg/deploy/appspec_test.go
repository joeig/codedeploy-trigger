package deploy

import (
	"testing"
)

func TestAppSpec_NewECS(t *testing.T) {
	appSpec := NewECS("this:is:the:arn", "containerName", 1337)

	result, err := appSpec.MarshalJSON()
	if err != nil {
		t.Error(err)
	}

	want := `{"version":"0.0","Resources":[{"TargetService":{"Type":"AWS::ECS::Service","Properties":{"TaskDefinition":"this:is:the:arn","LoadBalancerInfo":{"ContainerName":"containerName","ContainerPort":1337}}}}]}`
	if string(result) != want {
		t.Error("resulting JSON is wrong")
	}
}

func TestAppSpec_NewLambda(t *testing.T) {
	appSpec := NewLambda("function-name", "function-alias", "42", "43")

	result, err := appSpec.MarshalJSON()
	if err != nil {
		t.Error(err)
	}

	want := `{"version":"0.0","Resources":[{"TargetService":{"Type":"AWS::Lambda::Function","Properties":{"Name":"function-name","Alias":"function-alias","CurrentVersion":"42","TargetVersion":"43"}}}]}`
	if string(result) != want {
		t.Error("resulting JSON is wrong")
	}
}
