package main

import (
	"log"

	"github.com/idanyas/go-heleket"
)

// RunRefundExamples executes examples for refunding payments.
// NOTE: These operations will fail unless the target invoice/wallet
// has a refundable balance.
func RunRefundExamples() {
	// 1. Issue a refund for a completed payment.
	log.Println("1. Attempting to refund a payment...")
	refundPayment("some-order-id-that-was-paid") // Replace with a real, paid order ID

	// 2. Issue a refund from a blocked static wallet.
	log.Println("\n2. Attempting to refund from a blocked address...")
	refundFromBlockedAddress("some-wallet-order-id") // Replace with a real, blocked wallet order ID
}

// refundPayment demonstrates refunding a specific invoice.
func refundPayment(orderID string) {
	req := &heleket.RefundRequest{
		OrderId:    orderID,
		Address:    "TDD97yguPESTpcrJMqU6h2ozZbibv4Vaqm", // Address to send the refund to
		IsSubtract: true,                                 // Deduct commission from merchant balance
	}
	success, err := client.Refund(req)
	if err != nil {
		log.Printf("Refund failed (this is expected if the invoice is not refundable): %v", err)
		return
	}
	log.Printf("Refund request successful: %v", success)
}

// refundFromBlockedAddress shows how to refund funds from a blocked static wallet.
func refundFromBlockedAddress(orderID string) {
	req := &heleket.BlockedAddressRefundRequest{
		OrderId: orderID,
		Address: "TK8...", // Address to send the refund to
	}
	response, err := client.BlockedAddressRefund(req)
	if err != nil {
		log.Printf("BlockedAddressRefund failed (this is expected if the wallet has no funds to refund): %v", err)
		return
	}
	log.Println("Blocked address refund request successful:")
	prettyPrint(response)
}
