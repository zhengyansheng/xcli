package apiserver

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/xx/pkg/config"
)

// NewListCmd 创建 list 命令
func NewListCmd(configManager *config.ConfigManager) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "列出所有 APIserver",
		Long:  "显示所有已配置的 APIserver 和默认 APIserver",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(configManager)
		},
	}
}

func runList(cm *config.ConfigManager) error {
	cfg, err := cm.GetConfig()
	if err != nil {
		return fmt.Errorf("获取配置失败: %v", err)
	}

	// 使用 tabwriter 格式化输出
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "Default\tName\tVersion\tURL")

	for _, server := range cfg.APIServerInfo {
		isDefault := " "
		if server.URL == cfg.DefaultAPIServer {
			isDefault = "*"
		}
		if server.Token == "" {
			continue
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			isDefault,
			server.Name,
			server.Version,
			server.URL)
	}

	return w.Flush()
}
