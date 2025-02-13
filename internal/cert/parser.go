package cert

import (
	"fmt"
	"net/url"
)

func extractHostAndPort(rawURL string) (string, error) {
	// 解析 URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("解析 URL 失败: %v", err)
	}

	// 获取主机和端口
	host := parsedURL.Hostname()
	port := parsedURL.Port()

	// 如果端口为空，使用默认端口（HTTP 为 80，HTTPS 为 443）
	if port == "" {
		if parsedURL.Scheme == "https" {
			port = "443"
		} else if parsedURL.Scheme == "http" {
			port = "80"
		}
	}

	// 拼接主机和端口
	result := host
	if port != "" {
		result = fmt.Sprintf("%s:%s", host, port)
	}

	return result, nil
}
