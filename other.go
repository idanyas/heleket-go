package heleket

import (
	"encoding/json"
	"fmt"
)

const (
	balanceEndpoint             = "/balance"
	discountListEndpoint        = "/payment/discount/list"
	discountSetEndpoint         = "/payment/discount/set"
	exchangeRateListEndpointFmt = "/exchange-rate/%s/list"
)

// Balance
type WalletBalance struct {
	UUID         string `json:"uuid"`
	Balance      string `json:"balance"`
	CurrencyCode string `json:"currency_code"`
}

type BalanceInfo struct {
	Merchant []*WalletBalance `json:"merchant"`
	User     []*WalletBalance `json:"user"`
}

type balanceRawResponse struct {
	State  int8 `json:"state"`
	Result []struct {
		Balance *BalanceInfo `json:"balance"`
	} `json:"result"`
}

// Discount
type Discount struct {
	Network  string `json:"network"`
	Currency string `json:"currency"`
	Discount int8   `json:"discount"`
}

type discountListRawResponse struct {
	State  int8        `json:"state"`
	Result []*Discount `json:"result"`
}

type SetDiscountRequest struct {
	Currency        string `json:"currency"`
	Network         string `json:"network"`
	DiscountPercent int8   `json:"discount_percent"`
}

type setDiscountRawResponse struct {
	State  int8      `json:"state"`
	Result *Discount `json:"result"`
}

// Exchange Rate
type ExchangeRate struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Course string `json:"course"`
}

type exchangeRateListRawResponse struct {
	State  int8            `json:"state"`
	Result []*ExchangeRate `json:"result"`
}

func (c *Heleket) GetBalance() (*BalanceInfo, error) {
	res, err := c.fetch("POST", balanceEndpoint, make(map[string]any), c.paymentApiKey)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	response := &balanceRawResponse{}
	if err = json.NewDecoder(res.Body).Decode(response); err != nil {
		return nil, err
	}

	if len(response.Result) > 0 {
		return response.Result[0].Balance, nil
	}

	return nil, nil
}

func (c *Heleket) GetDiscountsList() ([]*Discount, error) {
	res, err := c.fetch("POST", discountListEndpoint, make(map[string]any), c.paymentApiKey)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	response := &discountListRawResponse{}
	if err = json.NewDecoder(res.Body).Decode(response); err != nil {
		return nil, err
	}

	return response.Result, nil
}

func (c *Heleket) SetDiscount(req *SetDiscountRequest) (*Discount, error) {
	res, err := c.fetch("POST", discountSetEndpoint, req, c.paymentApiKey)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	response := &setDiscountRawResponse{}
	if err = json.NewDecoder(res.Body).Decode(response); err != nil {
		return nil, err
	}

	return response.Result, nil
}

func (c *Heleket) GetExchangeRates(currency string) ([]*ExchangeRate, error) {
	endpoint := fmt.Sprintf(exchangeRateListEndpointFmt, currency)
	// GET request with no body, so payload is nil
	res, err := c.fetch("GET", endpoint, nil, c.paymentApiKey)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	response := &exchangeRateListRawResponse{}
	if err = json.NewDecoder(res.Body).Decode(response); err != nil {
		return nil, err
	}

	return response.Result, nil
}
