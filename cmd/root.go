package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xx/cmd/apiserver"
	"github.com/xx/cmd/bhosts"
	"github.com/xx/cmd/bjobs"
	setConfig "github.com/xx/cmd/config"
	"github.com/xx/cmd/xsub"
	"github.com/xx/pkg/config"
)

var (
	rootCmd = &cobra.Command{
		Use:   "cli",
		Short: "CLI tool for APIserver operations",
	}
	configManager *config.ConfigManager
)

// Execute 执行根命令
func Execute(cm *config.ConfigManager) error {
	if cm == nil {
		return fmt.Errorf("配置管理器未初始化")
	}
	configManager = cm
	cobra.EnableCommandSorting = false
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// 初始化子命令
	initCommands()

	return rootCmd.Execute()
}

// initCommands 初始化所有子命令
func initCommands() {

	// 添加子命令
	rootCmd.AddCommand(getAPIServerCmd())
	rootCmd.AddCommand(getConfigCmd())
	rootCmd.AddCommand(bjobs.NewBJobsCmd(configManager))
	rootCmd.AddCommand(bhosts.NewBHostsCmd(configManager))
	rootCmd.AddCommand(xsub.NewXSubCmd(configManager))

}

// getConfigCmd 返回配置子命令
func getConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "配置管理命令",
	}

	// 添加配置相关子命令
	cmd.AddCommand(setConfig.NewSetCmd(configManager))
	return cmd
}

// getAPIServerCmd 返回 apiserver 子命令
func getAPIServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apiserver",
		Short: "APIserver 相关操作",
	}

	// 禁用自动生成标签（这通常用于文档生成，但可以尝试以防止排序）
	cmd.DisableAutoGenTag = true

	// 手动设置子命令的添加顺序
	cmd.AddCommand(apiserver.NewLogonCmd(configManager))
	cmd.AddCommand(apiserver.NewLogoutCmd(configManager))
	cmd.AddCommand(apiserver.NewListCmd(configManager))

	return cmd
}
