package yhttp

import (
	"crypto/tls"
	"time"

	resty "github.com/go-resty/resty/v2"
)

const (
	UAHeader  = `User-Agent`
	UAContent = `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36 Edg/130.0.2849.56 Azeroth-SHA`
)

type (
	Client   = resty.Client
	Request  = resty.Request
	Response = resty.Response
)

// New 创建并配置一个新的 HTTP 客户端实例
// 返回值为 *Client 类型，表示配置好的 HTTP 客户端
func New() *Client {
	cli := resty.New()
	cli = cli.SetTimeout(time.Minute)                                    // 设置请求超时时间为 1 分钟
	cli = cli.SetHeader(UAHeader, UAContent)                             // 设置默认 User-Agent 头
	return cli.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // 配置 TLS 跳过证书验证
}
