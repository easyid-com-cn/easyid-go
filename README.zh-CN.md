# EasyID Go SDK

EasyID Go SDK 是易验云身份验证 API 的官方 Go 客户端。

English README: [README.md](README.md)

EasyID 提供身份证核验、手机号核验、人脸识别、银行卡核验、风控评分等能力。本 SDK 按业务模块分组，自动处理请求签名、认证头和错误解析，适合服务端集成。

## 安装

```bash
go get github.com/easyid-com-cn/easyid-go@latest
```

国内网络环境可先配置：

```bash
go env -w GOPROXY=https://goproxy.cn,direct
```

## 快速开始

```go
package main

import (
	"context"
	"fmt"

	"github.com/easyid-com-cn/easyid-go"
)

func main() {
	client := easyid.New("ak_xxx", "sk_xxx")

	result, err := client.IDCard.Verify2(context.Background(), &easyid.IDCardVerify2Request{
		Name:     "张三",
		IDNumber: "110101199001011234",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("是否匹配：", result.Match)
}
```

## 已支持接口

- `client.IDCard.Verify2`：身份证二要素核验
- `client.IDCard.Verify3`：身份证三要素核验
- `client.IDCard.OCR`：身份证 OCR
- `client.Phone.Status`：手机号状态查询
- `client.Phone.Verify3`：手机号三要素核验
- `client.Face.Liveness`：人脸活体检测
- `client.Face.Compare`：人脸比对
- `client.Face.Verify`：人脸核验
- `client.Bank.Verify4`：银行卡四要素核验
- `client.Risk.Score`：风控评分
- `client.Risk.StoreFingerprint`：存储设备指纹
- `client.Billing.Balance`：查询账户余额
- `client.Billing.Records`：查询账单记录

## 配置项

- `easyid.WithBaseURL(...)`：覆盖 API 地址，适合私有部署
- `easyid.WithTimeout(...)`：覆盖超时时间
- `easyid.WithHTTPClient(...)`：传入自定义 `http.Client`

## 错误处理

服务端业务错误会返回 `*easyid.APIError`。

```go
if apiErr, ok := easyid.IsAPIError(err); ok {
	fmt.Println(apiErr.Code, apiErr.Message, apiErr.RequestID)
}
```

## 安全说明

- 这是服务端 SDK，不要在浏览器、小程序、移动端或其他不可信客户端暴露 `secret`
- `keyID` 必须符合 `ak_[0-9a-f]+`
- SDK 会自动处理签名和认证请求头

## 官方资源

- 官网：`https://www.easyid.com.cn/`
- GitHub：`https://github.com/easyid-com-cn/`
