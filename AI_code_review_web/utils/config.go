package utils

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config 表示配置文件结构
type Config struct {
	MaxThreads  int          `yaml:"max_threads"`
	AIPlatforms []AIPlatform `yaml:"ai_platforms"`
}

// AIPlatform 表示一个AI平台的配置信息
type AIPlatform struct {
	Name     string `yaml:"name"`
	Endpoint string `yaml:"endpoint"`
	Model    string `yaml:"model"`
	APIKey   string `yaml:"api_key"`
}

// LoadConfig 加载配置文件
func LoadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
