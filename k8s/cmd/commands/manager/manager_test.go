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

// Package start contains the start command of vineyard operator
package manager

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func TestNewManagerCmd(t *testing.T) {
	tests := []struct {
		name string
		want *cobra.Command
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewManagerCmd(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewManagerCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_startManager(t *testing.T) {
	type args struct {
		mgr                  manager.Manager
		metricsAddr          string
		probeAddr            string
		enableLeaderElection bool
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startManager(tt.args.mgr, tt.args.metricsAddr, tt.args.probeAddr, tt.args.enableLeaderElection)
		})
	}
}

func Test_startScheduler(t *testing.T) {
	type args struct {
		mgr                 manager.Manager
		schedulerConfigFile string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startScheduler(tt.args.mgr, tt.args.schedulerConfigFile)
		})
	}
}
