package bhosts

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/xx/internal/client"
	"github.com/xx/pkg/config"
)

// NewBHostsCmd 创建主机查询命令
func NewBHostsCmd(configManager *config.ConfigManager) *cobra.Command {
	var (
		infoType string
		hostType string
		fullInfo bool
	)

	cmd := &cobra.Command{
		Use:   "bhosts",
		Short: "查询主机信息",
		Long: `查询主机信息，支持基本信息和详细信息。
示例:
  cli bhosts                    # 查询基本信息
  cli bhosts --type full        # 查询详细信息
  cli bhosts --host-type X86_64 # 按主机类型过滤`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 处理 --full 标志
			if fullInfo {
				infoType = "full"
			}
			return runBHosts(configManager, infoType, hostType)
		},
	}

	// 添加命令行参数
	flags := cmd.Flags()
	flags.StringVar(&infoType, "type", "basic", "信息类型 (basic/full)")
	flags.StringVar(&hostType, "host-type", "", "主机类型过滤 (X86_64/ARM)")
	flags.BoolVar(&fullInfo, "full", false, "显示详细信息")

	return cmd
}

func runBHosts(cm *config.ConfigManager, infoType, hostType string) error {
	// 获取配置
	cfg, err := cm.GetConfig()
	if err != nil {
		return fmt.Errorf("获取配置失败: %v", err)
	}

	// 检查是否有默认服务器
	if cfg.DefaultAPIServer == "" {
		return fmt.Errorf("未设置默认 APIserver，请先设置默认服务器或指定服务器")
	}

	// 获取服务器信息
	var serverInfo *config.APIServerInfo
	for _, server := range cfg.APIServerInfo {
		if server.URL == cfg.DefaultAPIServer {
			serverInfo = &server
			break
		}
	}

	if serverInfo == nil {
		return fmt.Errorf("未找到默认服务器信息")
	}

	if serverInfo.Token == "" {
		return fmt.Errorf("未登录到服务器，请先登录")
	}

	// 验证参数
	infoType = strings.ToLower(infoType)
	if infoType != "basic" && infoType != "full" {
		return fmt.Errorf("无效的信息类型: %s，必须是 basic 或 full", infoType)
	}

	// 创建查询参数
	queryParams := make(map[string]string)
	queryParams["type"] = infoType

	if hostType != "" {
		hostType = strings.ToUpper(hostType)
		if hostType != "X86_64" && hostType != "ARM" {
			return fmt.Errorf("无效的主机类型: %s，必须是 X86_64 或 ARM", hostType)
		}
		queryParams["filter"] = fmt.Sprintf("hostType:eq:%s", hostType)
	}

	// 创建 API 客户端并查询主机信息
	apiClient := client.NewAPIClient(cfg.DefaultAPIServer)
	hosts, err := apiClient.GetHosts(serverInfo.Token, queryParams)
	if err != nil {
		return fmt.Errorf("查询主机信息失败: %v", err)
	}

	// 显示结果
	printHosts(hosts)
	return nil
}

func printHosts(hosts *client.HostsResponse) {
	if hosts.Count == 0 {
		fmt.Println("没有找到主机")
		return
	}

	// 打印表头
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "HOST_NAME\tTYPE\tMODEL\tCPUS\tMEM(MB)\tSWAP(MB)\tRESOURCES")

	// 打印主机信息
	for _, host := range hosts.Data {
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%d\t%d\t%s\n",
			host.HostName,
			host.HostType,
			host.HostModel,
			host.MaxCpus,
			host.MaxMem/1024,  // 转换为MB
			host.MaxSwap/1024, // 转换为MB
			strings.Join(append(host.Resources, host.DResources...), " "))
	}
	w.Flush()
}
