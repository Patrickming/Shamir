package test

import (
	"reflect"
	"shamir/sss" // 导入sss包
	"testing"
)

// TestToByteArray 测试ToByteArray函数的正确性
func TestToByteArray(t *testing.T) {
	// 定义一组测试用例
	tests := []struct {
		name     string      // 测试用例的名字，用于标识测试
		input    interface{} // 测试输入，使用interface{}允许输入任何类型
		expected []byte      // 期望的输出，字节数组形式
		wantErr  bool        // 是否期望此测试产生错误
	}{
		// 定义具体的测试实例
		{"Int", 123, []byte("123"), false},                          // 测试整数输入
		{"Float", 98.76, []byte("98.76"), false},                    // 测试浮点数输入
		{"String", "Hello, world!", []byte("Hello, world!"), false}, // 测试字符串输入
		{"HexString", "0x49cb693f7b4f21fde2a2b24c31af1440f861069f", []byte("0x49cb693f7b4f21fde2a2b24c31af1440f861069f"), false}, // 测试十六进制字符串输入
		{"Unsupported", []int{1, 2, 3}, nil, true}, // 测试不支持的类型输入，期望产生错误
	}

	// 遍历所有测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 执行ToByteArray函数
			got, err := sss.ToByteArray(tt.input)
			// 检查错误是否如预期发生
			if (err != nil) != tt.wantErr {
				t.Errorf("ToByteArray() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// 比较实际输出和期望输出是否一致
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ToByteArray() got = %v, want %v", got, tt.expected)
			}
		})
	}
}
