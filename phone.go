package easyid

import (
	"context"
	"net/http"
	"net/url"
)

// PhoneService handles phone number verification.
type PhoneService struct{ c *Client }

// PhoneStatusResult contains carrier and status info for a phone number.
type PhoneStatusResult struct {
	Status   string `json:"status"`  // real / virtual / empty / unknown
	Carrier  string `json:"carrier"` // 运营商
	Province string `json:"province"`
	Roaming  bool   `json:"roaming"`
}

// Status queries the carrier and real/virtual status of a phone number.
func (s *PhoneService) Status(ctx context.Context, phone string) (*PhoneStatusResult, error) {
	q := url.Values{"phone": {phone}}
	var out PhoneStatusResult
	if err := s.c.do(ctx, http.MethodGet, "/v1/phone/status", q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// PhoneVerify3Request contains the fields for a 3-element phone verification.
type PhoneVerify3Request struct {
	Name     string `json:"name"`
	IDNumber string `json:"id_number"`
	Mobile   string `json:"mobile"`
}

// PhoneVerify3Result is the result of a 3-element phone verification.
type PhoneVerify3Result struct {
	Result   bool    `json:"result"`
	Match    bool    `json:"match"`
	Supplier string  `json:"supplier"`
	Score    float64 `json:"score"`
}

// Verify3 validates name + ID number + mobile (三要素手机核验).
func (s *PhoneService) Verify3(ctx context.Context, req *PhoneVerify3Request) (*PhoneVerify3Result, error) {
	var out PhoneVerify3Result
	if err := s.c.do(ctx, http.MethodPost, "/v1/phone/verify3", nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
