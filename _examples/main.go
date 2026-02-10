package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"sync"

	"github.com/daga93/przelewy24"
)

// This is full example of pacakge. Setup of Client, handler registration and notification receiving.
//
// When you replace configuration, you should be able to register transaction on you P24 panel.
// Just copy this example to any main func(). I suggest using Sandbox for tests...
// Remember that your server needs to be exposed - I use ngrok for that.
//
// Run code with `go run _examples/main.go`
// Run ngrok (or expose server)
// Click payment link (printed in console, should be  something like https://sandbox.przelewy24.pl/trnRequest/OAISJDIJ-A73647-E2R4DAS-8A177384FA2)
// Proceed with payment
// Wait for Notification
// Check your P24 Panel
func main() {
	crcKey := "your-crc-key-from-panel"       // Your crc key
	reportKey := "your-report-key-from-panel" // Your report key
	posId := 123456                           // Your client ID (from P24 Panel)

	ctx := context.Background()

	// Construct API Client for Przelewy24
	c := przelewy24.NewClient(
		posId,
		posId,
		crcKey,
		reportKey,
		przelewy24.WithShopDomain("https://example.com"),
	)

	// This jsut runs a simple server waiting for notification
	go runServer(
		c.NotificationWebhookHandler(
		// For this example, no real verification occurs.
		// Data from notifications is passed.
		// You should verify that amount, SessionID, Currency is
		// the same as in your datastore.
		// You should probably store OrderID as well.
		func(notification *przelewy24.NotificationBody) (przelewy24.VerificationData, error) {
			a := przelewy24.VerificationData{
				OrderId:   notification.OrderId,
				SessionId: notification.SessionId,
				Amount:    notification.Amount,
				Currency:  notification.Currency,
			}
			return a, nil
			}),
		c.RefundWebhookHandler(
			func(notification *przelewy24.RefundNotificationBody) (przelewy24.RefundVerificationData, error) {
				a := przelewy24.RefundVerificationData{
					SessionId:   notification.SessionId,
					OrderId:     notification.OrderId,
					Amount:      notification.Amount,
					Currency:    notification.Currency,
					RefundsUUID: notification.RefundsUuid,
				}
				return a, nil
			}),
	)

	// Here the real job is done - trasaction is registered.
	//
	// redirectionUrl, err := c.RegisterTransaction(
	// 	ctx,
	// 	100_00,
	// 	przelewy24.PLN,
	// 	"customer@example.com",
	// 	"Test 123",
	// 	przelewy24.CountryPL,
	// 	przelewy24.LanguagePL,
	// 	fmt.Sprintf("SO-%d", time.Now().UnixNano()),
	// )

	// Get the details by SessionID
	a, err := c.GetTransactionDetails(ctx, "SO-1769620150474463000")
	fmt.Println(a)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Begin Refund process. Remember that you should be accessible from internet
	// as the webhook is being sent from P24, to verify status.
	b, err := c.Refund(
		ctx,
		a.OrderId,
		a.SessionId,
		"abc-123-def-463",
		100,
		"Test refund",
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Success", b)

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

func runServer(
	handleFunc func(w http.ResponseWriter, r *http.Request),
	handleRefundFunc func(w http.ResponseWriter, r *http.Request),
) {
	slog.Info("Starting Server...")
	http.HandleFunc("/p24/payment/transaction/notification", handleFunc)
	http.HandleFunc("/p24/refund/notification", handleRefundFunc)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("Succesfully Redirected")) })
	log.Fatal(http.ListenAndServe(":8080", nil))
}
