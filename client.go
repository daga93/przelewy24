package przelewy24

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"
)

type APIClientOptions struct {
	TestEnvironment             bool
	ShopDomain                  string
	TransactionNotificationPath string
	RedirectPath                string
	ClientTimeout               time.Duration
	Logger                      *slog.Logger
}

type APIClientOption func(*APIClientOptions)

func WithTestEnvironment(isTestEnvironment bool) APIClientOption {
	return func(ao *APIClientOptions) {
		ao.TestEnvironment = isTestEnvironment
	}
}

func WithShopDomain(shopDomain string) APIClientOption {
	return func(ao *APIClientOptions) {
		ao.ShopDomain = shopDomain
	}
}

func WithTransactionNotificationPath(p string) APIClientOption {
	return func(ao *APIClientOptions) {
		ao.TransactionNotificationPath = p
	}
}

func WithRedirectPath(p string) APIClientOption {
	return func(ao *APIClientOptions) {
		ao.RedirectPath = p
	}
}

func WithClientTimeout(t time.Duration) APIClientOption {
	return func(ao *APIClientOptions) {
		ao.ClientTimeout = t
	}
}

func WithLogger(logger *slog.Logger) APIClientOption {
	return func(ao *APIClientOptions) {
		ao.Logger = logger
	}
}

type APIClient struct {
	client *http.Client
	logger *slog.Logger

	merchantID int
	posID      int
	crcKey     string
	reportKey  string
	domain     string

	url                             string
	transactionRegisterEndpoint     string
	transactionRedirectionEndpoint  string
	transactionVerificationEndpoint string

	urlStatus string
	urlReturn string
}

type NotificationHandler func(notification *NotificationBody) (VerificationData, error)

// NewClient creates new Przelewy24 API Client.
// This is main client that is used to perform
// actions on Przelewy24 API.
func NewClient(posID, merchantID int, crcKey, reportKey string, setters ...APIClientOption) *APIClient {
	args := &APIClientOptions{
		TestEnvironment:             true,
		ShopDomain:                  "http://localhost:8080",
		TransactionNotificationPath: "/p24/payment/transaction/notification",
		ClientTimeout:               30 * time.Second,
		RedirectPath:                "/",
		Logger:                      slog.New(slog.NewTextHandler(io.Discard, nil)),
	}

	for _, setter := range setters {
		setter(args)
	}

	client := &APIClient{
		client: &http.Client{
			Timeout: args.ClientTimeout,
		},
		logger:                          args.Logger,
		merchantID:                      merchantID,
		posID:                           posID,
		crcKey:                          crcKey,
		reportKey:                       reportKey,
		url:                             "https://sandbox.przelewy24.pl",
		transactionRegisterEndpoint:     "/api/v1/transaction/register",
		transactionRedirectionEndpoint:  "/trnRequest",
		transactionVerificationEndpoint: "/api/v1/transaction/verify",
		domain:                          args.ShopDomain,
		urlStatus:                       fmt.Sprintf("%s%s", args.ShopDomain, args.TransactionNotificationPath),
		urlReturn:                       fmt.Sprintf("%s%s", args.ShopDomain, args.RedirectPath),
	}

	if args.TestEnvironment == false {
		client.url = "https://secure.przelewy24.pl"
	}
	return client
}

type RegisterTransactionOption func(*registerTransactionRequestBody)

func WithClient(c string) RegisterTransactionOption {
	return func(r *registerTransactionRequestBody) {
		r.Client = c
	}
}

func WithAddress(a string) RegisterTransactionOption {
	return func(r *registerTransactionRequestBody) {
		r.Address = a
	}
}
func WithZip(z string) RegisterTransactionOption {
	return func(r *registerTransactionRequestBody) {
		r.Zip = z
	}
}

func WithCity(c string) RegisterTransactionOption {
	return func(r *registerTransactionRequestBody) {
		r.City = c
	}
}

