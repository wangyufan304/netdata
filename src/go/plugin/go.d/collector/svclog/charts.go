package svclog

import (
	"fmt"
	"github.com/netdata/netdata/go/plugins/plugin/go.d/agent/module"
)

const (
	prioErrorLogs = module.Priority + iota
	prioWarningLogs
	prioInfoLogs
	prioTotalLogs
)

var logStatsChartsTmpl = module.Charts{
	errorLogsChartTmpl.Copy(),
	warningLogsChartTmpl.Copy(),
	infoLogsChartTmpl.Copy(),
	totalLogsChartTmpl.Copy(),
}

var (
	errorLogsChartTmpl = module.Chart{
		ID:       "error_logs",
		Title:    "Error Logs Count",
		Units:    "count",
		Fam:      "logs",
		Ctx:      "logs.error_logs",
		Priority: prioErrorLogs,
		Type:     module.Area,
		Dims: module.Dims{
			{ID: "error_count", Name: "error_count"},
		},
	}
	warningLogsChartTmpl = module.Chart{
		ID:       "warning_logs",
		Title:    "Warning Logs Count",
		Units:    "count",
		Fam:      "logs",
		Ctx:      "logs.warning_logs",
		Priority: prioWarningLogs,
		Type:     module.Area,
		Dims: module.Dims{
			{ID: "warning_count", Name: "warning_count"},
		},
	}
	infoLogsChartTmpl = module.Chart{
		ID:       "info_logs",
		Title:    "Info Logs Count",
		Units:    "count",
		Fam:      "logs",
		Ctx:      "logs.info_logs",
		Priority: prioInfoLogs,
		Type:     module.Area,
		Dims: module.Dims{
			{ID: "info_count", Name: "info_count"},
		},
	}
	totalLogsChartTmpl = module.Chart{
		ID:       "total_logs",
		Title:    "Total Logs Count",
		Units:    "count",
		Fam:      "logs",
		Ctx:      "logs.total_logs",
		Priority: prioTotalLogs,
		Type:     module.Area,
		Dims: module.Dims{
			{ID: "total_count", Name: "total_count"},
		},
	}
)

func newLogStatsCharts() *module.Charts {
	charts := logStatsChartsTmpl.Copy()

	for _, chart := range *charts {
		chart.Labels = []module.Label{
			{Key: "log_type", Value: "general"},
		}
	}

	return charts
}

// 统计日志中错误、警告、信息的频率，并填充到图表
func (c *Collector) logStatsToCharts(logStats map[string]int64) {
	// 创建新的图表
	charts := newLogStatsCharts()

	// 根据日志统计数据填充图表
	for _, chart := range *charts {
		for i, dim := range chart.Dims {
			// 根据维度ID来为维度赋值
			switch dim.ID {
			case "error_logs":
				chart.Dims[i].Mul = int(logStats["error"]) // 使用 Mul 来存储统计值
			case "warning_logs":
				chart.Dims[i].Mul = int(logStats["warning"])
			case "info_logs":
				chart.Dims[i].Mul = int(logStats["info"])
			case "total_logs":
				chart.Dims[i].Mul = int(logStats["total"])
			}
		}
	}

	// 添加图表到 Collector
	if err := c.Charts().Add(*charts...); err != nil {
		c.Warning(err)
	}
}

// 模拟日志采集并统计错误、警告、信息的频次
func (c *Collector) collectLogStats(logFilePath string) {
	// 假设我们从日志文件中读取内容并计算每种日志级别的频率
	// 以下是一个模拟的统计
	logStats := map[string]int64{
		"error":   10, // 模拟的错误日志计数
		"warning": 5,  // 模拟的警告日志计数
		"info":    20, // 模拟的信息日志计数
		"total":   35, // 总日志计数
	}

	c.logStatsToCharts(logStats)
}

func (c *Collector) Warning(err error) {
	// 打印警告信息
	fmt.Println("Warning:", err)
}
