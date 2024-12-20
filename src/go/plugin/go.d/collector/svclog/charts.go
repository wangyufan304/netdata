package svclog

import (
	"github.com/netdata/netdata/go/plugins/plugin/go.d/agent/module"
)

const (
	prioErrorLogs = module.Priority + iota
	prioWarningLogs
	prioInfoLogs
	prioTotalLogs
)

var logChartsTmpl = module.Charts{
	logChartTmple.Copy(),
}

// 假设 module.Chart 结构体已经定义好，这里直接初始化一个实例
var logChartTmple = module.Chart{
	ID:       "%s_log_id", // ID 字段
	Title:    "%s Log Data Chart",
	Units:    "counts",
	Fam:      "log.%s",
	Ctx:      "svclog.%s",
	Priority: 0,
	Type:     module.Area,
	Dims:     module.Dims{},
}

// func newLogStatsCharts(logStats map[string]int64) *module.Charts {
// 	charts := logStatsChartsTmpl.Copy()

// 	for _, chart := range *charts {
// 		for i, dim := range chart.Dims {
// 			switch dim.ID {
// 			case "error_logs":
// 				chart.Dims[i].Mul = int(logStats["error"]) // 使用 Mul 来存储统计值
// 			case "warning_logs":
// 				chart.Dims[i].Mul = int(logStats["warning"])
// 			case "info_logs":
// 				chart.Dims[i].Mul = int(logStats["info"])
// 			case "total_logs":
// 				chart.Dims[i].Mul = int(logStats["total"])
// 			}
// 		}
// 	}

// 	return charts
// }

// // 统计日志中错误、警告、信息的频率，并填充到图表
// func (c *Collector) logStatsToCharts(logStats map[string]int64) {
// 	// 创建新的图表
// 	charts := newLogStatsCharts(logStats)
// 	// 添加图表到 Collector
// 	if err := c.Charts().Add(*charts...); err != nil {
// 		c.Warning(err)
// 	}
// }
