package main

import (
	"fmt"
	"log"
	"time"

	"github.com/idanyas/go-heleket"
)

// RunStaticWalletExamples executes all examples related to Static Wallets.
func RunStaticWalletExamples() {
	// 1. Create a static wallet for a user/order.
	log.Println("1. Creating a static wallet...")
	orderID := fmt.Sprintf("wallet-%d", time.Now().Unix())
	wallet, err := createStaticWallet(orderID)
	if err != nil {
		log.Printf("Error creating static wallet: %v\n", err)
		return
	}

	// 2. Generate a QR code for the new static wallet address.
	log.Println("\n2. Generating QR code for the static wallet...")
	generateStaticWalletQRCode(wallet.UUID)

	// 3. Block the static wallet address. Further payments will be held.
	log.Println("\n3. Blocking the static wallet address...")
	blockWalletAddress(wallet.OrderId)
}

// createStaticWallet shows how to generate a persistent crypto address for an order_id.
func createStaticWallet(orderID string) (*heleket.StaticWalletResponse, error) {
	req := &heleket.StaticWalletRequest{
		Currency: "USDT",
		Network:  "tron",
		OrderId:  orderID,
		StaticWalletRequestOptions: &heleket.StaticWalletRequestOptions{
			UrlCallback: "https://your-backend.com/webhook/static-wallet",
		},
	}

	wallet, err := client.CreateStaticWallet(req)
	if err != nil {
		return nil, fmt.Errorf("CreateStaticWallet failed: %w", err)
	}

	log.Println("Static wallet created successfully:")
	prettyPrint(wallet)
	return wallet, nil
}

// generateStaticWalletQRCode generates a QR code for a static wallet UUID.
func generateStaticWalletQRCode(walletUUID string) {
	qrCodeBase64, err := client.GenerateStaticWalletQRCode(walletUUID)
	if err != nil {
		log.Printf("GenerateStaticWalletQRCode failed: %v", err)
		return
	}
	log.Println("Successfully generated QR code image (Base64 encoded):")
	fmt.Printf("%.100s...\n", qrCodeBase64) // Print first 100 chars
}

// blockWalletAddress demonstrates how to block a static wallet.
func blockWalletAddress(orderID string) {
	req := &heleket.BlockAddressRequest{
		OrderId: orderID,
	}
	response, err := client.BlockAddress(req)
	if err != nil {
		log.Printf("BlockAddress failed: %v", err)
		return
	}
	log.Println("Block address request successful:")
	prettyPrint(response)
}
