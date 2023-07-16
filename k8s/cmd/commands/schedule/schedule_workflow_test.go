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

package schedule

import (
	"context"
	"reflect"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/v6d-io/v6d/k8s/cmd/commands/flags"
	"github.com/v6d-io/v6d/k8s/cmd/commands/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestNewScheduleWorkflowCmd(t *testing.T) {
	tests := []struct {
		name string
		want *cobra.Command
	}{
		// TODO: Add test cases.
		{
			name: "Test Case 1",
			want: scheduleWorkflowCmd, // 指定预期的 *cobra.Command 值
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewScheduleWorkflowCmd(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewScheduleWorkflowCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

type mockClient struct {
	client.Client
	mock.Mock
}

// Create mock implementation
func (m *mockClient) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
	args := m.Called(ctx, obj)
	return args.Error(0)
}

// Get mock implementation
func (m *mockClient) Get(ctx context.Context, key client.ObjectKey, obj client.Object) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func TestSchedulingWorkflow(t *testing.T) {
	flags.KubeConfig = "/home/zhuyi/.kube/config"
	flags.WorkflowFile = "/home/zhuyi/v6d/k8s/test/e2e/scheduling-outside-cluster-demo/test.yaml"
	client := util.KubernetesClient()
	manifests, err := util.ReadFromFile(flags.WorkflowFile)
	objs, err := util.ParseManifestsToObjects([]byte(manifests))
	//fmt.Println(objs[0])
	err = SchedulingWorkflow(client, objs[2])

	// The function should return the error
	assert.NoError(t, err)

}
