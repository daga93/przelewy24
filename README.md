This package provides an API Client for Przelewy24 payment gateway.
When customer want to pay, he/she gets redirected to Przelewy24 payment panel. There a bunch of payments methods are available.
Upon successful payment, webhook notification from Przelewy24 is done to your server (this package provides handlers for it).

## How does this module work?
1. Create API Client:
```go
	client := przelewy24.NewClient(
		posId,
		posId,
		crcKey,
		reportKey,
		WithShopDomain("https://example.com"),
	)
```
2. Add Notification Handler in your server:
```go
handler := client.NotificationWebhookHandler(
    func(notification *NotificationOnTransactionRequest) (VerificationData, error) {
        d := przelewy24.VerificationData{} 
        // Your code that gets VerificationData from database here
        return d, nil
    }
)

// Register handler, for instance:
http.HandleFunc("/p24/payment/transaction/notification", handler)
```
3. Use RegisterTransaction in payment step:
```go
	redirectionUrl, err := c.RegisterTransaction(
		ctx,                                            // Context
		100_00,                                         // Amount
		przelewy24.PLN,                                 // Currency
		"customer@example.com",                         // Email
		"Test 123",                                     // Description
        przelewy24.CountryPL,                           // Country
		przelewy24.LanguagePL,                          // Language
		fmt.Sprintf("SO-%d", time.Now().UnixNano()),    // SessionID
        WithPhone("123456789"),                         // Other - optional params
	)
```
And that is it. If you have your notification endpoint setup correctly, you'll receive Webhook Notification in few minutes. You'll see transaction in P24 Panel confirmed.

## About this repo
This is non comercial package build in spare time. I do not guarantee that it will be maintained - but if you need something fixed/added contact me.
Currently only Payments are available. I plan to add refunds API in near future.

## TODOs
- Cover with tests
- Add some build
- Implement Refunds
