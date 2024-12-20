package svclog

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func ReadOffset(serviceName string) (int64, error) {

	if offset, exists := sv_off[serviceName]; exists {
		return offset, nil
	}

	// 打开文件
	file, err := os.Open(offset_filenam)
	if err != nil {
		return 0, fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	// 使用 bufio.Scanner 逐行读取文件
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// 忽略空行
		if len(line) == 0 {
			continue
		}

		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		var offset int64
		_, err := fmt.Sscanf(parts[1], "%d", &offset)
		if err != nil {
			log.Printf("读取偏移量失败: %v", err)
			continue
		}
		sv_off[key] = offset
	}

	// 检查是否发生扫描错误
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("读取文件时发生错误: %v", err)
	}

	// 返回对应服务名称的偏移量
	if offset, exists := sv_off[serviceName]; exists {
		return offset, nil
	}

	// 如果未找到对应服务的偏移量，返回错误
	return 0, fmt.Errorf("未找到服务 %s 的偏移量", serviceName)
}

func UpdateOffset(serviceName string, newOffset int64) error {
	// 更新 map 中的服务偏移量
	sv_off[serviceName] = newOffset

	// 打开文件进行写操作
	file, err := os.OpenFile(offset_filenam, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	// 读取文件内容并更新
	var lines []string
	scanner := bufio.NewScanner(file)
	updated := false
	for scanner.Scan() {
		line := scanner.Text()
		// 如果当前行是要更新的服务名称，则更新偏移量
		if strings.HasPrefix(line, serviceName+" ") {
			lines = append(lines, fmt.Sprintf("%s %d", serviceName, newOffset))
			updated = true
		} else {
			lines = append(lines, line)
		}
	}

	// 如果文件中没有该服务的记录，则添加一行
	if !updated {
		lines = append(lines, fmt.Sprintf("%s %d", serviceName, newOffset))
	}

	// 重写文件
	file.Seek(0, 0)  // 重置文件指针到文件开始位置
	file.Truncate(0) // 清空文件内容

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("写入文件失败: %v", err)
		}
	}
	writer.Flush()

	return nil
}
