package svclog

import (
	"context"
	_ "embed"
	"github.com/netdata/netdata/go/plugins/plugin/go.d/agent/module"
	_ "github.com/netdata/netdata/go/plugins/plugin/go.d/pkg/logs"
)

//go:embed "config_schema.json"
var configSchema string

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
		Config: Config{
			UpdateEvery: 5,
			LogPath:     "/var/log/test.log",
		},
	}
}

type Config struct {
	UpdateEvery int    `yaml:"update_every,omitempty" json:"update_every"` // 数据收集时间间隔（秒）
	Vnode       string `yaml:"vnode,omitempty" json:"vnode"`               // 虚拟节点名称
	LogPath     string `yaml:"log_path,omitempty" json:"log_path"`
}

type Collector struct {
	module.Base
	Config `yaml:",inline" json:""`
	charts *module.Charts
}

func (c *Collector) Configuration() any {
	return c.Config
}

func (c *Collector) Init(ctx context.Context) error {
	return nil
}

func (c *Collector) Check(ctx context.Context) error {
	return nil
}

func (c *Collector) Charts() *module.Charts {
	return c.charts
}

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
