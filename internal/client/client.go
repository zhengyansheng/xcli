package client

import (
	"fmt"
	"net/http"
	"strings"

	"crypto/tls"

	"resty.dev/v3"
)

// APIResponse 定义通用的API响应结构
type APIResponse struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
	Count int         `json:"count"`
}

// LogonData 定义登录响应中的数据结构
type LogonData struct {
	Token string `json:"token"`
	Path  string `json:"path"`
	IP    string `json:"ip"`
}

// LogonResponse 定义登录响应的结构
type LogonResponse struct {
	Code  int       `json:"code"`
	Msg   string    `json:"msg"`
	Data  LogonData `json:"data"`
	Count int       `json:"count"`
}

// JobSubmitData 定义作业提交响应中的数据结构
type JobSubmitData struct {
	JobID   int64  `json:"jobid"`
	Message string `json:"message"`
}

// JobSubmitResponse 定义作业提交响应
type JobSubmitResponse struct {
	Code  int           `json:"code"`
	Msg   string        `json:"msg"`
	Data  JobSubmitData `json:"data"`
	Count int           `json:"count"`
}

// JobSubmitRequest 定义作业提交请求
type JobSubmitRequest struct {
	Queue   string `json:"queue"`
	ResReq  string `json:"resreq"`
	Command string `json:"command"`
}

// Host 定义主机信息结构
type Host struct {
	HostName       string    `json:"hostName"`
	HostType       string    `json:"hostType"`
	HostModel      string    `json:"hostModel"`
	CpuFactor      float64   `json:"cpuFactor"`
	MaxCpus        int       `json:"maxCpus"`
	MaxMem         int64     `json:"maxMem"`
	MaxSwap        int64     `json:"maxSwap"`
	MaxTmp         int64     `json:"maxTmp"`
	NDisks         int       `json:"nDisks"`
	NRes           int       `json:"nRes"`
	Resources      []string  `json:"resources"`
	NDRes          int       `json:"nDRes"`
	DResources     []string  `json:"DResources"`
	Windows        string    `json:"windows"`
	NumIndx        int       `json:"numIndx"`
	BusyThreshold  []float64 `json:"busyThreshold"`
	IsServer       bool      `json:"isServer"`
	Cores          int       `json:"cores"`
	HostAddr       string    `json:"hostAddr"`
	Pprocs         int       `json:"pprocs"`
	CoresPerProc   int       `json:"cores_per_proc"`
	ThreadsPerCore int       `json:"threads_per_core"`
}

// HostsResponse 定义主机查询响应
type HostsResponse struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	Data  []Host `json:"data"`
	Count int    `json:"count"`
}

// Job 定义作业信息结构
type Job struct {
	JobID          int64  `json:"jobid"`
	User           string `json:"user"`
	Status         string `json:"status"`
	JobName        string `json:"jobname"`
	Queue          string `json:"queue"`
	ProjectName    string `json:"projectname"`
	Command        string `json:"command"`
	ResReq         string `json:"resreq"`
	SubmitTime     string `json:"submittime"`
	JobDescription string `json:"jobdescription"`
}

// JobsResponse 定义作业查询响应
type JobsResponse struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	Data  []Job  `json:"data"`
	Count int    `json:"count"`
}

// APIClient 定义API客户端
type APIClient struct {
	client  *resty.Client
	baseURL string
}

// NewAPIClient 创建新的API客户端
func NewAPIClient(baseURL string) *APIClient {
	c := &APIClient{
		client:  resty.New(),
		baseURL: baseURL,
	}
	c.client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	return c
}

// Logon 执行登录操作
func (c *APIClient) Logon(username, password string) (*LogonResponse, error) {
	// c.client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	var loginResp LogonResponse
	resp, err := c.client.R().
		SetBody(map[string]interface{}{
			"username": username,
			"password": password,
		}).
		SetHeader("Content-Type", "application/json").
		SetResult(&loginResp).
		Post(c.baseURL)

	if err != nil {
		return nil, fmt.Errorf("登录失败: %v", err)
	}

	if resp.StatusCode() != 200 || loginResp.Code != 200 {
		return nil, fmt.Errorf("登录失败: %s", loginResp.Msg)
	}

	return &loginResp, nil
}

// Logout 执行登出操作
func (c *APIClient) Logout(token string) error {
	resp, err := c.client.R().
		SetHeader("Authorization", "Bearer "+token).
		Post(c.baseURL + "/xce/v1/auth/logout")

	if err != nil {
		return fmt.Errorf("登出请求失败: %v", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("登出失败: HTTP %d - %s", resp.StatusCode(), resp.String())
	}

	return nil
}

// SubmitJob 提交作业
func (c *APIClient) SubmitJob(token string, req *JobSubmitRequest) (*JobSubmitResponse, error) {
	var resp JobSubmitResponse

	// 从 baseURL 中提取 ip:port 部分
	baseURL := c.baseURL
	// 如果 URL 包含路径，去掉路径部分
	if idx := strings.Index(baseURL, "/xce/v1"); idx != -1 {
		baseURL = baseURL[:idx]
	}
	// 构建完整的作业提交 URL
	submitURL := baseURL + "/xce/v1/jobs"

	httpResp, err := c.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		SetResult(&resp).
		Post(submitURL)

	if err != nil {
		return nil, fmt.Errorf("提交作业请求失败: %v", err)
	}

	if httpResp.StatusCode() != http.StatusOK || resp.Code != http.StatusCreated {
		return nil, fmt.Errorf("提交作业失败: %s", resp.Msg)
	}

	return &resp, nil
}

// GetHosts 查询主机信息
func (c *APIClient) GetHosts(token string, params map[string]string) (*HostsResponse, error) {
	var resp HostsResponse

	// 从 baseURL 中提取 ip:port 部分
	baseURL := c.baseURL
	if idx := strings.Index(baseURL, "/xce/v1"); idx != -1 {
		baseURL = baseURL[:idx]
	}
	// 构建完整的主机查询 URL
	hostsURL := baseURL + "/xce/v1/hosts"

	req := c.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetResult(&resp)

	// 添加查询参数
	for k, v := range params {
		req.SetQueryParam(k, v)
	}

	httpResp, err := req.Get(hostsURL)
	if err != nil {
		return nil, fmt.Errorf("查询主机请求失败: %v", err)
	}

	if httpResp.StatusCode() != 200 || resp.Code != 200 {
		return nil, fmt.Errorf("查询主机失败: %s", resp.Msg)
	}

	return &resp, nil
}

// GetJobs 查询作业信息
func (c *APIClient) GetJobs(token string, params map[string]string) (*JobsResponse, error) {
	var resp JobsResponse

	// 从 baseURL 中提取 ip:port 部分
	baseURL := c.baseURL
	if idx := strings.Index(baseURL, "/xce/v1"); idx != -1 {
		baseURL = baseURL[:idx]
	}
	// 构建完整的作业查询 URL
	jobsURL := baseURL + "/xce/v1/jobs"

	req := c.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetResult(&resp)

	// 处理过滤条件
	if filter, ok := params["filter"]; ok {
		req.SetQueryParam("filter", filter)
	}

	// 处理字段选择
	if fields, ok := params["fields"]; ok {
		req.SetQueryParam("fields", fields)
	}

	httpResp, err := req.Get(jobsURL)
	if err != nil {
		return nil, fmt.Errorf("查询作业请求失败: %v", err)
	}

	if httpResp.StatusCode() != 200 || resp.Code != 200 {
		return nil, fmt.Errorf("查询作业失败: %s", resp.Msg)
	}

	return &resp, nil
}
