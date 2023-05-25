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