func WithPhone(p string) RegisterTransactionOption {
	return func(r *registerTransactionRequestBody) {
		r.Phone = p
	}
}
func WithMethod(m int) RegisterTransactionOption {
	return func(r *registerTransactionRequestBody) {
		r.Method = m
	}
}

func WithUrlReturn(u string) RegisterTransactionOption {
	return func(r *registerTransactionRequestBody) {
		r.UrlReturn = u
	}
}

func WithUrlStatus(u string) RegisterTransactionOption {
	return func(r *registerTransactionRequestBody) {
		r.UrlStatus = u
	}
}

func WithUrlNotify(u string) RegisterTransactionOption {
	return func(r *registerTransactionRequestBody) {
		r.UrlNotify = u
	}
}

// In minutes.
func WithTimeLimit(t int) RegisterTransactionOption {
	return func(r *registerTransactionRequestBody) {
		r.TimeLimit = t
	}
}

func WithChannel(c Channel) RegisterTransactionOption {
	return func(r *registerTransactionRequestBody) {
		r.Channel = c
	}
}
func WithWaitForResult(w bool) RegisterTransactionOption {
	return func(r *registerTransactionRequestBody) {
		r.WaitForResult = w
	}
}
func WithRegulationAccept(w bool) RegisterTransactionOption {
	return func(r *registerTransactionRequestBody) {
		r.RegulationAccept = w
	}
}

func WithShipping(s int) RegisterTransactionOption {
	return func(r *registerTransactionRequestBody) {
		r.Shipping = s
	}
}

func WithTransferLabel(t string) RegisterTransactionOption {
	return func(r *registerTransactionRequestBody) {
		r.TransferLabel = t
	}
}

func WithEncoding(e string) RegisterTransactionOption {
	return func(r *registerTransactionRequestBody) {
		r.Encoding = e
	}
}

func WithMethodRefId(m string) RegisterTransactionOption {
	return func(r *registerTransactionRequestBody) {
		r.MethodRefId = m
	}
}
func WithCart(c Cart) RegisterTransactionOption {
	return func(r *registerTransactionRequestBody) {
		r.Cart = c
	}
}

func WitAdditional(a *AdditionalData) RegisterTransactionOption {
	return func(r *registerTransactionRequestBody) {
		r.Additional = a
	}
}

