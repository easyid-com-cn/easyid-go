package easyid_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	easyid "github.com/easyid-com-cn/easyid-go"
)

// --- helpers ---

// mockServer starts a test HTTP server. handler receives the request and returns
// a raw JSON body (just the "data" field value).
func mockServer(t *testing.T, wantPath, wantMethod string, dataJSON string) (*httptest.Server, func()) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != wantPath {
			t.Errorf("path: got %q, want %q", r.URL.Path, wantPath)
		}
		if r.Method != wantMethod {
			t.Errorf("method: got %q, want %q", r.Method, wantMethod)
		}
		// verify required auth headers
		if r.Header.Get("X-Key-ID") == "" {
			t.Error("missing X-Key-ID header")
		}
		if r.Header.Get("X-Timestamp") == "" {
			t.Error("missing X-Timestamp header")
		}
		if r.Header.Get("X-Signature") == "" {
			t.Error("missing X-Signature header")
		}
		w.Header().Set("Content-Type", "application/json")
		resp := `{"code":0,"message":"success","request_id":"test-rid","data":` + dataJSON + `}`
		w.Write([]byte(resp))
	}))
	return srv, srv.Close
}

func newTestClient(t *testing.T, srv *httptest.Server) *easyid.Client {
	t.Helper()
	return easyid.New("ak_3f9a2b1c7d4e8f0a", "sk_test",
		easyid.WithBaseURL(srv.URL),
		easyid.WithTimeout(5*time.Second),
	)
}

// --- signer ---

func TestSign_EmptyBody(t *testing.T) {
	// Verifies sign does not panic on empty inputs — indirect via a real request
	srv, close := mockServer(t, "/v1/phone/status", http.MethodGet,
		`{"status":"real","carrier":"移动","province":"广东","roaming":false}`)
	defer close()

	client := newTestClient(t, srv)
	res, err := client.Phone.Status(context.Background(), "13800138000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Status != "real" {
		t.Errorf("status: got %q, want %q", res.Status, "real")
	}
}

// --- error handling ---

func TestAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"code":1001,"message":"invalid key_id","request_id":"err-rid","data":null}`))
	}))
	defer srv.Close()

	client := newTestClient(t, srv)
	_, err := client.Phone.Status(context.Background(), "13800138000")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr, ok := easyid.IsAPIError(err)
	if !ok {
		t.Fatalf("expected *APIError, got %T: %v", err, err)
	}
	if apiErr.Code != 1001 {
		t.Errorf("code: got %d, want 1001", apiErr.Code)
	}
	if apiErr.RequestID != "err-rid" {
		t.Errorf("request_id: got %q, want err-rid", apiErr.RequestID)
	}
}

// --- IDCard ---

func TestIDCard_Verify2(t *testing.T) {
	srv, close := mockServer(t, "/v1/idcard/verify2", http.MethodPost,
		`{"result":true,"match":true,"supplier":"aliyun","score":0.98}`)
	defer close()

	client := newTestClient(t, srv)
	res, err := client.IDCard.Verify2(context.Background(), &easyid.IDCardVerify2Request{
		Name:     "张三",
		IDNumber: "110101199001011234",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Result {
		t.Error("expected result=true")
	}
	if res.Supplier != "aliyun" {
		t.Errorf("supplier: got %q, want aliyun", res.Supplier)
	}
}

func TestIDCard_Verify3(t *testing.T) {
	srv, close := mockServer(t, "/v1/idcard/verify3", http.MethodPost,
		`{"result":true,"match":true,"supplier":"tencent","score":0.95}`)
	defer close()

	client := newTestClient(t, srv)
	res, err := client.IDCard.Verify3(context.Background(), &easyid.IDCardVerify3Request{
		Name:     "张三",
		IDNumber: "110101199001011234",
		Mobile:   "13800138000",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Match {
		t.Error("expected match=true")
	}
}

func TestIDCard_OCR(t *testing.T) {
	srv, close := mockServer(t, "/v1/ocr/idcard", http.MethodPost,
		`{"side":"front","name":"张三","id_number":"110101199001011234"}`)
	defer close()

	client := newTestClient(t, srv)
	res, err := client.IDCard.OCR(context.Background(),
		easyid.IDCardSideFront,
		strings.NewReader("fake-image-bytes"),
		"id_front.jpg",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Name != "张三" {
		t.Errorf("name: got %q, want 张三", res.Name)
	}
}

// --- Phone ---

func TestPhone_Status(t *testing.T) {
	srv, close := mockServer(t, "/v1/phone/status", http.MethodGet,
		`{"status":"virtual","carrier":"联通","province":"北京","roaming":false}`)
	defer close()

	client := newTestClient(t, srv)
	res, err := client.Phone.Status(context.Background(), "17000000000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Status != "virtual" {
		t.Errorf("status: got %q, want virtual", res.Status)
	}
}

func TestPhone_Verify3(t *testing.T) {
	srv, close := mockServer(t, "/v1/phone/verify3", http.MethodPost,
		`{"result":true,"match":true,"supplier":"aliyun","score":0.99}`)
	defer close()

	client := newTestClient(t, srv)
	res, err := client.Phone.Verify3(context.Background(), &easyid.PhoneVerify3Request{
		Name:     "张三",
		IDNumber: "110101199001011234",
		Mobile:   "13800138000",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Result {
		t.Error("expected result=true")
	}
}

// --- Face ---

func TestFace_Liveness(t *testing.T) {
	srv, close := mockServer(t, "/v1/face/liveness", http.MethodPost,
		`{"liveness":true,"score":0.97,"method":"passive","frames_analyzed":10,"attack_type":null}`)
	defer close()

	client := newTestClient(t, srv)
	res, err := client.Face.Liveness(context.Background(),
		easyid.LivenessModePassive,
		strings.NewReader("fake-video"),
		"video.mp4",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Liveness {
		t.Error("expected liveness=true")
	}
}

func TestFace_Compare(t *testing.T) {
	srv, close := mockServer(t, "/v1/face/compare", http.MethodPost,
		`{"match":true,"score":0.92}`)
	defer close()

	client := newTestClient(t, srv)
	res, err := client.Face.Compare(context.Background(),
		strings.NewReader("img1"), strings.NewReader("img2"),
		"face1.jpg", "face2.jpg",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Match {
		t.Error("expected match=true")
	}
}

func TestFace_Verify(t *testing.T) {
	srv, close := mockServer(t, "/v1/face/verify", http.MethodPost,
		`{"result":true,"supplier":"aliyun","score":0.96}`)
	defer close()

	client := newTestClient(t, srv)
	res, err := client.Face.Verify(context.Background(), &easyid.FaceVerifyRequest{
		IDNumber: "110101199001011234",
		MediaKey: "oss://bucket/key",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Result {
		t.Error("expected result=true")
	}
}

// --- Bank ---

func TestBank_Verify4(t *testing.T) {
	srv, close := mockServer(t, "/v1/bank/verify4", http.MethodPost,
		`{"result":true,"match":true,"bank_name":"工商银行","supplier":"aliyun","score":0.99,"masked_bank_card":"6222****1234","card_type":"debit"}`)
	defer close()

	client := newTestClient(t, srv)
	res, err := client.Bank.Verify4(context.Background(), &easyid.BankVerify4Request{
		Name:     "张三",
		IDNumber: "110101199001011234",
		BankCard: "6222021234567890",
		Mobile:   "13800138000",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.BankName != "工商银行" {
		t.Errorf("bank_name: got %q, want 工商银行", res.BankName)
	}
	if res.CardType != "debit" {
		t.Errorf("card_type: got %q, want debit", res.CardType)
	}
}

// --- Risk ---

func TestRisk_Score(t *testing.T) {
	srv, close := mockServer(t, "/v1/risk/score", http.MethodPost,
		`{"risk_score":30,"reasons":["new_device"],"recommendation":"allow","details":{"rule_score":null,"ml_score":null}}`)
	defer close()

	client := newTestClient(t, srv)
	res, err := client.Risk.Score(context.Background(), &easyid.RiskScoreRequest{
		IP:       "1.2.3.4",
		DeviceID: "dev_abc",
		Action:   "login",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.RiskScore != 30 {
		t.Errorf("risk_score: got %d, want 30", res.RiskScore)
	}
	if res.Recommendation != "allow" {
		t.Errorf("recommendation: got %q, want allow", res.Recommendation)
	}
}

func TestRisk_StoreFingerprint(t *testing.T) {
	srv, close := mockServer(t, "/v1/device/fingerprint", http.MethodPost,
		`{"device_id":"dev_abc","stored":true}`)
	defer close()

	fp, _ := json.Marshal(map[string]string{"canvas": "hash123", "webgl": "hash456"})
	client := newTestClient(t, srv)
	res, err := client.Risk.StoreFingerprint(context.Background(), &easyid.StoreFingerprintRequest{
		DeviceID:    "dev_abc",
		Fingerprint: fp,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Stored {
		t.Error("expected stored=true")
	}
}

// --- Billing ---

func TestBilling_Balance(t *testing.T) {
	srv, close := mockServer(t, "/v1/billing/balance", http.MethodGet,
		`{"app_id":"app_001","available_cents":100000}`)
	defer close()

	client := newTestClient(t, srv)
	res, err := client.Billing.Balance(context.Background(), "app_001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.AvailableCents != 100000 {
		t.Errorf("available_cents: got %d, want 100000", res.AvailableCents)
	}
}

func TestBilling_Records(t *testing.T) {
	srv, close := mockServer(t, "/v1/billing/records", http.MethodGet,
		`{"total":1,"page":1,"records":[{"id":1,"app_id":"app_001","change_cents":-100,"balance_before":100100,"balance_after":100000,"reason":"idcard_verify2","operator":"system","created_at":1711900000}]}`)
	defer close()

	client := newTestClient(t, srv)
	res, err := client.Billing.Records(context.Background(), "app_001", 1, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Total != 1 {
		t.Errorf("total: got %d, want 1", res.Total)
	}
	if len(res.Records) != 1 {
		t.Fatalf("records len: got %d, want 1", len(res.Records))
	}
	if res.Records[0].ChangeCents != -100 {
		t.Errorf("change_cents: got %d, want -100", res.Records[0].ChangeCents)
	}
}

// --- client options ---

func TestWithHTTPClient(t *testing.T) {
	custom := &http.Client{Timeout: 1 * time.Second}
	client := easyid.New("ak_3f9a2b1c7d4e8f0a", "sk_test", easyid.WithHTTPClient(custom))
	if client == nil {
		t.Fatal("client is nil")
	}
}

// --- keyID validation ---

func TestNew_InvalidKeyID(t *testing.T) {
	cases := []string{
		"",
		"sk_abc",              // wrong prefix
		"ak_test",             // non-hex suffix
		"ak_\r\nEvil: 1",     // header injection attempt
		"ak_UPPERCASE",        // uppercase not allowed
	}
	for _, id := range cases {
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("New(%q) expected panic, got none", id)
				}
			}()
			easyid.New(id, "sk_secret")
		}()
	}
}

func TestNew_EmptySecret(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("New with empty secret should panic")
		}
	}()
	easyid.New("ak_3f9a2b1c7d4e8f0a", "")
}

// --- HTTP error handling ---

func TestHTTP5xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("<html>503 Service Unavailable</html>"))
	}))
	defer srv.Close()

	client := newTestClient(t, srv)
	_, err := client.Phone.Status(context.Background(), "13800138000")
	if err == nil {
		t.Fatal("expected error on 503, got nil")
	}
	if _, ok := easyid.IsAPIError(err); ok {
		t.Error("503 HTML body should not produce APIError")
	}
}

func TestHTTP5xx_WithJSONBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"code":5000,"message":"internal server error","request_id":"err-500","data":null}`))
	}))
	defer srv.Close()

	client := newTestClient(t, srv)
	_, err := client.Phone.Status(context.Background(), "13800138000")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr, ok := easyid.IsAPIError(err)
	if !ok {
		t.Fatalf("expected *APIError, got %T: %v", err, err)
	}
	if apiErr.Code != 5000 {
		t.Errorf("code: got %d, want 5000", apiErr.Code)
	}
}

// --- User-Agent header ---

func TestUserAgentHeader(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ua := r.Header.Get("User-Agent")
		if !strings.HasPrefix(ua, "easyid-go/") {
			t.Errorf("User-Agent: got %q, want prefix easyid-go/", ua)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"code":0,"message":"success","request_id":"r","data":{"status":"real","carrier":"","province":"","roaming":false}}`))
	}))
	defer srv.Close()

	client := newTestClient(t, srv)
	client.Phone.Status(context.Background(), "13800138000") //nolint:errcheck
}
