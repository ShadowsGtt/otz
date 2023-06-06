package otz

import (
	"github.com/ShadowsGtt/otz/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sync/atomic"
)

// Config 本地配置
type Config struct {
	Ip   string     `yaml:"ip"`
	Port string     `yaml:"port"`
	Log  log.Config `yaml:"log"`
}

const (
	defaultConfigFile = "./otz_server.yaml"
)

var (
	GlobalServerConfigFile = defaultConfigFile
	globalServerConfig     atomic.Value
)

// LoadConfig 加载服务配置
func LoadConfig(filePath string) (*Config, error) {
	// 解析配置文件
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	err = yaml.Unmarshal(content, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// SetGlobalConfig 存储到全局变量
func SetGlobalConfig(cfg *Config) {
	globalServerConfig.Store(cfg)
}

// GetGlobalConfig 服务配置文件
func GetGlobalConfig() *Config {
	return globalServerConfig.Load().(*Config)
}
