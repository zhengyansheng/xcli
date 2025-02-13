package xsub

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xx/internal/client"
	"github.com/xx/pkg/config"
)

// NewXSubCmd 创建作业提交命令
func NewXSubCmd(configManager *config.ConfigManager) *cobra.Command {
	var (
		queue   string
		resReq  string
		command string
	)

	cmd := &cobra.Command{
		Use:   "xsub",
		Short: "提交作业到 APIserver",
		Long: `提交作业到 APIserver。
示例: 
  cli xsub -q q1 -R "select(!mg)" sleep 10`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("请提供要执行的命令")
			}
			// 将参数组合成命令字符串
			// command := strings.Join(args, " ")
			return runXSub(configManager, queue, resReq, command)
		},
	}

	// 添加命令行参数
	flags := cmd.Flags()
	flags.StringVarP(&queue, "queue", "q", "", "指定作业队列")
	flags.StringVarP(&resReq, "resreq", "R", "", "指定资源需求")
	flags.StringVarP(&command, "command", "c", "", "指定要执行的命令")
	// 设置必需参数
	//cmd.MarkFlagRequired("queue")

	return cmd
}

func runXSub(cm *config.ConfigManager, queue, resReq, command string) error {
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
	// fmt.Println(serverInfo)

	if serverInfo == nil {
		return fmt.Errorf("未找到默认服务器信息")
	}

	if serverInfo.Token == "" {
		return fmt.Errorf("未登录到服务器，请先登录")
	}

	// 创建作业提交请求
	jobReq := &client.JobSubmitRequest{
		Queue:   queue,
		ResReq:  resReq,
		Command: command,
	}
	// fmt.Println(jobReq)
	// fmt.Println(serverInfo.Token)
	// fmt.Println(cfg.DefaultAPIServer)

	// 创建 API 客户端并提交作业
	apiClient := client.NewAPIClient(cfg.DefaultAPIServer)
	jobResp, err := apiClient.SubmitJob(serverInfo.Token, jobReq)
	if err != nil {
		return fmt.Errorf("提交作业失败: %v", err)
	}

	fmt.Printf("作业提交成功，作业ID: %d\n%s\n", jobResp.Data.JobID, jobResp.Data.Message)
	return nil
}
