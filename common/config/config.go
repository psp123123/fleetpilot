package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	// 全局变量
	GlobalCfg *Config
)

// Config 结构体对应 config.yaml
type Config struct {
	Log   LogConfig
	Mysql MysqlConfig
}

// logger 相关配置内容
type LogConfig struct {
	Logger struct {
		Level string `yaml:"level"`
	} `yaml:"log"`
}

// mysql相关配置内容
type MysqlConfig struct {
	Mysqler struct {
		Address         string `yaml:"address"`
		Username        string `yaml:"username"`
		Password        string `yaml:"password"`
		Dbname          string `yaml:"dbname"`
		Timeout         int    `yaml:"timeout"`
		MultiStatements bool   `yaml:"multiStatements"`
	} `yaml:"mysql"`
}

// ensureConfigExists 检查并初始化配置文件
func EnsureConfigExists(path string) error {
	// 如果文件存在，直接返回
	if _, err := os.Stat(path); err == nil {
		return nil
	}

	fmt.Println("配置文件不存在，需要创建配置:", path)

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

// GetConfig 返回全局配置
func GetConfig() *Config {
	return GlobalCfg
}
