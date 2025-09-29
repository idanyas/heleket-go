package main

import (
	"fmt"
	"log"
	"time"

	"github.com/idanyas/go-heleket"
)

// RunPayoutExamples executes all examples related to the Payout API.
func RunPayoutExamples() {
	// 1. Create a payout request.
	log.Println("1. Creating a payout request...")
	payout, err := createPayout()
	if err != nil {
		log.Printf("Error creating payout: %v\n", err)
		// Don't return, so other examples can still run
	}

	// 2. Get information about the payout (if created successfully).
	if payout != nil {
		log.Println("\n2. Getting payout info by UUID...")
		getPayoutInfo(payout.UUID, "")
	}

	// 3. Get payout history.
	log.Println("\n3. Getting payout history...")
	getPayoutHistory()

	// 4. Get the list of all available payout services.
	log.Println("\n4. Getting list of payout services...")
	getPayoutServices()
}

// createPayout demonstrates creating a payout to a specified address.
// NOTE: This will attempt to make a real payout from your balance.
// It will likely fail if your balance is insufficient.
func createPayout() (*heleket.Payout, error) {
	req := &heleket.PayoutRequest{
		Amount:     "1.2",
		Currency:   "USDT",
		Network:    "tron",
		OrderId:    fmt.Sprintf("payout-%d", time.Now().Unix()),
		Address:    "TTEtddVZyNtLD9wbq4PzomjBhtxenSMXbb", // Use a valid recipient address
		IsSubtract: true,                                 // Commission will be subtracted from the amount
	}

	payout, err := client.CreatePayout(req)
	if err != nil {
		return nil, fmt.Errorf("CreatePayout failed: %w", err)
	}

	log.Println("Payout request created successfully:")
	prettyPrint(payout)
	return payout, nil
}

// getPayoutInfo shows how to fetch details of a specific payout.
func getPayoutInfo(uuid, orderID string) {
	req := &heleket.PayoutInfoRequest{
		PayoutUUID: uuid,
		OrderId:    orderID,
	}
	payoutInfo, err := client.GetPayoutInfo(req)
	if err != nil {
		log.Printf("GetPayoutInfo failed: %v", err)
		return
	}
	log.Println("Successfully retrieved payout info:")
	prettyPrint(payoutInfo)
}

// getPayoutHistory retrieves a list of past payouts.
func getPayoutHistory() {
	dateTo := time.Now()
	dateFrom := dateTo.Add(-24 * time.Hour)
	history, err := client.GetPayoutHistory(dateFrom, dateTo)
	if err != nil {
		log.Printf("GetPayoutHistory failed: %v", err)
		return
	}
	log.Println("Successfully retrieved payout history:")
	prettyPrint(history)
}

// getPayoutServices lists all available currencies and networks for payouts.
func getPayoutServices() {
	services, err := client.GetPayoutServicesList()
	if err != nil {
		log.Printf("GetPayoutServicesList failed: %v", err)
		return
	}
	log.Println("Successfully retrieved list of payout services (showing first 2):")
	if len(services) > 2 {
		prettyPrint(services[:2])
	} else {
		prettyPrint(services)
	}
}
