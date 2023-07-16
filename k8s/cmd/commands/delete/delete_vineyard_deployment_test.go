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
package delete

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"github.com/v6d-io/v6d/k8s/cmd/commands/deploy"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestNewDeleteVineyardDeploymentCmd(t *testing.T) {
	tests := []struct {
		name string
		want *cobra.Command
	}{
		// TODO: Add test cases.
		{
			name: "Test Case 1",
			want: deleteVineyardDeploymentCmd, // 指定预期的 *cobra.Command 值
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDeleteVineyardDeploymentCmd(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDeleteVineyardDeploymentCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}


func Test_deleteVineyarddFromTemplate(t *testing.T) {
	type args struct {
		c client.Client
	}
	
	objects, err := deploy.GetVineyardDeploymentObjectsFromTemplate()
	if err != nil {
		t.Fatalf("Failed to get vineyardd resources from template: %v", err)
	}

	// 将 []*unstructured.Unstructured 转换为 []runtime.Object
	runtimeObjects := make([]runtime.Object, len(objects))
	for i, obj := range objects {
		runtimeObjects[i] = obj
	}

	client := fake.NewFakeClient(runtimeObjects...) // 将对象添加到假的 Kubernetes 集群中
	
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Test Case 1",
			args: args{
				// 提供测试所需的参数值
				c:    client,
			},
			wantErr: false, // 设置预期的错误结果
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := deleteVineyarddFromTemplate(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("deleteVineyarddFromTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

