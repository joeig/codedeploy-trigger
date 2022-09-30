package deploy

import "encoding/json"

type Version string

const DefaultVersion = "0.0"

type TargetServiceType string

const (
	ECSTargetServiceType    TargetServiceType = "AWS::ECS::Service"
	LambdaTargetServiceType TargetServiceType = "AWS::Lambda::Function"
)

type LoadBalancerInfo struct {
	ContainerName string `json:"ContainerName"`
	ContainerPort int    `json:"ContainerPort"`
}

type ECSProperties struct {
	TaskDefinition   string           `json:"TaskDefinition"`
	LoadBalancerInfo LoadBalancerInfo `json:"LoadBalancerInfo"`
}

type LambdaProperties struct {
	Name           string `json:"Name"`
	Alias          string `json:"Alias"`
	CurrentVersion string `json:"CurrentVersion"`
	TargetVersion  string `json:"TargetVersion"`
}

type TargetService struct {
	Type       TargetServiceType `json:"Type"`
	Properties any               `json:"Properties"`
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

// NewECS creates a new instance of AppSpec incorporating commonly used default values for ECS service deployments.
func NewECS(taskDefinitionARN, containerName string, containerPort int) *AppSpec {
	return &AppSpec{
		Version: DefaultVersion,
		Resources: []Resource{
			{
				TargetService: TargetService{
					Type: ECSTargetServiceType,
					Properties: ECSProperties{
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

// NewLambda creates a new instance of AppSpec incorporating commonly used default values for Lambda function deployments.
func NewLambda(functionName, functionAlias, currentVersion, targetVersion string) *AppSpec {
	return &AppSpec{
		Version: DefaultVersion,
		Resources: []Resource{
			{
				TargetService: TargetService{
					Type: LambdaTargetServiceType,
					Properties: LambdaProperties{
						Name:           functionName,
						Alias:          functionAlias,
						CurrentVersion: currentVersion,
						TargetVersion:  targetVersion,
					},
				},
			},
		},
	}
}

func (a *AppSpec) MarshalJSON() ([]byte, error) {
	return json.Marshal(*a)
}
