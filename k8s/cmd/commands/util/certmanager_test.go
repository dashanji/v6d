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
package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCertManagerManifestsFromLocal(t *testing.T) {
	// 执行 getCertManagerManifestsFromLocal 函数
	manifests, err := getCertManagerManifestsFromLocal()

	// 验证错误是否为 nil
	assert.Nil(t, err, "Unexpected error")

	// 验证返回的 manifests 是否符合预期
	assert.NotNil(t, manifests, "Manifests should not be nil")
	assert.NotEmpty(t, manifests, "Manifests should not be empty")
}

func TestGetCertManager(t *testing.T) {
	// 执行 GetCertManager 函数
	manifests, err := GetCertManager()

	// 验证错误是否为 nil
	assert.Nil(t, err, "Unexpected error")

	// 验证返回的 manifests 是否符合预期
	assert.NotNil(t, manifests, "Manifests should not be nil")
	assert.NotEmpty(t, manifests, "Manifests should not be empty")
}
