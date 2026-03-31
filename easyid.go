// Package easyid provides a Go SDK for the EasyID API platform.
//
// # Quick start
//
//	client := easyid.New("ak_xxx", "sk_xxx")
//	result, err := client.IDCard.Verify2(ctx, &easyid.IDCardVerify2Request{
//	    Name:     "张三",
//	    IDNumber: "110101199001011234",
//	})
//
// # Authentication
//
// Every request is signed with HMAC-SHA256 using your secret key.
// The SDK handles signing automatically; you only need to supply keyID and secret.
//
// # Error handling
//
// Service errors are returned as *APIError with a numeric Code and Message.
// Use IsAPIError to inspect them:
//
//	if apiErr, ok := easyid.IsAPIError(err); ok {
//	    fmt.Println(apiErr.Code, apiErr.Message)
//	}
//
// # Thread safety
//
// A Client is safe for concurrent use. Create one and reuse it across goroutines.
package easyid

import (
	"net/http"
	"regexp"
	"time"
)

const (
	defaultBaseURL = "https://api.easyid.com"
	// Version is the current SDK version, sent in the User-Agent header.
	Version = "1.0.0"
)

var reKeyID = regexp.MustCompile(`^ak_[0-9a-f]+$`)

// Client is the EasyID API client. Create one with New() and reuse it across goroutines.
type Client struct {
	keyID   string
	secret  string
	baseURL string
	http    *http.Client

	// Service groups — mirrors the API surface.
	IDCard  *IDCardService
	Phone   *PhoneService
	Face    *FaceService
	Bank    *BankService
	Risk    *RiskService
	Billing *BillingService
}

// Option configures the client.
type Option func(*Client)

// WithBaseURL overrides the API base URL (useful for testing / private deployment).
func WithBaseURL(url string) Option {
	return func(c *Client) { c.baseURL = url }
}

// WithHTTPClient replaces the default HTTP client (useful for custom TLS / proxy).
func WithHTTPClient(h *http.Client) Option {
	return func(c *Client) { c.http = h }
}

// WithTimeout sets a request timeout (default 30s).
func WithTimeout(d time.Duration) Option {
	return func(c *Client) { c.http.Timeout = d }
}

// New creates a new EasyID client.
// keyID must have the form "ak_<hex>"; secret must be non-empty.
// Panics if either value is invalid — fail fast at startup.
//
//	client := easyid.New("ak_xxx", "sk_xxx")
func New(keyID, secret string, opts ...Option) *Client {
	if !reKeyID.MatchString(keyID) {
		panic("easyid: keyID must match ak_<hex>, got: " + keyID)
	}
	if secret == "" {
		panic("easyid: secret must not be empty")
	}
	c := &Client{
		keyID:   keyID,
		secret:  secret,
		baseURL: defaultBaseURL,
		http:    &http.Client{Timeout: 30 * time.Second},
	}
	for _, o := range opts {
		o(c)
	}
	c.IDCard = &IDCardService{c: c}
	c.Phone = &PhoneService{c: c}
	c.Face = &FaceService{c: c}
	c.Bank = &BankService{c: c}
	c.Risk = &RiskService{c: c}
	c.Billing = &BillingService{c: c}
	return c
}
