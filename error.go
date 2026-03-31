package easyid

import "fmt"

// APIError 服务端返回的业务错误
type APIError struct {
	Code      int    // 业务错误码，非 0
	Message   string // 错误描述
	RequestID string // 便于排查的请求 ID
}

func (e *APIError) Error() string {
	return fmt.Sprintf("easyid: code=%d message=%s request_id=%s", e.Code, e.Message, e.RequestID)
}

// IsAPIError 判断 err 是否为服务端业务错误
func IsAPIError(err error) (*APIError, bool) {
	if e, ok := err.(*APIError); ok {
		return e, true
	}
	return nil, false
}
