package heleket

import (
	"encoding/json"
	"errors"
	"net/url"
	"time"
)

const (
	createInvoiceEndpoit          = "/payment"
	generateInvoiceQRCodeEndpoint = "/payment/qr"
	paymentInfoEndpoint           = "/payment/info"
	paymentHistoryEndpoint        = "/payment/list"
	paymentServicesListEndpoint   = "/payment/services"
)

type InvoiceRequest struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
	OrderId  string `json:"order_id"`
	*InvoiceRequestOptions
}

type InvoiceRequestOptions struct {
	Network                string     `json:"network,omitempty"`
	UrlReturn              string     `json:"url_return,omitempty"`
	UrlSuccess             string     `json:"url_success,omitempty"`
	UrlCallback            string     `json:"url_callback,omitempty"`
	IsPaymentMultiple      bool       `json:"is_payment_multiple,omitempty"`
	Lifetime               uint16     `json:"lifetime,omitempty"`
	ToCurrency             string     `json:"to_currency,omitempty"`
	Subtract               uint8      `json:"subtract,omitempty"`
	AccuracyPaymentPercent uint8      `json:"accuracy_payment_percent,omitempty"`
	AdditionalData         string     `json:"additional_data,omitempty"`
	Currencies             []Currency `json:"currencies,omitempty"`
	ExceptCurrencies       []Currency `json:"except_currencies,omitempty"`
	CourseSource           string     `json:"course_source,omitempty"`
	FromReferralCode       string     `json:"from_referral_code,omitempty"`
	DiscountPercent        int8       `json:"discount_percent,omitempty"`
	IsRefresh              bool       `json:"is_refresh,omitempty"`
	PayerEmail             string     `json:"payer_email,omitempty"`
}

type Currency struct {
	Currency string `json:"currency"`
	Network  string `json:"network,omitempty"`
}

type PaymentConvert struct {
	ToCurrency string `json:"to_currency"`
	Commission string `json:"commission"`
	Rate       string `json:"rate"`
	Amount     string `json:"amount"`
}

type Payment struct {
	UUID                    string          `json:"uuid"`
	OrderId                 string          `json:"order_id"`
	Amount                  string          `json:"amount"`
	PaymentAmount           string          `json:"payment_amount,omitempty"`
	PaymentAmountUSD        string          `json:"payment_amount_usd,omitempty"`
	PayerAmount             string          `json:"payer_amount,omitempty"`
	PayerAmountExchangeRate string          `json:"payer_amount_exchange_rate,omitempty"`
	DiscountPercent         int8            `json:"discount_percent,omitempty"`
	Discount                string          `json:"discount,omitempty"`
	PayerCurrency           string          `json:"payer_currency,omitempty"`
	Currency                string          `json:"currency"`
	MerchantAmount          string          `json:"merchant_amount,omitempty"`
	Network                 string          `json:"network,omitempty"`
	Address                 string          `json:"address,omitempty"`
	From                    string          `json:"from,omitempty"`
	TxId                    string          `json:"txid,omitempty"`
	PaymentStatus           string          `json:"payment_status"`
	Status                  string          `json:"status,omitempty"`
	Url                     string          `json:"url"`
	ExpiredAt               float64         `json:"expired_at"`
	IsFinal                 bool            `json:"is_final"`
	AdditionalData          string          `json:"additional_data,omitempty"`
	Comments                string          `json:"comments,omitempty"`
	CreatedAt               time.Time       `json:"created_at"`
	UpdatedAt               time.Time       `json:"updated_at"`
	AddressQrCode           string          `json:"address_qr_code,omitempty"`
	Commission              string          `json:"commission,omitempty"`
	Convert                 *PaymentConvert `json:"convert,omitempty"`
}

type invoiceRawResponse struct {
	Result *Payment
	State  int8
}

type paymentQRCodeRawResponse struct {
	Result struct {
		Image string `json:"image"`
	} `json:"result"`
	State int8 `json:"state"`
}

type PaymentInfoRequest struct {
	PaymentUUID string `json:"uuid,omitempty"`
	OrderId     string `json:"order_id,omitempty"`
}

