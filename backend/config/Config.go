package config

import (
	"card-authorization/log"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

var SystemConfig Config

type Config struct {
	HTTPPort    string `yaml:"http_port"`
	EmailConfig struct {
		SMTPHost     string `yaml:"smtp_host"`
		SMTPPort     int    `yaml:"smtp_port"` // 改为 int 以匹配 YAML
		AuthEmail    string `yaml:"auth_email"`
		AuthPassword string `yaml:"auth_password"`
	} `yaml:"email"`
}

// LoadConfig 加载外部配置文件
// LoadConfig 读取 YAML 配置文件
func LoadConfig() error {
	// 读取文件
	data, err := os.ReadFile("./config.yaml")
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析 YAML
	if err := yaml.Unmarshal(data, &SystemConfig); err != nil {
		return fmt.Errorf("解析 YAML 失败: %w", err)
	}
	log.Info("HTTP Port: %s", SystemConfig.HTTPPort)
	log.Info("SMTP Host: %s", SystemConfig.EmailConfig.SMTPHost)
	log.Info("SMTP Port: %d", SystemConfig.EmailConfig.SMTPPort)
	log.Info("Auth Email: %s", SystemConfig.EmailConfig.AuthEmail)
	log.Info("Auth Password: %s", "******")
	return nil
}
