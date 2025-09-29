package heleket

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
)

func (c *Heleket) signRequest(apiKey string, reqBody []byte) string {
	data := base64.StdEncoding.EncodeToString(reqBody)
	hash := md5.Sum([]byte(data + apiKey))
	return hex.EncodeToString(hash[:])
}

func (c *Heleket) VerifySign(apiKey string, reqBody []byte) error {
	var jsonBody map[string]any
	decoder := json.NewDecoder(bytes.NewReader(reqBody))
	decoder.UseNumber()
	if err := decoder.Decode(&jsonBody); err != nil {
		return err
	}

	reqSign, ok := jsonBody["sign"].(string)
	if !ok {
		return errors.New("missing signature field in request body")
	}
	delete(jsonBody, "sign")

	// Re-marshal the body without the 'sign' field.
	// This produces a canonical JSON with alphabetically sorted keys.
	// This assumes the server also signs a canonicalized JSON.
	bodyWithoutSign, err := json.Marshal(jsonBody)
	if err != nil {
		return err
	}

	// The webhook documentation warns that the source JSON (from PHP) may have
	// escaped forward slashes, and our signature generation must match that.
	// Go's json.Marshal does not escape them, so we do it manually.
	// This is a simplification and might not handle all edge cases of escaped characters
	// but is necessary to match the documented signature scheme for payloads containing URLs.
	bodyWithEscapedSlashes := bytes.Replace(bodyWithoutSign, []byte("/"), []byte("\\/"), -1)

	expectedSign := c.signRequest(apiKey, bodyWithEscapedSlashes)
	if reqSign != expectedSign {
		return errors.New("invalid signature")
	}
	return nil
}
