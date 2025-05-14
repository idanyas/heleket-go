package heleket

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const apiUrl = "https://api.heleket.com/v1"

type Heleket struct {
	merchant      string
	paymentApiKey string
	payoutApiKey  string
	client        *http.Client
}

func New(client *http.Client, merchant, paymentApiKey, payoutApiKey string) *Heleket {
	return &Heleket{
		client:        client,
		merchant:      merchant,
		paymentApiKey: paymentApiKey,
		payoutApiKey:  payoutApiKey,
	}
}

func (c *Heleket) fetch(method string, endpoint string, payload any) (*http.Response, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	sign := c.signRequest(c.paymentApiKey, body)
	req, err := http.NewRequest(method, apiUrl+endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("merchant", c.merchant)
	req.Header.Set("sign", sign)
	res, err := c.client.Do(req)
	return res, err
}
