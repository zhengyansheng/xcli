package main

import (
	"fmt"
	"os"

	"github.com/xx/cmd"
	"github.com/xx/pkg/config"
)

func main() {
	// 初始化配置管理器
	configManager, err := config.NewConfigManager()
	if err != nil {
		fmt.Fprintf(os.Stderr, "初始化配置失败: %v\n", err)
		os.Exit(1)
	}

	// 执行命令
	if err := cmd.Execute(configManager); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
