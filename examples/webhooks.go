package main

import (
	"log"

	"github.com/idanyas/heleket-go"
)

// A sample webhook payload from the Heleket documentation.
// The signature is for demonstration and will not match your API key.
var sampleWebhookPayload = []byte(`
{
  "type": "payment",
  "uuid": "62f88b36-a9d5-4fa6-aa26-e040c3dbf26d",
  "order_id": "97a75bf8eda5cca41ba9d2e104840fcd",
  "amount": "3.00000000",
  "payment_amount": "3.00000000",
  "payment_amount_usd": "0.23",
  "merchant_amount": "2.94000000",
  "commission": "0.06000000",
  "is_final": true,
  "status": "paid",
  "from": "THgEWubVc8tPKXLJ4VZ5zbiiAK7AgqSeGH",
  "wallet_address_uuid": null,
  "network": "tron",
  "currency": "TRX",
  "payer_currency": "TRX",
  "additional_data": null,
  "convert": {
    "to_currency": "USDT",
    "commission": null,
    "rate": "0.07700000",
    "amount": "0.22638000"
  },
  "txid": "6f0d9c8374db57cac0d806251473de754f361c83a03cd805f74aa9da3193486b",
  "sign": "a76c0d77f3e8e1a419b138af04ab600a"
}
`)

// RunWebhookExamples executes examples for handling and testing webhooks.
func RunWebhookExamples() {
	// 1. Parse an incoming webhook payload without signature verification.
	log.Println("1. Parsing a webhook without signature verification...")
	parseWebhook(false)

	// 2. Parse a webhook with signature verification (expected to fail with sample data).
	log.Println("\n2. Parsing a webhook with signature verification (expected to fail)...")
	parseWebhook(true)

	// 3. Resend a webhook for a specific payment.
	log.Println("\n3. Resending a webhook...")
	resendWebhook("some-order-id") // Replace with a real order ID

	// Use a service like https://webhook.site to get a test URL
	testWebhookURL := "https://webhook.site/your-unique-url"

	// 4. Send a test payment webhook to a specified URL.
	log.Println("\n4. Sending a test payment webhook...")
	testPaymentWebhook(testWebhookURL)

	// 5. Send a test payout webhook to a specified URL.
	log.Println("\n5. Sending a test payout webhook...")
	testPayoutWebhook(testWebhookURL)

	// 6. Send a test static wallet webhook to a specified URL.
	log.Println("\n6. Sending a test static wallet webhook...")
	testWalletWebhook(testWebhookURL)
}

// parseWebhook shows how to parse a raw webhook body and optionally verify its signature.
func parseWebhook(verifySign bool) {
	webhook, err := client.ParseWebhook(sampleWebhookPayload, verifySign)
	if err != nil {
		log.Printf("ParseWebhook failed: %v", err)
		return
	}
	log.Println("Webhook parsed successfully:")
	prettyPrint(webhook)
}

// resendWebhook demonstrates how to request a webhook to be sent again.
func resendWebhook(orderID string) {
	req := &heleket.ResendWebhookRequest{
		OrderId: orderID,
	}
	success, err := client.ResendWebhook(req)
	if err != nil {
		log.Printf("ResendWebhook failed (expected if order_id is invalid or has no callback URL): %v", err)
		return
	}
	log.Printf("Webhook resend request successful: %v", success)
}

// testPaymentWebhook sends a simulated webhook to your endpoint for testing.
func testPaymentWebhook(callbackURL string) {
	req := &heleket.TestWebhookRequest{
		UrlCallback: callbackURL,
		Currency:    "USDT",
		Network:     "tron",
		Status:      "paid",
	}
	response, err := client.TestPaymentWebhook(req)
	if err != nil {
		log.Printf("TestPaymentWebhook failed: %v", err)
		return
	}
	log.Println("Test payment webhook sent successfully:")
	prettyPrint(response)
}

// testPayoutWebhook sends a simulated payout webhook to your endpoint for testing.
func testPayoutWebhook(callbackURL string) {
	req := &heleket.TestWebhookRequest{
		UrlCallback: callbackURL,
		Currency:    "USDT",
		Network:     "tron",
		Status:      "paid",
	}
	response, err := client.TestPayoutWebhook(req)
	if err != nil {
		log.Printf("TestPayoutWebhook failed: %v", err)
		return
	}
	log.Println("Test payout webhook sent successfully:")
	prettyPrint(response)
}

// testWalletWebhook sends a simulated wallet webhook to your endpoint for testing.
func testWalletWebhook(callbackURL string) {
	req := &heleket.TestWebhookRequest{
		UrlCallback: callbackURL,
		Currency:    "USDT",
		Network:     "tron",
		Status:      "paid",
	}
	response, err := client.TestWalletWebhook(req)
	if err != nil {
		log.Printf("TestWalletWebhook failed: %v", err)
		return
	}
	log.Println("Test wallet webhook sent successfully:")
	prettyPrint(response)
}
