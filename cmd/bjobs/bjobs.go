package bjobs

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

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
		Long: `查询作业信息，支持按用户和队列过滤。
示例:
  cli bjobs                                # 查询所有作业
  cli bjobs -u user1                      # 查询指定用户的作业
  cli bjobs -q queue1 -u user1            # 查询指定用户在指定队列的作业
  cli bjobs jobid,status,queue,command    # 查询指定字段`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 如果有位置参数，作为字段列表
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
	params := make(map[string]string)

	// 处理过滤条件
	var filters []string
	if user != "" {
		filters = append(filters, fmt.Sprintf("user:eq:%s", user))
	}
	if queue != "" {
		filters = append(filters, fmt.Sprintf("queue:eq:%s", queue))
	}
	if len(filters) > 0 {
		params["filter"] = fmt.Sprintf("[%s]", strings.Join(filters, ","))
	}

	// 处理字段选择
	if fields != "" {
		params["fields"] = fields
	}

	// 创建 API 客户端并查询作业信息
	apiClient := client.NewAPIClient(cfg.DefaultAPIServer)
	jobs, err := apiClient.GetJobs(serverInfo.Token, params)
	if err != nil {
		return fmt.Errorf("查询作业失败: %v", err)
	}

	// 显示结果
	printJobs(jobs)
	return nil
}

func printJobs(jobs *client.JobsResponse) {
	if jobs.Count == 0 {
		fmt.Println("没有找到作业")
		return
	}

	// 打印表头
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "JOBID\tUSER\tSTATUS\tQUEUE\tCOMMAND")

	// 打印作业信息
	for _, job := range jobs.Data {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
			job.JobID,
			job.User,
			job.Status,
			job.Queue,
			job.Command)
	}
	w.Flush()
}
