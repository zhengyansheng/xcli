package config

import (
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/xx/pkg/config"
)

// NewSetCmd 创建配置设置命令
func NewSetCmd(configManager *config.ConfigManager) *cobra.Command {
	var (
		defaultAPIServer string
		defaultQueryAll  string
		caCert           string
	)

	cmd := &cobra.Command{
		Use:   "set",
		Short: "设置 APIserver 配置文件中的值",
		Long: `设置 APIserver 配置文件中的个别值。
示例:
  cli config set --defaultapiserver http://tt1.test.com:8080  # 设置默认 APIserver
  cli config set --defaultqueryall y                          # 设置默认查询所有
  cli config set --cacert /usr/cacert.pem                     # 设置证书路径`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSet(configManager, defaultAPIServer, defaultQueryAll, caCert)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&defaultAPIServer, "defaultapiserver", "", "设置默认的 APIserver")
	flags.StringVar(&defaultQueryAll, "defaultqueryall", "", "设置是否默认查询所有 (y/n)")
	flags.StringVar(&caCert, "cacert", "", "设置用于验证 APIserver 的 CA 证书路径")

	return cmd
}

func runSet(cm *config.ConfigManager, defaultAPIServer, defaultQueryAll, caCert string) error {
	cfg, err := cm.GetConfig()
	if err != nil {
		return fmt.Errorf("获取配置失败: %v", err)
	}

	// 设置默认 APIserver
	if defaultAPIServer != "" {
		// 验证 URL 格式
		if !isValidURL(defaultAPIServer) {
			return fmt.Errorf("无效的 APIserver URL: %s", defaultAPIServer)
		}

		found := false
		for _, server := range cfg.APIServerInfo {
			if server.URL == defaultAPIServer {
				found = true
				cfg.DefaultAPIServer = defaultAPIServer
				break
			}
		}
		if !found {
			return fmt.Errorf("未找到指定的 APIserver: %s", defaultAPIServer)
		}
	}

	// 设置默认查询所有
	if defaultQueryAll != "" {
		switch defaultQueryAll {
		case "y", "Y":
			cfg.DefaultQueryAll = true
		case "n", "N":
			cfg.DefaultQueryAll = false
		default:
			return fmt.Errorf("defaultqueryall 参数无效，请使用 y 或 n")
		}
	}

	// 设置 CA 证书
	if caCert != "" {
		// 检查证书文件是否存在
		if _, err := os.Stat(caCert); err != nil {
			return fmt.Errorf("证书文件不存在: %s", caCert)
		}
		cfg.CACert = caCert
	}

	// 保存配置
	if err := cm.SaveConfig(cfg); err != nil {
		return fmt.Errorf("保存配置失败: %v", err)
	}

	fmt.Println("配置已更新")
	return nil
}

// isValidURL 验证 URL 格式
func isValidURL(url string) bool {
	re := regexp.MustCompile(`^http[s]?://[a-zA-Z0-9.-]+(:[0-9]+)?(/.*)?$`)
	return re.MatchString(url)
}
