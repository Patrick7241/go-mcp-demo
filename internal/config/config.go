package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Ollama OllamaConfig `yaml:"ollama"`
	Mysql  MysqlConfig  `yaml:"mysql"`
	Sql    SqlConfig    `yaml:"sql"`
}

type OllamaConfig struct {
	Model string `yaml:"model"`
}

type MysqlConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type SqlConfig struct {
	SqlFilePath string `yaml:"sqlFilePath"`
}

var AppConfig *Config

func LoadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	AppConfig = &cfg
	return nil
}
