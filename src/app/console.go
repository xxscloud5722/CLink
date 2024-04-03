package app

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Kafka struct {
		Server   string `yaml:"server"` // Kafka 地址.
		Port     int    `yaml:"port"`   // Kafka 端口.
		Consumer struct {
			GroupId string   `yaml:"group-id"` // Kafka 消费者分组ID.
			Topic   []string `yaml:"topic"`    // Kafka 消费则主题.
		} `json:"consumer"`
	} `yaml:"kafka"`
	Clickhouse struct {
		Server   string            `yaml:"host"`     // Clickhouse 地址.
		Port     int               `yaml:"port"`     // Clickhouse 端口.
		Username string            `yaml:"username"` // Clickhouse 用户名称.
		Password string            `yaml:"password"` // Clickhouse 用户密码.
		Database string            `yaml:"database"` // Clickhouse 数据库名称.
		Table    string            `yaml:"table"`    // Clickhouse 表名称.
		Fields   map[string]string `yaml:"fields"`   // Clickhouse 数据库字段关系对应表.
	} `yaml:"clickhouse"`
	Pattern       string   `yaml:"pattern"`       // 正则过滤器.
	PatternFields []string `yaml:"pattern-index"` // 正则解析的字段名称.
	debug         bool     // 是否启用调试模式.
}

type LogMessage struct {
	body string
}

func Run(configPath *string, debug bool) error {
	yamlFile, err := os.ReadFile(*configPath)
	if err != nil {
		return errors.New(fmt.Sprintf("[Config] Error reading YAML file: %v", err))
	}
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return errors.New(fmt.Sprintf("[Config] Error parse YAML file:  %v", err))
	}
	config.debug = debug

	// Create receiver.

	// Create filters.

	// Create a transmitter.

	return nil
}
