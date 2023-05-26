/*
* Copyright 2020-2023 Alibaba Group Holding Limited.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package sidecar

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
	"github.com/v6d-io/v6d/k8s/apis/k8s/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Test_validateFormat(t *testing.T) {
	type args struct {
		format string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateFormat(tt.args.format); (err != nil) != tt.wantErr {
				t.Errorf("validateFormat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getWorkloadResource(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getWorkloadResource()
			if (err != nil) != tt.wantErr {
				t.Errorf("getWorkloadResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getWorkloadResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetWorkloadObj(t *testing.T) {
	type args struct {
		workload string
	}
	tests := []struct {
		name    string
		args    args
		want    *unstructured.Unstructured
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetWorkloadObj(tt.args.workload)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetWorkloadObj() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetWorkloadObj() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetManifestFromTemplate(t *testing.T) {
	type args struct {
		workload string
	}
	tests := []struct {
		name    string
		args    args
		want    OutputManifests
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetManifestFromTemplate(tt.args.workload)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetManifestFromTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetManifestFromTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseManifestsAsYAML(t *testing.T) {
	type args struct {
		om OutputManifests
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseManifestsAsYAML(tt.args.om)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseManifestsAsYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseManifestsAsYAML() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_deployDuringInjection(t *testing.T) {
	type args struct {
		om *OutputManifests
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := deployDuringInjection(tt.args.om); (err != nil) != tt.wantErr {
				t.Errorf("deployDuringInjection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_outputInjectedResult(t *testing.T) {
	type args struct {
		om OutputManifests
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputInjectedResult(tt.args.om); (err != nil) != tt.wantErr {
				t.Errorf("outputInjectedResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_buildSidecar(t *testing.T) {
	type args struct {
		namespace string
	}
	tests := []struct {
		name    string
		args    args
		want    *v1alpha1.Sidecar
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildSidecar(tt.args.namespace)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildSidecar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildSidecar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInjectSidecarConfig(t *testing.T) {
	type args struct {
		sidecar     *v1alpha1.Sidecar
		workloadObj *unstructured.Unstructured
		sidecarObj  *unstructured.Unstructured
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InjectSidecarConfig(tt.args.sidecar, tt.args.workloadObj, tt.args.sidecarObj); (err != nil) != tt.wantErr {
				t.Errorf("InjectSidecarConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewInjectCmd(t *testing.T) {
	tests := []struct {
		name string
		want *cobra.Command
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewInjectCmd(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInjectCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}
