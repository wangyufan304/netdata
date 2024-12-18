package svclog

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

func (c *Collector) collect() (map[string]int64, error) {

	mx := make(map[string]int64)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() { defer wg.Done(); c.logTest(mx) }()

	wg.Wait()

	return mx, nil
}
func (c *Collector) logTest(mx map[string]int64) {
	// 打开日志文件
	file, err := os.Open(c.Config.LogPath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	// 创建一个扫描器来逐行读取文件
	scanner := bufio.NewScanner(file)

	// 定义要统计的关键字
	keywords := []string{"error", "warning", "critical", "info"}

	// 逐行扫描文件
	for scanner.Scan() {
		line := scanner.Text() // 获取当前行的文本

		// 遍历每个关键字并检查是否出现在当前行中
		for _, keyword := range keywords {
			if strings.Contains(strings.ToLower(line), keyword) {
				// 如果该行包含关键字，更新对应的频率
				mx[keyword]++
			}
		}
	}

	// 检查扫描过程中是否有错误
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}
	c.logStatsToCharts(mx)
}
