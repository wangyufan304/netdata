package svclog

import (
	"context"
	_ "embed"

	"github.com/netdata/netdata/go/plugins/plugin/go.d/agent/module"
	_ "github.com/netdata/netdata/go/plugins/plugin/go.d/pkg/logs"
)

//go:embed "config_schema.json"
var configSchema string

// 定义服务文件便宜地址
var offset_filenam string = "/var/log/service.offset"

var sv_off = make(map[string]int64)

func init() {
	module.Register("svclog", module.Creator{
		JobConfigSchema: configSchema,
		Defaults: module.Defaults{
			UpdateEvery: 5,
		},
		Create: func() module.Module {
			return New()
		},
		Config: func() any {
			return &Config{}
		},
	})
}

func New() *Collector {
	return &Collector{
		charts: &module.Charts{},
		sm:     make(map[string]bool),
		Config: Config{
			Services: []Service{
				{
					ServiceName: "nginx",
					LogAddress:  "/var/log/nginx/access.log",
					Rules: map[string]string{
						"rule1": "pattern1",
						"rule2": "pattern2",
					},
				},
				{
					ServiceName: "dpvs",
					LogAddress:  "/var/log/nginx/access.log",
					Rules: map[string]string{
						"das": "pattern1",
						"453": "pattern2",
					},
				},
			},
		},
	}
}

type (
	Config struct {
		Services []Service `yaml:"services" json:"services"` // 服务配置的数组
	}

	// Service 表示每个服务的配置
	Service struct {
		ServiceName string            `yaml:"service_name" json:"service_name"` // 服务名称
		LogAddress  string            `yaml:"log_address" json:"log_address"`   // 日志文件地址
		Rules       map[string]string `yaml:"rules" json:"rules"`               // 服务日志解析规则
	}
)

type Collector struct {
	module.Base
	Config `yaml:",inline" json:""`
	charts *module.Charts
	sm     map[string]bool
}

func (c *Collector) Configuration() any {
	return c.Config
}

func (c *Collector) Init(ctx context.Context) error {
	// 读取配置文件
	for _, v := range c.Services {
		ReadOffset(v.ServiceName)
	}
	return nil
}

func (c *Collector) Check(ctx context.Context) error {
	return nil
}

func (c *Collector) Charts() *module.Charts {
	return c.charts
}

// func (c *Collector) Charts() *module.Charts {
// 	return &module.Charts{
// 		{
// 			ID:    "random",
//
// 			Dims: module.Dims{
// 				{ID: "random0", Name: "random 0"},
// 				{ID: "random1", Name: "random 1"},
// 			},
// 		},
// 	}
// }

// func (c *Collector) Collect(context.Context) map[string]int64 {
// 	return map[string]int64{
// 		"error": rand.Int63n(100),
// 		"info":  rand.Int63n(100),
// 		"warn":  rand.Int63n(100),
// 		"total": rand.Int63n(100),
// 	}
// }

func (c *Collector) Collect(context.Context) map[string]int64 {
	mx, err := c.collect()
	if err != nil {
		c.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx

}

func (c *Collector) Cleanup(context.Context) {}
