package easyid

import (
	"context"
	"net/http"
)

// BankService handles bank card verification.
type BankService struct{ c *Client }

// BankVerify4Request contains the fields for a 4-element bank card verification.
type BankVerify4Request struct {
	Name     string `json:"name"`
	IDNumber string `json:"id_number"`
	BankCard string `json:"bank_card"`
	Mobile   string `json:"mobile,omitempty"` // 预留手机号（部分银行必填）
	TraceID  string `json:"trace_id,omitempty"`
}

// BankVerify4Result is the result of a 4-element bank card verification.
type BankVerify4Result struct {
	Result         bool    `json:"result"`
	Match          bool    `json:"match"`
	BankName       string  `json:"bank_name"`
	Supplier       string  `json:"supplier"`
	Score          float64 `json:"score"`
	MaskedBankCard string  `json:"masked_bank_card"`
	CardType       string  `json:"card_type"`
}

// Verify4 validates name + ID number + bank card + mobile (四要素银行卡核验).
func (s *BankService) Verify4(ctx context.Context, req *BankVerify4Request) (*BankVerify4Result, error) {
	var out BankVerify4Result
	if err := s.c.do(ctx, http.MethodPost, "/v1/bank/verify4", nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
