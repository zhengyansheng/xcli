package bjobs

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xx/internal/client"
	"github.com/xx/pkg/config"
)

// NewBJobsCmd 创建作业查询命令
func NewBJobsCmd(configManager *config.ConfigManager) *cobra.Command {
	var (
		user   string
		queue  string
		fields string
	)

	cmd := &cobra.Command{
		Use:   "bjobs",
		Short: "查询作业信息",
		Long: `查询作业信息，支持按用户和队列过滤，以及自定义显示字段。
示例:
  cli bjobs                                    # 查询所有作业
  cli bjobs -u user1                          # 查询指定用户的作业
  cli bjobs -q queue1 -u user1                # 查询指定用户在指定队列的作业
  cli bjobs jobid,status,queue,command        # 指定显示字段`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 处理自定义字段
			if len(args) > 0 {
				fields = args[0]
			}
			return runBJobs(configManager, user, queue, fields)
		},
	}

	// 添加命令行参数
	flags := cmd.Flags()
	flags.StringVarP(&user, "user", "u", "", "按用户过滤")
	flags.StringVarP(&queue, "queue", "q", "", "按队列过滤")

	return cmd
}

func runBJobs(cm *config.ConfigManager, user, queue, fields string) error {
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

	// 构建查询参数
	queryParams := make(map[string]string)

	// 处理过滤条件
	var filters []string
	if user != "" {
		filters = append(filters, fmt.Sprintf("user:eq:%s", user))
	}
	if queue != "" {
		filters = append(filters, fmt.Sprintf("queue:eq:%s", queue))
	}
	if len(filters) > 0 {
		queryParams["filter"] = fmt.Sprintf("[%s]", strings.Join(filters, ","))
	}

	// 处理字段选择
	if fields != "" {
		queryParams["fields"] = fields
	}

	// 创建 API 客户端并查询作业信息
	apiClient := client.NewAPIClient(cfg.DefaultAPIServer)
	jobs, err := apiClient.GetJobs(serverInfo.Token, queryParams)
	if err != nil {
		return fmt.Errorf("查询作业失败: %v", err)
	}

	// 显示结果
	printJobs(jobs)
	return nil
}

func printJobs(jobs *client.JobsResponse) {
	if len(jobs.Jobs) == 0 {
		fmt.Println("没有找到作业")
		return
	}

	// 打印表头
	fmt.Printf("%-10s %-10s %-15s %s\n", "JOBID", "STATUS", "QUEUE", "COMMAND")

	// 打印作业信息
	for _, job := range jobs.Jobs {
		fmt.Printf("%-10d %-10s %-15s %s\n",
			job.JobID,
			job.Status,
			job.Queue,
			job.Command)
	}
}
