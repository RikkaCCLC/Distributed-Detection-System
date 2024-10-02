package models

import (
	"regexp"
)

// 用户配置的日志策略，可以是agent 本地的yaml，也可以是通过接口过来的
type LogStrategy struct {
	Id         int64  `json:"id" yaml:"-"`
	MetricName string `json:"metric_name" yaml:"metric_name"`
	MetricHelp string `json:"metric_help" yaml:"metric_help"`
	FilePath   string `json:"file_path" yaml:"file_path"`
	Pattern    string `json:"pattern" yaml:"pattern"`
	Func       string `json:"func" yaml:"func"`

	Creator string `json:"creator"`

	// 上面是yaml或者前端的配置
	TagJson    string                    `json:"tag_json" yaml:"-"`      // 这是给db用的
	Tags       map[string]string         `json:"-" yaml:"tags" xorm:"-"` // 这是yaml用的
	PatternReg *regexp.Regexp            `json:"-" yaml:"-" xorm:"-"`    // 主正则
	TagRegs    map[string]*regexp.Regexp `json:"-" yaml:"-" xorm:"-"`    // 标签的正则 ，这是最终的

}

func (lm *LogStrategy) TableName() string {
	return "log_job"
}

// 带参数查询一条记录函数 level=3 and path=/0
func LogJobGets(where string, args ...interface{}) ([]*LogStrategy, error) {
	var obj []*LogStrategy
	err := DB["stree"].Table("log_job").Where(where, args...).Find(&obj)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 插入一条记录
func (lm *LogStrategy) AddOne() (int64, error) {
	_, err := DB["stree"].InsertOne(lm)
	return lm.Id, err
}
