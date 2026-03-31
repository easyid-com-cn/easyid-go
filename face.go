package easyid

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

// FaceService handles face liveness detection, comparison, and verification.
type FaceService struct{ c *Client }

// LivenessMode controls liveness detection mode.
type LivenessMode string

const (
	LivenessModeActive  LivenessMode = "active"
	LivenessModePassive LivenessMode = "passive"
)

// LivenessResult is the result of a liveness detection request.
type LivenessResult struct {
	Liveness       bool    `json:"liveness"`
	Score          float64 `json:"score"`
	Method         string  `json:"method"`
	FramesAnalyzed int     `json:"frames_analyzed"`
	AttackType     *string `json:"attack_type"`
}

// Liveness uploads a video/image and detects whether it is a live face.
// mode is optional; pass empty string to use server default.
func (s *FaceService) Liveness(ctx context.Context, mode LivenessMode, media io.Reader, filename string) (*LivenessResult, error) {
	var out LivenessResult
	err := s.c.doMultipart(ctx, "/v1/face/liveness", func(w *multipart.Writer) error {
		if mode != "" {
			if err := w.WriteField("mode", string(mode)); err != nil {
				return err
			}
		}
		part, err := w.CreateFormFile("media", filename)
		if err != nil {
			return err
		}
		if _, err := io.Copy(part, media); err != nil {
			return fmt.Errorf("copy media: %w", err)
		}
		return nil
	}, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// CompareResult is the result of a face comparison request.
type CompareResult struct {
	Match bool    `json:"match"`
	Score float64 `json:"score"`
}

// Compare uploads two face images and returns a similarity score.
func (s *FaceService) Compare(ctx context.Context, image1, image2 io.Reader, filename1, filename2 string) (*CompareResult, error) {
	var out CompareResult
	err := s.c.doMultipart(ctx, "/v1/face/compare", func(w *multipart.Writer) error {
		p1, err := w.CreateFormFile("image1", filename1)
		if err != nil {
			return err
		}
		if _, err := io.Copy(p1, image1); err != nil {
			return fmt.Errorf("copy image1: %w", err)
		}
		p2, err := w.CreateFormFile("image2", filename2)
		if err != nil {
			return err
		}
		if _, err := io.Copy(p2, image2); err != nil {
			return fmt.Errorf("copy image2: %w", err)
		}
		return nil
	}, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// FaceVerifyRequest contains the fields for face + ID verification.
type FaceVerifyRequest struct {
	IDNumber    string `json:"id_number"`
	MediaKey    string `json:"media_key,omitempty"`
	CallbackURL string `json:"callback_url,omitempty"`
}

// FaceVerifyResult is the result of a face + ID verification.
type FaceVerifyResult struct {
	Result   bool    `json:"result"`
	Supplier string  `json:"supplier"`
	Score    float64 `json:"score"`
}

// Verify validates a face against an ID number (人脸核验).
func (s *FaceService) Verify(ctx context.Context, req *FaceVerifyRequest) (*FaceVerifyResult, error) {
	var out FaceVerifyResult
	if err := s.c.do(ctx, http.MethodPost, "/v1/face/verify", nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
