package main

import (
	"flag"
	"testing"
	"time"
)

func Test_checkTarget(t *testing.T) {
	type args struct {
		flagName  string
		flagValue string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ECS target",
			args: args{
				flagName:  "test",
				flagValue: ECSTarget,
			},
			wantErr: false,
		},
		{
			name: "Lambda target",
			args: args{
				flagName:  "test",
				flagValue: LambdaTarget,
			},
			wantErr: false,
		},
		{
			name: "Unknown target",
			args: args{
				flagName:  "test",
				flagValue: "unknown",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkTarget(tt.args.flagName, tt.args.flagValue); (err != nil) != tt.wantErr {
				t.Errorf("checkTarget() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checkNotEmpty(t *testing.T) {
	type args struct {
		flagName  string
		flagValue string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "available",
			args: args{
				flagName:  "test",
				flagValue: "foo",
			},
			wantErr: false,
		},
		{
			name: "empty",
			args: args{
				flagName:  "test",
				flagValue: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkNotEmpty(tt.args.flagName, tt.args.flagValue); (err != nil) != tt.wantErr {
				t.Errorf("checkNotEmpty() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checkPortRange(t *testing.T) {
	type args struct {
		flagName  string
		flagValue int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "lower boundary",
			args: args{
				flagName:  "test",
				flagValue: 0,
			},
			wantErr: false,
		},
		{
			name: "upper boundary",
			args: args{
				flagName:  "test",
				flagValue: 65535,
			},
			wantErr: false,
		},
		{
			name: "below longer boundary",
			args: args{
				flagName:  "test",
				flagValue: -1,
			},
			wantErr: true,
		},
		{
			name: "above upper boundary",
			args: args{
				flagName:  "test",
				flagValue: 65536,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkPortRange(tt.args.flagName, tt.args.flagValue); (err != nil) != tt.wantErr {
				t.Errorf("checkPortRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checkDuration(t *testing.T) {
	type args struct {
		flagName  string
		flagValue time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "boundary",
			args: args{
				flagName:  "test",
				flagValue: 1,
			},
			wantErr: false,
		},
		{
			name: "out of boundary",
			args: args{
				flagName:  "test",
				flagValue: 0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkDuration(tt.args.flagName, tt.args.flagValue); (err != nil) != tt.wantErr {
				t.Errorf("checkDuration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFlagContext_Parse(t *testing.T) {
	flagContext := &FlagContext{FlagSet: flag.CommandLine}
	arguments := []string{
		"-maxWaitDuration",
		"15m",
		"-applicationName",
		"my-app",
		"-deploymentGroupName",
		"my-group",
		"-target",
		"Lambda",
		"-appSpecFileName",
		"my-file",
		"-taskDefinitionARN",
		"my-task-def",
		"-containerName",
		"my-container",
		"-containerPort",
		"1337",
		"-functionName",
		"my-function",
		"-functionAlias",
		"my-alias",
		"-currentVersion",
		"my-current-version",
		"-targetVersion",
		"my-target-version",
	}

	if err := flagContext.Parse(arguments); err != nil {
		t.Error("unexpected error")
	}

	if flagContext.maxWaitDuration.String() != "15m0s" {
		t.Error("unexpected max wait duration")
	}

	if *flagContext.applicationName != "my-app" {
		t.Error("unexpected application name")
	}

	if *flagContext.deploymentGroupName != "my-group" {
		t.Error("unexpected deployment group name")
	}

	if *flagContext.target != "Lambda" {
		t.Error("unexpected target")
	}

	if *flagContext.appSpecFileName != "my-file" {
		t.Error("unexpected app spec file name")
	}

	if *flagContext.taskDefinitionARN != "my-task-def" {
		t.Error("unexpected task definition ARN")
	}

	if *flagContext.containerName != "my-container" {
		t.Error("unexpected container name")
	}

	if *flagContext.containerPort != 1337 {
		t.Error("unexpected container Port")
	}

	if *flagContext.functionName != "my-function" {
		t.Error("unexpected function name")
	}

	if *flagContext.functionAlias != "my-alias" {
		t.Error("unexpected function alias")
	}

	if *flagContext.currentVersion != "my-current-version" {
		t.Error("unexpected current version")
	}
	if *flagContext.targetVersion != "my-target-version" {
		t.Error("unexpected target version")
	}
}
