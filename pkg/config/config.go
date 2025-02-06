package config

import (
	"encoding/json"
	"os"
	"os/user"
	"path"
)

type APIServerInfo struct {
	URL          string `json:"url"`
	Name         string `json:"name"`
	JobIDRange   string `json:"jobid_range"`
	ClusterIndex string `json:"clusterIndex"`
	Version      string `json:"version"`
	Token        string `json:"token"`
}

type Config struct {
	Account          string          `json:"account"`
	DefaultAPIServer string          `json:"defaultAPIserver"`
	DefaultQueryAll  bool            `json:"defaultqueryall"`
	CACert           string          `json:"cacert"`
	APIServerInfo    []APIServerInfo `json:"APIserverInfo"`
}

// ConfigManager 用于统一管理配置
type ConfigManager struct {
	config     *Config
	configPath string
}

// NewConfigManager 创建配置管理器
func NewConfigManager() (*ConfigManager, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	configDir := path.Join(usr.HomeDir, ".cli")
	configPath := path.Join(configDir, "config.json")

	return &ConfigManager{
		configPath: configPath,
	}, nil
}

// GetConfig 获取配置
func (cm *ConfigManager) GetConfig() (*Config, error) {
	if cm.config != nil {
		return cm.config, nil
	}

	config, err := cm.loadConfig()
	if err != nil {
		return nil, err
	}

	cm.config = config
	return config, nil
}

// SaveConfig 保存配置
func (cm *ConfigManager) SaveConfig(config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	dir := path.Dir(cm.configPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	if err := os.WriteFile(cm.configPath, data, 0600); err != nil {
		return err
	}

	cm.config = config
	return nil
}

// loadConfig 从文件加载配置
func (cm *ConfigManager) loadConfig() (*Config, error) {
	// 如果配置文件不存在，创建默认配置
	if _, err := os.Stat(cm.configPath); os.IsNotExist(err) {
		return &Config{
			APIServerInfo: make([]APIServerInfo, 0),
		}, nil
	}

	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) Save() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	configDir := path.Join(usr.HomeDir, ".cli")
	configPath := path.Join(configDir, "config.json")

	// 确保目录存在
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0600)
}
