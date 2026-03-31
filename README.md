# EasyID Go SDK

Official Go SDK for the EasyID identity verification API.

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

## More Docs

- [Integration Guide](/Users/nbt-mingyi/mingyi.wu/easyid/sdk/docs/integration-guide.md)
- [Publishing Strategy](/Users/nbt-mingyi/mingyi.wu/easyid/sdk/docs/repository-publishing-strategy.md)
