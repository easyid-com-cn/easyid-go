# EasyID Go SDK

Official Go SDK for the EasyID identity verification API.

EasyID 易验云 focuses on identity verification and security risk control APIs, including real-name verification, liveness detection, face recognition, phone verification, and fraud-risk related capabilities.

Chinese documentation: [README.zh-CN.md](README.zh-CN.md)

## Install

```bash
go get github.com/easyid-com-cn/easyid-go@latest
```

## Quick Start

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

    fmt.Println(result.Match)
}
```

## Supported APIs

- IDCard: `Verify2`, `Verify3`, `OCR`
- Phone: `Status`, `Verify3`
- Face: `Liveness`, `Compare`, `Verify`
- Bank: `Verify4`
- Risk: `Score`, `StoreFingerprint`
- Billing: `Balance`, `Records`

## Configuration

- `easyid.WithBaseURL(...)`
- `easyid.WithTimeout(...)`
- `easyid.WithHTTPClient(...)`

## Error Handling

Service-side business errors are returned as `*easyid.APIError`.

```go
if apiErr, ok := easyid.IsAPIError(err); ok {
    fmt.Println(apiErr.Code, apiErr.Message)
}
```

## Security Notice

This is a server-side SDK. Never expose `secret` in browsers, mobile apps, or other untrusted clients.

## Official Resources

- Official website: `https://www.easyid.com.cn/`
- GitHub organization: `https://github.com/easyid-com-cn/`
