package easyid

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// sign 根据服务端规则计算签名
// 规则：to_sign = timestamp + "\n" + sorted_query + body
//
//	signature = hex(HMAC-SHA256(secret, to_sign))
func sign(secret, timestamp string, query url.Values, body []byte) string {
	keys := make([]string, 0, len(query))
	for k := range query {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for i, k := range keys {
		if i > 0 {
			sb.WriteByte('&')
		}
		sb.WriteString(k)
		sb.WriteByte('=')
		sb.WriteString(query.Get(k))
	}
	if len(body) > 0 {
		if sb.Len() > 0 {
			sb.WriteByte('&')
		}
		sb.Write(body)
	}

	toSign := fmt.Sprintf("%s\n%s", timestamp, sb.String())
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(toSign))
	return hex.EncodeToString(mac.Sum(nil))
}
