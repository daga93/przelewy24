package przelewy24

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAPIClient_RegisterTransaction(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		posID      int
		merchantID int
		crcKey     string
		reportKey  string
		setters    []APIClientOption
		// Named input parameters for target function.
		amount      int
		currency    Currency
		email       string
		description string
		country     Country
		language    Language
		sessionID   string
		opts        []RegisterTransactionOption
		want        string
		wantErr     bool
	}{
		{
			name:        "Happy Path",
			posID:       123,
			merchantID:  123,
			crcKey:      "test-crc-key",
			reportKey:   "test-report-key",
			setters:     []APIClientOption{},
			amount:      100,
			currency:    PLN,
			email:       "example@example.com",
			description: "Test Description",
			country:     CountryPL,
			language:    LanguagePL,
			sessionID:   "sess-1",
			opts:        []RegisterTransactionOption{},
			want:        "test-token",
			wantErr:     false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("Expected POST request, got %s", r.Method)
				}

				resp := map[string]interface{}{
					"data": map[string]string{
						"token": tt.want,
					},
				}
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(resp)
			}))

			c := NewClient(tt.posID, tt.merchantID, tt.crcKey, tt.reportKey, withURL(ts.URL))
			got, gotErr := c.RegisterTransaction(context.Background(), tt.amount, tt.currency, tt.email, tt.description, tt.country, tt.language, tt.sessionID, tt.opts...)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("RegisterTransaction() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("RegisterTransaction() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			wantURL := fmt.Sprintf("%s/trnRequest/%s", ts.URL, tt.want)
			if wantURL != got {
				t.Errorf("RegisterTransaction() = %v, want %v", got, wantURL)
			}

			ts.Close()
		})
	}
}

func TestNotificationWebhookHandler(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		resp := verifyTransactionResponseBody{
			Data:         verifyTransactionDataField{Status: "success"},
			ResponseCode: 0,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	client := NewClient(123, 456, "crc", "report", withURL(ts.URL))

	mockHandler := func(n *NotificationBody) (VerificationData, error) {
		return VerificationData{Amount: 1000, Currency: "PLN", SessionId: "Session-123", OrderId: 123}, nil
	}

	handler := client.NotificationWebhookHandler(mockHandler)

	// Tworzymy "sztuczne" żądanie od P24
	payload := `{"sessionId": "Session-123", "currency": "PLN", "amount": 1000, "orderId": 123}`
	req := httptest.NewRequest("POST", "/p24/webhook", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "5.252.202.255:1234"

	rr := httptest.NewRecorder()

	// Act
	handler(rr, req)

	// Assert
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", rr.Code)
	}
}
