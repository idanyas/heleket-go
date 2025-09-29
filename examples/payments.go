package main

import (
	"fmt"
	"log"
	"time"

	"github.com/idanyas/go-heleket"
)

// RunPaymentExamples executes all examples related to the Payment API.
func RunPaymentExamples() {
	// 1. Create a new invoice. This is the central object for payments.
	log.Println("1. Creating a payment invoice...")
	orderID := fmt.Sprintf("sdk-test-%d", time.Now().Unix())
	invoice, err := createPaymentInvoice(orderID)
	if err != nil {
		log.Printf("Error creating invoice: %v\n", err)
		return // Stop if we can't create an invoice
	}

	// 2. Generate a QR code for the created invoice's payment address.
	log.Println("\n2. Generating QR code for the invoice...")
	generateInvoiceQRCode(invoice.UUID)

	// 3. Get information about the invoice we just created.
	log.Println("\n3. Getting payment info by UUID...")
	getPaymentInfo(invoice.UUID, "")

	log.Println("\n4. Getting payment info by Order ID...")
	getPaymentInfo("", invoice.OrderId)

	// 5. Get payment history for the last 24 hours.
	log.Println("\n5. Getting payment history...")
	getPaymentHistory()

	// 6. Get the list of all available payment services (currencies/networks).
	log.Println("\n6. Getting list of payment services...")
	getPaymentServices()
}

// createPaymentInvoice demonstrates how to create a payment invoice.
func createPaymentInvoice(orderID string) (*heleket.Payment, error) {
	req := &heleket.InvoiceRequest{
		Amount:   "10.50",
		Currency: "USD",
		OrderId:  orderID,
		InvoiceRequestOptions: &heleket.InvoiceRequestOptions{
			// To create an invoice for a specific crypto, provide ToCurrency and Network
			ToCurrency:  "USDT",
			Network:     "tron",
			UrlCallback: "https://your-backend.com/webhook/heleket",
			UrlSuccess:  "https://your-shop.com/payment/success",
			Lifetime:    3600, // 1 hour in seconds
		},
	}

	invoice, err := client.CreateInvoice(req)
	if err != nil {
		return nil, fmt.Errorf("CreateInvoice failed: %w", err)
	}

	log.Println("Invoice created successfully:")
	prettyPrint(invoice)
	return invoice, nil
}

// generateInvoiceQRCode shows how to generate a QR code for an invoice UUID.
func generateInvoiceQRCode(invoiceUUID string) {
	qrCodeBase64, err := client.GeneratePaymentQRCode(invoiceUUID)
	if err != nil {
		log.Printf("GeneratePaymentQRCode failed: %v", err)
		return
	}
	log.Println("Successfully generated QR code image (Base64 encoded):")
	fmt.Printf("%.100s...\n", qrCodeBase64) // Print first 100 chars
}

// getPaymentInfo demonstrates fetching invoice details by either UUID or Order ID.
func getPaymentInfo(uuid, orderID string) {
	req := &heleket.PaymentInfoRequest{
		PaymentUUID: uuid,
		OrderId:     orderID,
	}
	paymentInfo, err := client.GetPaymentInfo(req)
	if err != nil {
		log.Printf("GetPaymentInfo failed: %v", err)
		return
	}
	log.Println("Successfully retrieved payment info:")
	prettyPrint(paymentInfo)
}

// getPaymentHistory shows how to retrieve a list of past payments with pagination.
func getPaymentHistory() {
	dateTo := time.Now()
	dateFrom := dateTo.Add(-24 * time.Hour)

	history, err := client.GetPaymentHistory(dateFrom, dateTo, "")
	if err != nil {
		log.Printf("GetPaymentHistory failed: %v", err)
		return
	}

	log.Println("Successfully retrieved payment history (first page):")
	prettyPrint(history)

	// Example of fetching the next page using the cursor
	if history.Paginate != nil && history.Paginate.HasPages && history.Paginate.NextCursor != "" {
		log.Println("\nFetching next page of payment history...")
		nextPageHistory, err := client.GetPaymentHistory(dateFrom, dateTo, history.Paginate.NextCursor)
		if err != nil {
			log.Printf("GetPaymentHistory (next page) failed: %v", err)
			return
		}
		log.Println("Successfully retrieved next page:")
		prettyPrint(nextPageHistory)
	}
}

// getPaymentServices lists all available currencies and networks for payments.
func getPaymentServices() {
	services, err := client.GetPaymentServicesList()
	if err != nil {
		log.Printf("GetPaymentServicesList failed: %v", err)
		return
	}
	log.Println("Successfully retrieved list of payment services (showing first 2):")
	if len(services) > 2 {
		prettyPrint(services[:2])
	} else {
		prettyPrint(services)
	}
}
