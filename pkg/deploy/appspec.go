package deploy

import "encoding/json"

type Version string

const DefaultVersion = "0.0"

type TargetServiceType string

const ECSTargetServiceType TargetServiceType = "AWS::ECS::Service"

type LoadBalancerInfo struct {
	ContainerName string `json:"ContainerName"`
	ContainerPort int    `json:"ContainerPort"`
}

type Properties struct {
	TaskDefinition   string           `json:"TaskDefinition"`
	LoadBalancerInfo LoadBalancerInfo `json:"LoadBalancerInfo"`
}

type TargetService struct {
	Type       TargetServiceType `json:"Type"`
	Properties Properties        `json:"Properties"`
}

type Resource struct {
	TargetService TargetService `json:"TargetService"`
}

// AppSpec provides an application specification for CodeDeploy.
// Reference: https://docs.aws.amazon.com/codedeploy/latest/userguide/reference-appspec-file.html
type AppSpec struct {
	Version   Version    `json:"version"`
	Resources []Resource `json:"Resources"`
}

// NewECS creates a new instance of AppSpec incorporating commonly used default values for ECS deployments.
func NewECS(taskDefinitionARN, containerName string, containerPort int) *AppSpec {
	return &AppSpec{
		Version: DefaultVersion,
		Resources: []Resource{
			{
				TargetService: TargetService{
					Type: ECSTargetServiceType,
					Properties: Properties{
						TaskDefinition: taskDefinitionARN,
						LoadBalancerInfo: LoadBalancerInfo{
							ContainerName: containerName,
							ContainerPort: containerPort,
						},
					},
				},
			},
		},
	}
}

func (a *AppSpec) MarshalJSON() ([]byte, error) {
	return json.Marshal(*a)
}
