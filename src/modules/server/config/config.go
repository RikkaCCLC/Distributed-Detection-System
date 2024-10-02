package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	MysqlS        []*MySQLConf         `yaml:"mysql_s"`
	RpcAddr       string               `yaml:"rpc_addr"`
	HttpAddr      string               `yaml:"http_addr"`
	PCC           *PublicCloudSyncConf `yaml:"public_cloud_sync"`
	IndexModules  []*IndexModuleConf   `yaml:"index_modules"`
	ProberTargets []*Targets           `yaml:"prober_targets"`
}
type Targets struct {
	ProberType string   `yaml:"prober_type"`
	Region     string   `yaml:"region"`
	Target     []string `yaml:"target"`
}

type IndexModuleConf struct {
	Enable       bool   `yaml:"enable"`
	ResourceName string `yaml:"resource_name"`
	Modulus      int    `yaml:"modulus"`
	Num          int    `yaml:"num"`
}

type PublicCloudSyncConf struct {
	Enable bool `yaml:"enable"`
}
type MySQLConf struct {
	Name  string `yaml:"name"`
	Addr  string `yaml:"addr"`
	Max   int    `yaml:"max"`
	Idle  int    `yaml:"idle"`
	Debug bool   `yaml:"debug"`
}

// 根据io read读取配置文件后的字符串解析yaml
func Load(s []byte) (*Config, error) {
	cfg := &Config{}

	err := yaml.Unmarshal(s, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// 根据conf路径读取内容
func LoadFile(filename string) (*Config, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg, err := Load(content)
	if err != nil {
		fmt.Printf("[parsing Yaml file err...][err:%v]\n", err)
		return nil, err
	}
	return cfg, nil
}
