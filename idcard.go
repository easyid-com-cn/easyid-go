package easyid

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

// IDCardService handles identity card verification.
type IDCardService struct{ c *Client }

// --- Verify2 (二要素：姓名 + 身份证号) ---

type IDCardVerify2Request struct {
	Name     string `json:"name"`
	IDNumber string `json:"id_number"`
	TraceID  string `json:"trace_id,omitempty"`
}

// --- Verify3 (三要素：姓名 + 身份证号 + 手机号) ---

type IDCardVerify3Request struct {
	Name     string `json:"name"`
	IDNumber string `json:"id_number"`
	Mobile   string `json:"mobile"`
	TraceID  string `json:"trace_id,omitempty"`
}

// IDCardVerifyResult is shared by Verify2 and Verify3.
type IDCardVerifyResult struct {
	Result   bool            `json:"result"`
	Match    bool            `json:"match"`
	Supplier string          `json:"supplier"`
	Score    float64         `json:"score"`
	Raw      json.RawMessage `json:"raw,omitempty"`
}

// Verify2 validates name + ID number (二要素核验).
func (s *IDCardService) Verify2(ctx context.Context, req *IDCardVerify2Request) (*IDCardVerifyResult, error) {
	var out IDCardVerifyResult
	if err := s.c.do(ctx, http.MethodPost, "/v1/idcard/verify2", nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Verify3 validates name + ID number + mobile (三要素核验).
func (s *IDCardService) Verify3(ctx context.Context, req *IDCardVerify3Request) (*IDCardVerifyResult, error) {
	var out IDCardVerifyResult
	if err := s.c.do(ctx, http.MethodPost, "/v1/idcard/verify3", nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// --- OCR ---

// IDCardSide is the side of an ID card image.
type IDCardSide string

const (
	IDCardSideFront IDCardSide = "front"
	IDCardSideBack  IDCardSide = "back"
)

// OCRIDCardResult contains fields extracted from an ID card image.
type OCRIDCardResult struct {
	Side     string          `json:"side"`
	Name     string          `json:"name,omitempty"`
	IDNumber string          `json:"id_number,omitempty"`
	Gender   string          `json:"gender,omitempty"`
	Nation   string          `json:"nation,omitempty"`
	Birth    string          `json:"birth,omitempty"`
	Address  string          `json:"address,omitempty"`
	Issue    string          `json:"issue,omitempty"`
	Valid    string          `json:"valid,omitempty"`
	Raw      json.RawMessage `json:"raw,omitempty"`
}

// OCR uploads an ID card image and returns the extracted fields.
// image should be an open *os.File or any io.Reader; filename is used as the multipart filename.
func (s *IDCardService) OCR(ctx context.Context, side IDCardSide, image io.Reader, filename string) (*OCRIDCardResult, error) {
	var out OCRIDCardResult
	err := s.c.doMultipart(ctx, "/v1/ocr/idcard", func(w *multipart.Writer) error {
		if err := w.WriteField("side", string(side)); err != nil {
			return err
		}
		part, err := w.CreateFormFile("image", filename)
		if err != nil {
			return err
		}
		if _, err := io.Copy(part, image); err != nil {
			return fmt.Errorf("copy image: %w", err)
		}
		return nil
	}, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
