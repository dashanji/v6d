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
package deploy

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestNewDeployRecoverJobCmd(t *testing.T) {
	tests := []struct {
		name string
		want *cobra.Command
	}{
		// TODO: Add test cases.
		{
			name: "Test Case 1",
			want: deployRecoverJobCmd, // 指定预期的 *cobra.Command 值
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDeployRecoverJobCmd(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDeployRecoverJobCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getRecoverObjectsFromTemplate(t *testing.T) {
	type args struct {
		c client.Client
	}
	tests := []struct {
		name    string
		args    args
		want    []*unstructured.Unstructured
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getRecoverObjectsFromTemplate(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("getRecoverObjectsFromTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getRecoverObjectsFromTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_waitRecoverJobReady(t *testing.T) {
	type args struct {
		c client.Client
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
			if err := waitRecoverJobReady(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("waitRecoverJobReady() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_createMappingTableConfigmap(t *testing.T) {
	type args struct {
		c  client.Client
		cs kubernetes.Clientset
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
			if err := createMappingTableConfigmap(tt.args.c, tt.args.cs); (err != nil) != tt.wantErr {
				t.Errorf("createMappingTableConfigmap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
