package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/idanyas/go-heleket"
)

var (
	// client is the shared Heleket SDK client instance.
	client *heleket.Heleket
)

func main() {
	// --- Configuration ---
	// IMPORTANT: Replace with your actual credentials.
	// It's recommended to use environment variables for security.
	merchantID := os.Getenv("HELEKET_MERCHANT_ID")
	paymentAPIKey := os.Getenv("HELEKET_PAYMENT_API_KEY")
	payoutAPIKey := os.Getenv("HELEKET_PAYOUT_API_KEY")

	if merchantID == "" || paymentAPIKey == "" || payoutAPIKey == "" {
		log.Fatal("Please set HELEKET_MERCHANT_ID, HELEKET_PAYMENT_API_KEY, and HELEKET_PAYOUT_API_KEY environment variables.")
	}

	// Initialize the Heleket client
	httpClient := &http.Client{}
	client = heleket.New(httpClient, merchantID, paymentAPIKey, payoutAPIKey)

	log.Println("--- Running Payment API Examples ---")
	RunPaymentExamples()

	log.Println("\n--- Running Payout API Examples ---")
	RunPayoutExamples()

	log.Println("\n--- Running Static Wallet API Examples ---")
	RunStaticWalletExamples()

	log.Println("\n--- Running Refund API Examples ---")
	RunRefundExamples()

	log.Println("\n--- Running Webhook API Examples ---")
	RunWebhookExamples()

	log.Println("\n--- Running Other Features Examples ---")
	RunOtherFeaturesExamples()
}

// prettyPrint is a helper function to print structs in a readable JSON format.
func prettyPrint(v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Printf("Error pretty printing: %v", err)
		fmt.Printf("%+v\n", v) // fallback to default print
		return
	}
	fmt.Println(string(b))
}
