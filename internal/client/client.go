package client

import (
	"fmt"

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

// JobSubmitResponse 定义作业提交响应
type JobSubmitResponse struct {
	JobID string `json:"jobid"`
}

// JobSubmitRequest 定义作业提交请求
type JobSubmitRequest struct {
	Queue   string `json:"queue"`
	ResReq  string `json:"resreq"`
	Command string `json:"command"`
}

// HostsResponse 定义主机查询响应
type HostsResponse struct {
	// TODO: 根据实际API响应定义字段
	Hosts []interface{} `json:"hosts"`
}

// Job 定义作业信息
type Job struct {
	JobID   int    `json:"jobid"`
	Status  string `json:"status"`
	Queue   string `json:"queue"`
	Command string `json:"command"`
}

// JobsResponse 定义作业查询响应
type JobsResponse struct {
	Jobs []Job `json:"jobs"`
}

// APIClient 定义API客户端
type APIClient struct {
	client  *resty.Client
	baseURL string
}

// NewAPIClient 创建新的API客户端
func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		client:  resty.New(),
		baseURL: baseURL,
	}
}

// Logon 执行登录操作
func (c *APIClient) Logon(username, password string) (*LogonResponse, error) {
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
	httpResp, err := c.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetBody(req).
		SetResult(&resp).
		Post(c.baseURL + "/xce/v1/jobs")

	if err != nil {
		return nil, fmt.Errorf("提交作业请求失败: %v", err)
	}

	if httpResp.StatusCode() != 200 {
		return nil, fmt.Errorf("提交作业失败: HTTP %d - %s", httpResp.StatusCode(), httpResp.String())
	}

	return &resp, nil
}

// GetHosts 查询主机信息
func (c *APIClient) GetHosts(token string, params map[string]string) (*HostsResponse, error) {
	var resp HostsResponse
	req := c.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetResult(&resp)

	// 添加查询参数
	for k, v := range params {
		req.SetQueryParam(k, v)
	}

	httpResp, err := req.Get(c.baseURL + "/xce/v1/hosts")
	if err != nil {
		return nil, fmt.Errorf("查询主机请求失败: %v", err)
	}

	if httpResp.StatusCode() != 200 {
		return nil, fmt.Errorf("查询主机失败: HTTP %d - %s", httpResp.StatusCode(), httpResp.String())
	}

	return &resp, nil
}

// GetJobs 查询作业信息
func (c *APIClient) GetJobs(token string, params map[string]string) (*JobsResponse, error) {
	var resp JobsResponse
	req := c.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetResult(&resp)

	// 添加查询参数
	for k, v := range params {
		req.SetQueryParam(k, v)
	}

	httpResp, err := req.Get(c.baseURL + "/xce/v1/jobs")
	if err != nil {
		return nil, fmt.Errorf("查询作业请求失败: %v", err)
	}

	if httpResp.StatusCode() != 200 {
		return nil, fmt.Errorf("查询作业失败: HTTP %d - %s", httpResp.StatusCode(), httpResp.String())
	}

	return &resp, nil
}
