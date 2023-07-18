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

	"github.com/v6d-io/v6d/k8s/apis/k8s/v1alpha1"

	"github.com/v6d-io/v6d/k8s/cmd/commands/flags"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestBuildVineyard(t *testing.T) {
	opts := &flags.VineyarddOpts

	tests := []struct {
		name    string
		want    *v1alpha1.Vineyardd
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Test Case 1",
			want: &v1alpha1.Vineyardd{
				ObjectMeta: metav1.ObjectMeta{
					Name:      flags.VineyarddName,
					Namespace: flags.GetDefaultVineyardNamespace(),
				},
				Spec: *opts,
			}, // 指定预期的 *cobra.Command 值
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildVineyard()
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildVineyard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildVineyard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildVineyardManifestFromInput(t *testing.T) {
	opts := &flags.VineyarddOpts
	tests := []struct {
		name    string
		want    *v1alpha1.Vineyardd
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Test Case 1",
			want: &v1alpha1.Vineyardd{
				ObjectMeta: metav1.ObjectMeta{
					Name:      flags.VineyarddName,
					Namespace: flags.GetDefaultVineyardNamespace(),
				},
				Spec: *opts,
			}, // 指定预期的 *cobra.Command 值
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildVineyardManifestFromInput()
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildVineyardManifestFromInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildVineyardManifestFromInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildVineyardManifestFromFile(t *testing.T) {
	tests := []struct {
		name    string
		want    *v1alpha1.Vineyardd
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Test Case 1",
			want: &v1alpha1.Vineyardd{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "vineyard-system",
				},
				Spec: v1alpha1.VineyarddSpec{
					Replicas:     0,
					EtcdReplicas: 0,
					Service:      v1alpha1.ServiceConfig{},
					Vineyard:     v1alpha1.VineyardConfig{},
					PluginImage:  v1alpha1.PluginImageConfig{},
					Metric:       v1alpha1.MetricConfig{},
					Volume:       v1alpha1.VolumeConfig{PvcName: "", MountPath: ""},
				},
				Status: v1alpha1.VineyarddStatus{
					ReadyReplicas: 0,
					Conditions:    []appsv1.DeploymentCondition{},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags.Namespace = "vineyard-system"
			flags.VineyarddFile = "/home/zhuyi/v6d/k8s/config/crd/bases/k8s.v6d.io_vineyardds.yaml"
			got, err := BuildVineyardManifestFromFile()
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildVineyardManifestFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			a, _ := got.CreationTimestamp.Marshal()
			b, _ := tt.want.CreationTimestamp.Marshal()
			if !reflect.DeepEqual(a, b) {
				t.Errorf("BuildVineyardManifestFromFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

/*func TestNewDeployVineyarddCmd(t *testing.T) {
	tests := []struct {
		name string
		want *cobra.Command
	}{
		// TODO: Add test cases.
		{
			name: "Test Case 1",
			want: deployVineyarddCmd, // 指定预期的 *cobra.Command 值
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDeployVineyarddCmd(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDeployVineyarddCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}*/
