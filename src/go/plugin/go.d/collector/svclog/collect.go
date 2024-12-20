package svclog

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sync"

	"github.com/netdata/netdata/go/plugins/plugin/go.d/agent/module"
)

var gp int = 0

func (c *Collector) collect() (map[string]int64, error) {
	mu := &sync.Mutex{}
	mx := make(map[string]int64)
	var wg sync.WaitGroup

	wg.Add(len(c.Services))
	for _, service := range c.Services {
		go func(sv Service) { defer wg.Done(); c.logTest(sv, mx, mu) }(service)
	}

	wg.Wait()

	return mx, nil
}

func (c *Collector) logTest(sv Service, mx map[string]int64, mu *sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()
	if !c.sm[sv.ServiceName] {
		c.sm[sv.ServiceName] = true
		c.addLogCharts(sv)
	}
	newoffset, err := ReadAndMatchFromOffset(sv.LogAddress, sv_off[sv.ServiceName], sv.Rules, mx)
	if err == nil {
		UpdateOffset(sv.ServiceName, newoffset)
	}
}

func newLogCharts(sv Service) *module.Charts {
	charts := logChartsTmpl.Copy()
	var dims []*module.Dim
	for key := range sv.Rules {
		dims = append(dims, &module.Dim{
			ID:   key,
			Name: key + "_count",
		})
	}
	for _, chart := range *charts {
		chart.ID = fmt.Sprintf(chart.ID, sv.ServiceName)
		chart.Title = fmt.Sprintf(chart.Title, sv.ServiceName)
		chart.Fam = fmt.Sprintf(chart.Fam, sv.ServiceName)
		chart.Ctx = fmt.Sprintf(chart.Ctx, sv.ServiceName)
		chart.Priority = gp
		gp++
		chart.Labels = []module.Label{
			{Key: "service_name", Value: sv.ServiceName},
		}
		chart.Dims = dims
	}
	return charts
}

func (c *Collector) addLogCharts(sv Service) {
	charts := newLogCharts(sv)
	if err := c.Charts().Add(*charts...); err != nil {
		c.Warning(err)
	}
}

// 定义函数来从指定文件的偏移量读取数据，并根据规则匹配更新结果
func ReadAndMatchFromOffset(filePath string, offset int64, rules map[string]string, resultMap map[string]int64) (int64, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return offset, fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	// 将文件指针移动到指定偏移量位置
	_, err = file.Seek(offset, 0)
	if err != nil {
		return offset, fmt.Errorf("移动文件指针到偏移量 %d 时出错: %v", offset, err)
	}

	// 使用 bufio.Scanner 逐行读取文件
	scanner := bufio.NewScanner(file)
	var newOffset int64 = offset

	// 遍历文件中的每一行
	for scanner.Scan() {
		line := scanner.Text()
		// 遍历规则并尝试匹配每行数据
		for key, regexPattern := range rules {
			// 编译正则表达式
			re, err := regexp.Compile(regexPattern)
			if err != nil {
				log.Printf("正则表达式编译失败: %v", err)
				continue
			}

			// 如果匹配成功，则更新 map
			if re.MatchString(line) {
				// 匹配成功，更新匹配次数
				resultMap[key]++
				log.Printf("匹配到: %s, 匹配次数: %d", key, resultMap[key])
			}
		}

		// 更新偏移量：扫描的下一行的开始位置
		newOffset, err = file.Seek(0, 1)
		if err != nil {
			return newOffset, fmt.Errorf("获取新的偏移量时出错: %v", err)
		}
	}

	// 检查文件扫描过程中是否发生了错误
	if err := scanner.Err(); err != nil {
		return newOffset, fmt.Errorf("扫描文件时发生错误: %v", err)
	}

	// 返回新的偏移量
	return newOffset, nil
}
