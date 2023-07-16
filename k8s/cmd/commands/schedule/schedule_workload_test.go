/** Copyright 2020-2023 Alibaba Group Holding Limited.

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
	"reflect"
	"testing"

	"github.com/spf13/cobra"
	"github.com/v6d-io/v6d/k8s/cmd/commands/flags"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/stretchr/testify/assert"

	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	core "k8s.io/client-go/testing"
	"k8s.io/client-go/tools/clientcmd"
)

func TestNewScheduleWorkloadCmd(t *testing.T) {
	tests := []struct {
		name string
		want *cobra.Command
	}{
		// TODO: Add test cases.
		{
			name: "Test Case 1",
			want: scheduleWorkloadCmd, // 指定预期的 *cobra.Command 值
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewScheduleWorkloadCmd(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewScheduleWorkloadCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

/*func TestValidateWorkloadKind(t *testing.T) {
	type args struct {
		kind string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateWorkloadKind(tt.args.kind); got != tt.want {
				t.Errorf("ValidateWorkloadKind() = %v, want %v", got, tt.want)
			}
		})
	}
}*/

func TestValidateWorkloadKind(t *testing.T) {
	tests := []struct {
		name     string
		kind     string
		expected bool
	}{
		{
			name:     "Valid Deployment kind",
			kind:     "Deployment",
			expected: true,
		},
		{
			name:     "Valid StatefulSet kind",
			kind:     "StatefulSet",
			expected: true,
		},
		{
			name:     "Valid ReplicaSet kind",
			kind:     "ReplicaSet",
			expected: true,
		},
		{
			name:     "Valid Job kind",
			kind:     "Job",
			expected: true,
		},
		{
			name:     "Valid CronJob kind",
			kind:     "CronJob",
			expected: true,
		},
		{
			name:     "Invalid kind",
			kind:     "InvalidKind",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateWorkloadKind(tt.kind)
			assert.Equal(t, tt.expected, result)
		})
	}
}

/*func Test_getWorkload(t *testing.T) {
	type args struct {
		workload string
	}
	tests := []struct {
		name    string
		args    args
		want    *unstructured.Unstructured
		want1   bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := getWorkload(tt.args.workload)
			if (err != nil) != tt.wantErr {
				t.Errorf("getWorkload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getWorkload() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getWorkload() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}*/

func TestGetWorkload(t *testing.T) {
	workload := `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
spec:
  replicas: 3
  template:
    metadata:
      labels:
        app: my-app
    spec:
      containers:
      - name: my-container
        image: nginx:latest
`

	obj, isWorkload, err := getWorkload(workload)
	assert.NoError(t, err)
	assert.NotNil(t, obj)
	assert.True(t, isWorkload)
	assert.Equal(t, "Deployment", obj.GetKind())
	assert.Equal(t, "my-deployment", obj.GetName())
}

/*func TestSchedulingWorkload(t *testing.T) {
	type args struct {
		c               client.Client
		unstructuredObj *unstructured.Unstructured
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SchedulingWorkload(tt.args.c, tt.args.unstructuredObj)
			if (err != nil) != tt.wantErr {
				t.Errorf("SchedulingWorkload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SchedulingWorkload() = %v, want %v", got, tt.want)
			}
		})
	}
}*/

func TestSchedulingWorkload(t *testing.T) {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      flags.VineyarddName,
			Namespace: flags.VineyarddNamespace,
		},
	}

	clientMock := &fake.Clientset{}
	clientMock.AddReactor("get", "deployments", func(action core.Action) (bool, runtime.Object, error) {
		return true, deployment, nil
	})

	unstructuredObj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "Deployment",
			"apiVersion": "apps/v1",
			"metadata": map[string]interface{}{
				"name":      "example-deployment",
				"namespace": "default",
			},
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"spec": map[string]interface{}{
						"affinity": map[string]interface{}{
							"podAffinity": map[string]interface{}{
								"requiredDuringSchedulingIgnoredDuringExecution": []interface{}{},
							},
						},
					},
				},
			},
		},
	}

	expectedResult := `{"apiVersion":"apps/v1","kind":"Deployment","metadata":{"name":"example-deployment","namespace":"default"},"spec":{"template":{"spec":{"affinity":{"podAffinity":{"requiredDuringSchedulingIgnoredDuringExecution":[{"labelSelector":{"matchExpressions":[{"key":"app.kubernetes.io/instance","operator":"In","values":["vineyard-system-vineyardd-sample"]}]},"topologyKey":"kubernetes.io/hostname"}]}}}}}}` + "\n"

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	kubeconfig := filepath.Join(homeDir, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)

	clientScheme := runtime.NewScheme()
	_ = scheme.AddToScheme(clientScheme)
	c, err := client.New(config, client.Options{Scheme: clientScheme})

	if err != nil {
		t.Fatalf("Cannot create client, error: %v", err)
	}
	result, err := SchedulingWorkload(c, unstructuredObj)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}
