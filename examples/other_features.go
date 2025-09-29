package main

import (
	"log"

	"github.com/idanyas/heleket-go"
)

// RunOtherFeaturesExamples executes examples for miscellaneous API endpoints.
func RunOtherFeaturesExamples() {
	// 1. Get your merchant and user account balances.
	log.Println("1. Getting account balance...")
	getBalances()

	// 2. Get the list of configured payment discounts.
	log.Println("\n2. Getting list of discounts...")
	getDiscounts()

	// 3. Set a new discount for a specific payment method.
	log.Println("\n3. Setting a discount...")
	setDiscount()

	// 4. Get current exchange rates for a currency.
	log.Println("\n4. Getting exchange rates for BTC...")
	getExchangeRates("BTC")
}

// getBalances fetches balances for all currencies in your account.
func getBalances() {
	balance, err := client.GetBalance()
	if err != nil {
		log.Printf("GetBalance failed: %v", err)
		return
	}
	log.Println("Successfully retrieved balances:")
	prettyPrint(balance)
}

// getDiscounts lists all currently active payment method discounts.
func getDiscounts() {
	discounts, err := client.GetDiscountsList()
	if err != nil {
		log.Printf("GetDiscountsList failed: %v", err)
		return
	}
	log.Println("Successfully retrieved discounts list:")
	prettyPrint(discounts)
}

// setDiscount shows how to apply a discount or commission to a payment method.
func setDiscount() {
	req := &heleket.SetDiscountRequest{
		Currency:        "USDT",
		Network:         "bsc",
		DiscountPercent: -5, // -5 means a 5% commission is added for the user
	}
	discount, err := client.SetDiscount(req)
	if err != nil {
		log.Printf("SetDiscount failed: %v", err)
		return
	}
	log.Println("Successfully set discount:")
	prettyPrint(discount)
}

// getExchangeRates fetches the latest conversion rates for a given currency.
func getExchangeRates(currency string) {
	rates, err := client.GetExchangeRates(currency)
	if err != nil {
		log.Printf("GetExchangeRates failed: %v", err)
		return
	}
	log.Printf("Successfully retrieved exchange rates for %s:", currency)
	prettyPrint(rates)
}
