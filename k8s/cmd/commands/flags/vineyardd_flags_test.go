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
package flags

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/v6d-io/v6d/k8s/apis/k8s/v1alpha1"
)

func TestApplyVineyardContainerOpts(t *testing.T) {
	type args struct {
		c      *v1alpha1.VineyardConfig
		prefix string
		cmd    *cobra.Command
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ApplyVineyardContainerOpts(tt.args.c, tt.args.prefix, tt.args.cmd)
		})
	}
}

func TestApplyServiceOpts(t *testing.T) {
	type args struct {
		s      *v1alpha1.ServiceConfig
		prefix string
		cmd    *cobra.Command
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ApplyServiceOpts(tt.args.s, tt.args.prefix, tt.args.cmd)
		})
	}
}

func TestApplyVolumeOpts(t *testing.T) {
	type args struct {
		v      *v1alpha1.VolumeConfig
		prefix string
		cmd    *cobra.Command
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ApplyVolumeOpts(tt.args.v, tt.args.prefix, tt.args.cmd)
		})
	}
}

func TestApplyMetricContainerOpts(t *testing.T) {
	type args struct {
		m      *v1alpha1.MetricConfig
		prefix string
		cmd    *cobra.Command
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ApplyMetricContainerOpts(tt.args.m, tt.args.prefix, tt.args.cmd)
		})
	}
}

func TestApplyPluginImageOpts(t *testing.T) {
	type args struct {
		cmd *cobra.Command
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ApplyPluginImageOpts(tt.args.cmd)
		})
	}
}

func TestApplyVineyarddNameOpts(t *testing.T) {
	type args struct {
		cmd *cobra.Command
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ApplyVineyarddNameOpts(tt.args.cmd)
		})
	}
}

func TestApplyVineyarddOpts(t *testing.T) {
	type args struct {
		cmd *cobra.Command
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ApplyVineyarddOpts(tt.args.cmd)
		})
	}
}
