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
)

func TestNewDeployVineyardClusterCmd(t *testing.T) {
	tests := []struct {
		name string
		want *cobra.Command
	}{
		// TODO: Add test cases.
		{
			name: "Test Case 1",
			want: deployVineyardClusterCmd, // 指定预期的 *cobra.Command 值
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDeployVineyardClusterCmd(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDeployVineyardClusterCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}
