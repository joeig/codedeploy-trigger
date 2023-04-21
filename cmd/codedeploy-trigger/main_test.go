package main

import (
	"testing"
	"time"
)

func Test_failIfWrongTarget(t *testing.T) {
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
			if err := failIfWrongTarget(tt.args.flagName, tt.args.flagValue); (err != nil) != tt.wantErr {
				t.Errorf("failIfWrongTarget() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_failIfEmptyFlag(t *testing.T) {
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
			if err := failIfEmptyFlag(tt.args.flagName, tt.args.flagValue); (err != nil) != tt.wantErr {
				t.Errorf("failIfEmptyFlag() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_failIfOutOfPortRangeFlag(t *testing.T) {
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
			if err := failIfOutOfPortRangeFlag(tt.args.flagName, tt.args.flagValue); (err != nil) != tt.wantErr {
				t.Errorf("failIfOutOfPortRangeFlag() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_failIfInsufficientDuration(t *testing.T) {
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
			if err := failIfInsufficientDuration(tt.args.flagName, tt.args.flagValue); (err != nil) != tt.wantErr {
				t.Errorf("failIfInsufficientDuration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
