package main

import (
	"os"

	"k8s.io/klog/v2"

	"github.com/xx/cmd"
	"github.com/xx/pkg/config"
)

func main() {
	// 初始化配置管理器
	configManager, err := config.NewConfigManager()
	if err != nil {
		klog.Errorf("初始化配置失败: %v\n", err)
		os.Exit(1)
	}

	// 确保配置管理器被正确初始化
	if configManager == nil {
		klog.Error("配置管理器初始化失败")
		os.Exit(1)
	}

	// 执行命令
	if err := cmd.Execute(configManager); err != nil {
		klog.Fatal(err)
		os.Exit(1)
	}
}
