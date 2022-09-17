package deploy

import (
	"testing"
)

func TestAppSpec_MarshalJSON(t *testing.T) {
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
