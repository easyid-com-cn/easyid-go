package easyid

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const maxResponseBytes = 10 << 20 // 10 MB

// apiResponse is the envelope for every EasyID response.
type apiResponse struct {
	Code      int             `json:"code"`
	Message   string          `json:"message"`
	RequestID string          `json:"request_id"`
	Data      json.RawMessage `json:"data"`
}

// do executes a JSON request and unmarshals data into out (may be nil).
func (c *Client) do(ctx context.Context, method, path string, query url.Values, body interface{}, out interface{}) error {
	var bodyBytes []byte
	var err error
	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("easyid: marshal request: %w", err)
		}
	}
	return c.doRaw(ctx, method, path, query, bodyBytes, "application/json", out)
}

// doMultipart executes a multipart/form-data request.
func (c *Client) doMultipart(ctx context.Context, path string, mw func(*multipart.Writer) error, out interface{}) error {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	if err := mw(w); err != nil {
		return fmt.Errorf("easyid: build multipart: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("easyid: close multipart writer: %w", err)
	}
	return c.doRaw(ctx, http.MethodPost, path, nil, buf.Bytes(), w.FormDataContentType(), out)
}

func (c *Client) doRaw(ctx context.Context, method, path string, query url.Values, bodyBytes []byte, contentType string, out interface{}) error {
	ts := strconv.FormatInt(time.Now().Unix(), 10)

	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return fmt.Errorf("easyid: parse url: %w", err)
	}
	if len(query) > 0 {
		u.RawQuery = query.Encode()
	}

	sig := sign(c.secret, ts, query, bodyBytes)

	req, err := http.NewRequestWithContext(ctx, method, u.String(), bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("easyid: build request: %w", err)
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("X-Key-ID", c.keyID)
	req.Header.Set("X-Timestamp", ts)
	req.Header.Set("X-Signature", sig)
	req.Header.Set("User-Agent", "easyid-go/"+Version)

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("easyid: http: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes))
	if err != nil {
		return fmt.Errorf("easyid: read response: %w", err)
	}

	// Non-2xx with non-JSON body (e.g. nginx 502 HTML) → surface HTTP status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var env apiResponse
		if jsonErr := json.Unmarshal(raw, &env); jsonErr == nil && env.Code != 0 {
			return &APIError{Code: env.Code, Message: env.Message, RequestID: env.RequestID}
		}
		return fmt.Errorf("easyid: http status %d", resp.StatusCode)
	}

	var env apiResponse
	if err := json.Unmarshal(raw, &env); err != nil {
		return fmt.Errorf("easyid: decode response (status=%d): %w", resp.StatusCode, err)
	}
	if env.Code != 0 {
		return &APIError{Code: env.Code, Message: env.Message, RequestID: env.RequestID}
	}
	if out != nil && len(env.Data) > 0 {
		if err := json.Unmarshal(env.Data, out); err != nil {
			return fmt.Errorf("easyid: decode data: %w", err)
		}
	}
	return nil
}
