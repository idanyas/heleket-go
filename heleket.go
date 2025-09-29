package heleket

import (
	"bytes"
	"encoding/json"
	"io"
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

func (c *Heleket) fetch(method string, endpoint string, payload any, apiKey string) (*http.Response, error) {
	var body []byte
	var err error

	// For GET requests, payload should be nil. Signature is on an empty string.
	// For POST requests with no parameters, payload should be an empty map or struct, which marshals to "{}".
	if payload != nil {
		body, err = json.Marshal(payload)
		if err != nil {
			return nil, err
		}
	} else {
		body = []byte("")
	}

	sign := c.signRequest(apiKey, body)

	var reqBody io.Reader
	if method != http.MethodGet {
		reqBody = bytes.NewBuffer(body)
	}

	req, err := http.NewRequest(method, apiUrl+endpoint, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("merchant", c.merchant)
	req.Header.Set("sign", sign)
	res, err := c.client.Do(req)
	return res, err
}
