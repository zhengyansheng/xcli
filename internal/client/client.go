package client

import (
	"fmt"

	"resty.dev/v3"
)

// LogonResponse 定义登录响应的结构
type LogonResponse struct {
	Token   string `json:"token"`
	Version string `json:"version"`
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
		SetResult(&loginResp). // 设置结果结构体
		Post(c.baseURL + "/xce/v1/auth/logon")

	if err != nil {
		return nil, fmt.Errorf("登录失败: %v", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("登录失败: HTTP %d - %s", resp.StatusCode(), resp.String())
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
