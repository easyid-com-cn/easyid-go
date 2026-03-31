package easyid

import (
	"context"
	"encoding/json"
	"net/http"
)

// RiskService handles risk scoring and device fingerprint storage.
type RiskService struct{ c *Client }

// RiskScoreRequest contains the context signals for a risk score request.
// All fields are optional; provide as many as available for better accuracy.
type RiskScoreRequest struct {
	IP                string          `json:"ip,omitempty"`
	DeviceFingerprint string          `json:"device_fingerprint,omitempty"`
	DeviceID          string          `json:"device_id,omitempty"`
	Phone             string          `json:"phone,omitempty"`
	Email             string          `json:"email,omitempty"`
	UserAgent         string          `json:"user_agent,omitempty"`
	Action            string          `json:"action,omitempty"`
	Amount            int64           `json:"amount,omitempty"`
	Context           json.RawMessage `json:"context,omitempty"`
}

// RiskScoreResult is the result of a risk score request.
type RiskScoreResult struct {
	RiskScore      int      `json:"risk_score"`      // 0-100，越高风险越大
	Reasons        []string `json:"reasons"`
	Recommendation string   `json:"recommendation"` // allow / review / block
	Details        struct {
		RuleScore *int `json:"rule_score"`
		MLScore   *int `json:"ml_score"`
	} `json:"details"`
}

// Score returns a risk score for the given context signals.
func (s *RiskService) Score(ctx context.Context, req *RiskScoreRequest) (*RiskScoreResult, error) {
	var out RiskScoreResult
	if err := s.c.do(ctx, http.MethodPost, "/v1/risk/score", nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// StoreFingerprintRequest stores a device fingerprint for future risk analysis.
type StoreFingerprintRequest struct {
	DeviceID    string          `json:"device_id"`
	Fingerprint json.RawMessage `json:"fingerprint"` // SDK/browser 采集的原始指纹 JSON
}

// StoreFingerprintResult is the result of storing a device fingerprint.
type StoreFingerprintResult struct {
	DeviceID string `json:"device_id"`
	Stored   bool   `json:"stored"`
}

// StoreFingerprint stores a device fingerprint for risk analysis.
func (s *RiskService) StoreFingerprint(ctx context.Context, req *StoreFingerprintRequest) (*StoreFingerprintResult, error) {
	var out StoreFingerprintResult
	if err := s.c.do(ctx, http.MethodPost, "/v1/device/fingerprint", nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
