package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path"
)

type APIServerInfo struct {
	Name         string `json:"name"`
	URL          string `json:"url"`
	Token        string `json:"token"`
	Path         string `json:"path"`
	JobIDRange   string `json:"jobid_range,omitempty"`
	ClusterIndex string `json:"cluster_index,omitempty"`
	Version      string `json:"version,omitempty"`
}

type Config struct {
	Account          string          `json:"account"`
	DefaultAPIServer string          `json:"defaultAPIserver"`
	DefaultQueryAll  bool            `json:"defaultqueryall"`
	CACert           string          `json:"cacert"`
	APIServerInfo    []APIServerInfo `json:"servers"`
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
		return nil, fmt.Errorf("获取用户信息失败: %v", err)
	}

	configDir := path.Join(usr.HomeDir, ".cli")
	configPath := path.Join(configDir, "config.json")

	// 确保配置目录存在
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return nil, fmt.Errorf("创建配置目录失败: %v", err)
	}

	cm := &ConfigManager{
		configPath: configPath,
	}

	// 初始化时加载配置
	_, err = cm.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %v", err)
	}

	return cm, nil
}

// GetConfig 获取配置
func (cm *ConfigManager) GetConfig() (*Config, error) {
	if cm == nil {
		return nil, fmt.Errorf("配置管理器未初始化")
	}

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
		defaultConfig := &Config{
			APIServerInfo: make([]APIServerInfo, 0),
		}
		// 保存默认配置
		if err := cm.SaveConfig(defaultConfig); err != nil {
			return nil, fmt.Errorf("创建默认配置失败: %v", err)
		}
		return defaultConfig, nil
	}

	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
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
