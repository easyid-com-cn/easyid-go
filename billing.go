package easyid

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// BillingService handles balance and billing record queries.
type BillingService struct{ c *Client }

// BalanceResult contains the current available balance for an app.
type BalanceResult struct {
	AppID          string `json:"app_id"`
	AvailableCents int64  `json:"available_cents"` // 单位：分
}

// Balance returns the available balance for the authenticated app.
func (s *BillingService) Balance(ctx context.Context, appID string) (*BalanceResult, error) {
	q := url.Values{"app_id": {appID}}
	var out BalanceResult
	if err := s.c.do(ctx, http.MethodGet, "/v1/billing/balance", q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// BillingRecord is a single billing record.
type BillingRecord struct {
	ID            int64  `json:"id"`
	AppID         string `json:"app_id"`
	RequestID     string `json:"request_id"`
	ChangeCents   int64  `json:"change_cents"`   // 正数=充值，负数=扣费
	BalanceBefore int64  `json:"balance_before"`
	BalanceAfter  int64  `json:"balance_after"`
	Reason        string `json:"reason"`
	Operator      string `json:"operator"`
	CreatedAt     int64  `json:"created_at"` // Unix 秒
}

// BillingRecordsResult contains a paginated list of billing records.
type BillingRecordsResult struct {
	Total   int64           `json:"total"`
	Page    int             `json:"page"`
	Records []BillingRecord `json:"records"`
}

// Records returns paginated billing records. page and pageSize default to 1/20 if zero.
func (s *BillingService) Records(ctx context.Context, appID string, page, pageSize int) (*BillingRecordsResult, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	q := url.Values{
		"app_id":    {appID},
		"page":      {strconv.Itoa(page)},
		"page_size": {strconv.Itoa(pageSize)},
	}
	var out BillingRecordsResult
	if err := s.c.do(ctx, http.MethodGet, "/v1/billing/records", q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