// RegisterTransaction registers transaction in Przelewy24.
// Minor currency is used - amount should be in grosz (PLN),
// cents (EUR/USD) etc. Be sure to pass unique sessionID -
// Przelewy24 will respond with error when duplicated sessionID is passed.
func (c *APIClient) RegisterTransaction(
	ctx context.Context,
	amount int,
	currency Currency,
	email string,
	description string,
	country Country,
	language Language,
	sessionID string,
	opts ...RegisterTransactionOption,
) (string, error) {
	body := registerTransactionRequestBody{
		SessionId:   sessionID,
		MerchantId:  c.merchantID,
		PosId:       c.posID,
		Email:       email,
		Amount:      amount,
		Currency:    currency,
		Description: description,
		Country:     country,
		Language:    language,
		UrlReturn:   c.urlReturn,
		UrlStatus:   c.urlStatus,
	}

	for _, opt := range opts {
		opt(&body)
	}

	signature := signRegistration(body.SessionId, c.merchantID, amount, currency, c.crcKey)
	body.Sign = hex.EncodeToString(signature[:])

	b, err := json.Marshal(body)
	if err != nil {
		c.logger.Error("Marshalling Register Transaction body failed.", "error", err)
		return "", fmt.Errorf("p24: failed to marshal on registering transaction: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, "POST", c.url+c.transactionRegisterEndpoint, bytes.NewReader(b))
	if err != nil {
		c.logger.Error("Creating new request failed.", "error", err)
		return "", fmt.Errorf("p24: failed to create new request: %w", err)
	}
	c.setHeaders(req)

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("Failed to call Przelewy24 in register transaction request.", "error", err)
		return "", fmt.Errorf("p24: failed to make a request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("Failed to register transaction in Przelewy24.", "error", err, "status", resp.StatusCode)
		return "", fmt.Errorf("p24: error when creating transaction: %w", err)
	}
	r := &registerTransactionResponseBody{}
	err = json.Unmarshal(responseBody, r)
	if err != nil {
		c.logger.Error("Unmarshalling Register Transaction body failed.", "error", err)
		return "", fmt.Errorf("p24: failed to unmarshal transaction response: %w", err)
	}

	if resp.StatusCode != 200 {
		c.logger.Error("Unexpected error with transaction registration - aborting.", "error", r.ErrorData)
		return "", fmt.Errorf("p24: unexpected status code (%s) when registering transaction; message: '%s'", resp.Status, r.ErrorData)
	}
	c.logger.Info("Transaction successfully registered.", "session_id", body.SessionId)

	return fmt.Sprintf("%s%s/%s", c.url, c.transactionRedirectionEndpoint, r.Data.Token), nil
}

var trustedIPs = []string{
	"5.252.202.255",
	"5.252.202.254",
	"20.215.81.124",
}

// NotificationWebhookHandler returns a handler that handles Webhook Notification about Transaction Registration.
// The handler that is passed to this function should return data ([VerificationData]) from your database.
// This data is used to verify, that transaction was correct.
func (c *APIClient) NotificationWebhookHandler(handler NotificationHandler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body != nil {
			defer r.Body.Close()
		}

		ipStr := getIP(r)
		if !isIPInList(ipStr, trustedIPs) {
			c.logger.Warn("Notification from untrusted IP. Aborting.", "ip", ipStr)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		var notification NotificationBody
		if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
			c.logger.Error("Failed to decode p24 notification.", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		data, err := handler(&notification)
		if err != nil {
			c.logger.Error("Transaction could not be verified - error on execution of Notification Handler", "error", err)
			w.WriteHeader(http.StatusPreconditionFailed)
			return
		}

		err = c.verifyTransaction(r.Context(), data)
		if err != nil {
			c.logger.Error("Transaction could not Verified.", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		c.logger.Info("Transaction Verification Success.", "session_id", notification.SessionId)
	}
}

// verifyTransaction sends request to Przelewy24, with acknowledgement that the transaction is correct.
func (c *APIClient) verifyTransaction(ctx context.Context, rb VerificationData) error {
	signature := signVerification(rb.SessionId, rb.OrderId, rb.Amount, rb.Currency, c.crcKey)
	requestBody := verifyTransactionRequestBody{
		MerchantId: c.merchantID,
		PosId:      c.posID,
		SessionId:  rb.SessionId,
		Amount:     rb.Amount,
		Currency:   rb.Currency,
		OrderId:    rb.OrderId,
		Sign:       hex.EncodeToString(signature[:]),
	}

	b, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("transaction verification body marshalling error: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", c.url+c.transactionVerificationEndpoint, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("error while creating request to verify transaction: %w", err)
	}
	c.setHeaders(req)
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("error on making call to verify transaction: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error while reading body of verification response: %w", err)
	}
	r := &registerTransactionResponseBody{}
	err = json.Unmarshal(responseBody, r)
	if err != nil {
		return fmt.Errorf("transaction verification response unmarshalling error: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		// Tutaj warto sparsować body, żeby wyciągnąć kod błędu z P24
		return fmt.Errorf("p24 returned unexpected status: %d", resp.StatusCode)
	}

	return nil
}

func (c *APIClient) setHeaders(req *http.Request) {
	req.SetBasicAuth(fmt.Sprint(c.merchantID), c.reportKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
}

func getIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}

	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return strings.TrimSpace(r.RemoteAddr)
	}

	return strings.TrimSpace(host)
}

func isIPInList(ipStr string, trustedList []string) bool {
	// Trimujemy spacje na wypadek błędów w liście lub nagłówkach
	ipStr = strings.TrimSpace(ipStr)

	for _, trusted := range trustedList {
		if ipStr == strings.TrimSpace(trusted) {
			return true
		}
	}
	return false
}
