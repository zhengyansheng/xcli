package main

import (
	"k8s.io/klog/v2"
	"os"

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

	// 执行命令
	if err := cmd.Execute(configManager); err != nil {
		klog.Fatal(err)
		os.Exit(1)
	}
}
