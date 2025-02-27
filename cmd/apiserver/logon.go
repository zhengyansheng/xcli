package apiserver

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xx/internal/cert"
	"github.com/xx/internal/client"
	"github.com/xx/pkg/config"
)

// updateServerConfig 更新服务器配置
func updateServerConfig(cfg *config.Config, url, username string, loginResp *client.LogonResponse) error {
	// 检查是否已存在相同URL的服务器
	found := false
	for i, server := range cfg.APIServerInfo {
		if server.URL == url {
			// 更新现有服务器信息
			cfg.APIServerInfo[i].Token = loginResp.Data.Token
			cfg.APIServerInfo[i].Path = loginResp.Data.Path // 保存路径
			found = true
			break
		}
	}

	// 如果是新服务器，添加到列表
	if !found {
		// 生成新服务器名称
		serverName := fmt.Sprintf("apiserver%d", len(cfg.APIServerInfo)+1)

		newServer := config.APIServerInfo{
			Name:  serverName,
			URL:   url,
			Token: loginResp.Data.Token,
			Path:  loginResp.Data.Path, // 保存路径
		}

		// 添加到数组末尾
		cfg.APIServerInfo = append(cfg.APIServerInfo, newServer)
	}

	// 如果是第一个服务器，设置为默认服务器
	if len(cfg.APIServerInfo) == 1 {
		cfg.DefaultAPIServer = url
	}

	// 设置当前用户账号
	cfg.Account = username

	return nil
}

func NewLogonCmd(configManager *config.ConfigManager) *cobra.Command {
	var opts struct {
		username string
		password string
		url      string
	}

	cmd := &cobra.Command{
		Use:   "logon",
		Short: "登录到 APIserver",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogon(opts, configManager)
		},
	}

	// 设置命令行参数
	flags := cmd.Flags()
	flags.StringVarP(&opts.username, "username", "n", "", "用户名")
	flags.StringVarP(&opts.password, "password", "p", "", "密码")
	flags.StringVar(&opts.url, "url", "", "APIserver 地址")

	// 必填参数
	cmd.MarkFlagRequired("username")
	cmd.MarkFlagRequired("password")
	cmd.MarkFlagRequired("url")

	return cmd
}

func runLogon(opts struct{ username, password, url string }, cm *config.ConfigManager) error {
	if cm == nil {
		return fmt.Errorf("配置管理器未初始化")
	}

	// 获取配置
	cfg, err := cm.GetConfig()
	if err != nil {
		return fmt.Errorf("获取配置失败: %v", err)
	}

	// 创建 API 客户端
	apiClient := client.NewAPIClient(opts.url)
	apiClient.SetRootCAs(cfg.CACert)

	certPath, err := cert.GetCertPath()
	if err != nil {
		return fmt.Errorf("获取证书路径失败: %v", err)
	}
	err = cert.GeneratorCert(opts.url, certPath)
	if err != nil {
		return fmt.Errorf("生成证书失败: %v", err)
	}

	// 执行登录
	loginResp, err := apiClient.Logon(opts.username, opts.password)
	if err != nil {
		return fmt.Errorf("登录失败: %v", err)
	}

	// 更新服务器信息
	if err := updateServerConfig(cfg, opts.url, opts.username, loginResp); err != nil {
		return fmt.Errorf("更新服务器配置失败: %v", err)
	}

	cfg.CACert = certPath
	fmt.Println("证书路径: ", certPath)

	// 保存配置
	if err := cm.SaveConfig(cfg); err != nil {
		return fmt.Errorf("保存配置失败: %v", err)
	}

	fmt.Println("登录成功")
	return nil
}
