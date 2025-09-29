package tests

import (
	"testing"
	"time"

	"github.com/idanyas/go-heleket"

	"github.com/stretchr/testify/require"
)

func createTestInvoice(t *testing.T) *heleket.Payment {
	invoiceReq := &heleket.InvoiceRequest{
		Amount:   "10",
		Currency: "USD",
		OrderId:  "xxy",
		InvoiceRequestOptions: &heleket.InvoiceRequestOptions{
			Network:     "tron",
			UrlCallback: "https://example.com/heleket/callback",
		},
	}
	invoice, err := TestHeleket.CreateInvoice(invoiceReq)
	require.NoError(t, err)
	require.NotEmpty(t, invoice)

	return invoice
}

func TestCreateInvoice(t *testing.T) {
	createTestInvoice(t)
}

func TestGenerateInvoiceQRCode(t *testing.T) {
	invoice := createTestInvoice(t)
	qrCode, err := TestHeleket.GeneratePaymentQRCode(invoice.UUID)
	require.NoError(t, err)
	require.NotEmpty(t, qrCode)
}

func TestGetPaymentInfo(t *testing.T) {
	invoice := createTestInvoice(t)
	payment, err := TestHeleket.GetPaymentInfo(&heleket.PaymentInfoRequest{PaymentUUID: invoice.UUID})
	require.NoError(t, err)
	require.NotEmpty(t, payment)
}

func TestGetPaymentHistory(t *testing.T) {
	payments, err := TestHeleket.GetPaymentHistory(time.Now(), time.Now(), "")
	require.NoError(t, err)
	require.NotEmpty(t, payments)
}
