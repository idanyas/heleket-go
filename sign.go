package heleket

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
)

// signRequest generates a signature for the request body according to Heleket's algorithm:
// MD5(base64(requestBody) + apiKey)
func (c *Heleket) signRequest(apiKey string, reqBody []byte) string {
	data := base64.StdEncoding.EncodeToString(reqBody)
	hash := md5.Sum([]byte(data + apiKey))
	return hex.EncodeToString(hash[:])
}

// VerifySign verifies the webhook signature according to Heleket's algorithm:
// 1. Parse the webhook JSON
// 2. Extract and remove the 'sign' field
// 3. Re-encode the remaining data to JSON (without 'sign')
// 4. Generate signature: MD5(base64(json_without_sign) + apiKey)
// 5. Compare with the extracted signature using constant-time comparison
//
// This matches Heleket's PHP implementation:
// $sign = $data['sign'];
// unset($data['sign']);
// $hash = md5(base64_encode(json_encode($data, JSON_UNESCAPED_UNICODE)) . $apiPaymentKey);
func (c *Heleket) VerifySign(apiKey string, reqBody []byte) error {
	// Parse the webhook body into a map
	var jsonBody map[string]any
	if err := json.Unmarshal(reqBody, &jsonBody); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	// Extract the signature from the webhook
	reqSign, ok := jsonBody["sign"].(string)
	if !ok {
		return errors.New("missing or invalid 'sign' field in webhook")
	}

	// Remove the 'sign' field before computing the expected signature
	delete(jsonBody, "sign")

	// Re-encode the body without the 'sign' field
	// Go's json.Marshal produces the correct format that matches Heleket's server:
	// - No escaped forward slashes (unlike PHP's default json_encode)
	// - No escaped Unicode (matches PHP's JSON_UNESCAPED_UNICODE flag)
	bodyWithoutSign, err := json.Marshal(jsonBody)
	if err != nil {
		return fmt.Errorf("failed to re-encode JSON: %w", err)
	}

	// Generate the expected signature
	expectedSign := c.signRequest(apiKey, bodyWithoutSign)

	// Use constant-time comparison to prevent timing attacks
	if !constantTimeEqual(expectedSign, reqSign) {
		return fmt.Errorf("invalid signature: expected=%s received=%s", expectedSign, reqSign)
	}

	return nil
}

// constantTimeEqual performs constant-time string comparison to prevent timing attacks
func constantTimeEqual(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	var result byte
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}
	return result == 0
}
