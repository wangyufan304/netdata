package svclog

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

// 创建临时文件并写入测试数据
func createTempFileWithData(data string) (*os.File, error) {
	tmpFile, err := ioutil.TempFile("", "service.offset")
	if err != nil {
		return nil, err
	}

	// 写入数据
	_, err = tmpFile.WriteString(data)
	if err != nil {
		return nil, err
	}

	// 重置文件指针到开头
	tmpFile.Seek(0, 0)
	return tmpFile, nil
}

// 测试 ReadOffset 和 UpdateOffset 函数
func TestReadAndUpdateOffset(t *testing.T) {
	// 测试数据：服务名和偏移量
	testData := "service1 1000\nservice2 2000\nservice3 3000\n"

	// 创建临时文件并写入测试数据
	tmpFile, err := createTempFileWithData(testData)
	if err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // 确保测试结束后删除临时文件

	// 设置全局变量 offset_filenam 为临时文件路径
	offset_filenam = tmpFile.Name()

	// 测试 ReadOffset 函数
	t.Run("TestReadOffset", func(t *testing.T) {
		// 读取 service1 的偏移量
		offset, err := ReadOffset("service1")
		if err != nil {
			t.Fatalf("读取服务 'service1' 偏移量时出错: %v", err)
		}
		if offset != 1000 {
			t.Errorf("期望偏移量为 1000，但获取到 %d", offset)
		}

		// 读取 service2 的偏移量
		offset, err = ReadOffset("service2")
		if err != nil {
			t.Fatalf("读取服务 'service2' 偏移量时出错: %v", err)
		}
		if offset != 2000 {
			t.Errorf("期望偏移量为 2000，但获取到 %d", offset)
		}
	})

	// 测试 UpdateOffset 函数
	t.Run("TestUpdateOffset", func(t *testing.T) {
		// 更新 service1 的偏移量为 4000
		err := UpdateOffset("service1", 4000)
		if err != nil {
			t.Fatalf("更新服务 'service1' 偏移量时出错: %v", err)
		}

		// 再次读取 service1 的偏移量，应该为 4000
		offset, err := ReadOffset("service1")
		if err != nil {
			t.Fatalf("读取服务 'service1' 偏移量时出错: %v", err)
		}
		if offset != 4000 {
			t.Errorf("期望偏移量为 4000，但获取到 %d", offset)
		}

		// 验证文件内容是否更新
		fileData, err := ioutil.ReadFile(tmpFile.Name())
		if err != nil {
			t.Fatalf("读取文件时出错: %v", err)
		}

		// 检查文件中的内容是否已更新
		fileContent := string(fileData)
		if !strings.Contains(fileContent, "service1 4000") {
			t.Errorf("文件中没有找到更新后的 'service1 4000', 文件内容: %s", fileContent)
		}
	})
}
