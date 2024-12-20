package svclog

import (
	"fmt"
	"sync"
	"testing"
)

// TestLogTest 测试 logTest 方法
// func TestLogTest(t *testing.T) {
// 	co := New()
// 	mx, _ := co.collect()
// 	fmt.Println(mx)
// 	for _, chart := range *co.charts {
// 		fmt.Println(chart.Fam)
// 		fmt.Println(chart.Ctx)
// 	}
// 	co.Init(nil)
// 	fmt.Println(sv_off)

// }

func TestLogTest(t *testing.T) {
	// 示例文件路径和规则文件路径
	filePath := "/tmp/testlog.txt"
	offset := int64(0) // 从文件的开头开始读取
	rules := make(map[string]string)
	rules["service1"] = "^service1\\s+\\d+$"
	rules["service2"] = "^service2\\s+\\d+$"
	rules["service3"] = "^service3\\s+\\d+$"

	// 用于存储匹配结果的 map，存储匹配次数
	resultMap := make(map[string]int64)

	// 调用 ReadAndMatchFromOffset 函数进行匹配
	newOffset, err := ReadAndMatchFromOffset(filePath, offset, rules, resultMap)
	if err != nil {
		fmt.Println("读取和匹配时出错:", err)
		return
	}

	// 输出匹配的结果
	fmt.Println("匹配结果:")
	for key, count := range resultMap {
		fmt.Printf("Key: %s, 匹配次数: %d\n", key, count)
	}

	// 输出新的偏移量
	fmt.Printf("新的偏移量: %d\n", newOffset)
}

// Test function for logTest
func TestLogTes1t(t *testing.T) {
	// Prepare test data
	sv := Service{
		ServiceName: "nginxtest",
		LogAddress:  "/var/log/nginx/test.log",
		Rules:       map[string]string{"rule1": "GET.*", "rule2": "ERROR.*"},
	}

	// Initialize map and mutex
	mx := make(map[string]int64)
	mu := &sync.Mutex{}

	// Initialize Collector
	collector := New()
	collector.Init(nil)
	// Call logTest
	collector.logTest(sv, mx, mu)
	fmt.Println(mx)

}
