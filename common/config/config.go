package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 结构体对应 config.yaml
type Config struct {
	Log struct {
		Level string `yaml:"level"`
	} `yaml:"log"`
}

// 默认配置内容
var defaultConfig = Config{
	Log: struct {
		Level string `yaml:"level"`
	}{Level: "info"},
}

// ensureConfigExists 检查并初始化配置文件
func EnsureConfigExists(path string) error {
	// 如果文件存在，直接返回
	if _, err := os.Stat(path); err == nil {
		return nil
	}

	fmt.Println("配置文件不存在，正在创建默认配置:", path)

	data, err := yaml.Marshal(defaultConfig)
	if err != nil {
		return fmt.Errorf("生成默认配置失败: %v", err)
	}

	// 创建并写入默认配置
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("写入默认配置文件失败: %v", err)
	}
	return nil
}

// loadConfig 读取配置文件
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	return &cfg, nil
}
