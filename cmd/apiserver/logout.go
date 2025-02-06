package apiserver

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xx/internal/client"
	"github.com/xx/pkg/config"
)

// NewLogoutCmd 创建登出命令
func NewLogoutCmd(configManager *config.ConfigManager) *cobra.Command {
	var url string

	cmd := &cobra.Command{
		Use:   "logout",
		Short: "从 APIserver 登出",
		Long: `从指定的 APIserver 登出，使当前 token 失效。
示例: 
  cli apiserver logout --url http://tt1.test.com:8080
  cli apiserver logout --url https://tt1.test.com:8443`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogout(url, configManager)
		},
	}

	// 添加 url 参数
	cmd.Flags().StringVar(&url, "url", "", "指定 APIserver 地址")
	cmd.MarkFlagRequired("url")

	return cmd
}

func runLogout(url string, cm *config.ConfigManager) error {
	// 获取配置
	cfg, err := cm.GetConfig()
	if err != nil {
		return fmt.Errorf("获取配置失败: %v", err)
	}

	// 查找指定 URL 的服务器
	var token string
	found := false
	for i, server := range cfg.APIServerInfo {
		if server.URL == url {
			token = server.Token
			// 清除 token
			cfg.APIServerInfo[i].Token = ""
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("未找到指定的服务器: %s", url)
	}

	// 如果有 token，调用登出接口
	if token != "" {
		apiClient := client.NewAPIClient(url)
		if err := apiClient.Logout(token); err != nil {
			return fmt.Errorf("登出失败: %v", err)
		}
	}

	// 保存配置
	if err := cm.SaveConfig(cfg); err != nil {
		return fmt.Errorf("保存配置失败: %v", err)
	}

	fmt.Println("登出成功")
	return nil
}
