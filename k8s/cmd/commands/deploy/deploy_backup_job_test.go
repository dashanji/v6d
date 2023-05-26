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
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestNewDeployBackupJobCmd(t *testing.T) {
	tests := []struct {
		name string
		want *cobra.Command
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDeployBackupJobCmd(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDeployBackupJobCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getBackupObjectsFromTemplate(t *testing.T) {
	type args struct {
		c    client.Client
		args []string
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
			got, err := getBackupObjectsFromTemplate(tt.args.c, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("getBackupObjectsFromTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getBackupObjectsFromTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_waitBackupJobReady(t *testing.T) {
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
			if err := waitBackupJobReady(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("waitBackupJobReady() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
