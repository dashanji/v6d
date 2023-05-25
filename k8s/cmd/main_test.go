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

package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMainExecute(t *testing.T) {
	// 设置测试输入
	os.Args = []string{"vineyardctl", "create"}
	os.Args = []string{"vineyardctl", "delete"}
	os.Args = []string{"vineyardctl", "deploy"}
	os.Args = []string{"vineyardctl", "manager"}
	os.Args = []string{"vineyardctl", "schedule"}
	//os.Args = []string{"vineyardctl", "sidecar"}



	// 执行被测试的函数
	err := cmd.Execute()

	// 验证测试结果是否符合预期
	assert.NoError(t, err, "Expected no error")
}