type PaymentHistoryResponse struct {
	Payments []*Payment
	Paginate *PaymentHistoryPaginate
}

type PaymentHistoryPaginate struct {
	Count          int16  `json:"count"`
	HasPages       bool   `json:"hasPages"`
	NextCursor     string `json:"nextCursor,omitempty"`
	PreviousCursor string `json:"previousCursor,omitempty"`
	PerPage        int16  `json:"perPage"`
}

type paymentHistoryRawResponse struct {
	State  int8 `json:"state"`
	Result struct {
		Items    []*Payment              `json:"items"`
		Paginate *PaymentHistoryPaginate `json:"paginate"`
	} `json:"result"`
}

type PaymentService struct {
	Network     string                    `json:"network"`
	Currency    string                    `json:"currency"`
	IsAvailable bool                      `json:"is_available"`
	Limit       *PaymentServiceLimit      `json:"limit"`
	Commission  *PaymentServiceCommission `json:"commission"`
}

type PaymentServiceLimit struct {
	MinAmount string `json:"min_amount"`
	MaxAmount string `json:"max_amount"`
}

type PaymentServiceCommission struct {
	FeeAmount string `json:"fee_amount"`
	Percent   string `json:"percent"`
}

type paymentServiceListRawResponse struct {
	Result []*PaymentService `json:"result"`
	State  int8              `json:"state"`
}

func (c *Heleket) CreateInvoice(invoiceReq *InvoiceRequest) (*Payment, error) {
	res, err := c.fetch("POST", createInvoiceEndpoit, invoiceReq, c.paymentApiKey)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	response := &invoiceRawResponse{}
	if err = json.NewDecoder(res.Body).Decode(response); err != nil {
		return nil, err
	}

	return response.Result, nil
}

func (c *Heleket) GeneratePaymentQRCode(paymentUUID string) (string, error) {
	payload := map[string]any{"merchant_payment_uuid": paymentUUID}
	res, err := c.fetch("POST", generateInvoiceQRCodeEndpoint, payload, c.paymentApiKey)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	response := &paymentQRCodeRawResponse{}
	if err = json.NewDecoder(res.Body).Decode(response); err != nil {
		return "", err
	}

	return response.Result.Image, nil

}

func (c *Heleket) GetPaymentInfo(paymentInfoReq *PaymentInfoRequest) (*Payment, error) {
	if paymentInfoReq.PaymentUUID == "" && paymentInfoReq.OrderId == "" {
		return nil, errors.New("you should pass one of required values [PaymentUUID, OrderId]")
	}

	res, err := c.fetch("POST", paymentInfoEndpoint, paymentInfoReq, c.paymentApiKey)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	response := &invoiceRawResponse{}
	if err = json.NewDecoder(res.Body).Decode(response); err != nil {
		return nil, err
	}

	return response.Result, nil
}

func (c *Heleket) GetPaymentHistory(dateFrom, dateTo time.Time, cursor string) (*PaymentHistoryResponse, error) {
	const timeFormat = "2006-01-02 15:04:05"
	payload := map[string]any{"date_from": dateFrom.Format(timeFormat), "date_to": dateTo.Format(timeFormat)}

	endpoint := paymentHistoryEndpoint
	if cursor != "" {
		endpoint += "?cursor=" + url.QueryEscape(cursor)
	}

	res, err := c.fetch("POST", endpoint, payload, c.paymentApiKey)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	response := &paymentHistoryRawResponse{}
	if err = json.NewDecoder(res.Body).Decode(response); err != nil {
		return nil, err
	}

	paymentHistory := &PaymentHistoryResponse{
		Payments: response.Result.Items,
		Paginate: response.Result.Paginate,
	}
	return paymentHistory, nil
}

func (c *Heleket) GetPaymentServicesList() ([]*PaymentService, error) {
	payload := make(map[string]any)
	res, err := c.fetch("POST", paymentServicesListEndpoint, payload, c.paymentApiKey)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	response := &paymentServiceListRawResponse{}
	if err = json.NewDecoder(res.Body).Decode(response); err != nil {
		return nil, err
	}

	return response.Result, nil
}
