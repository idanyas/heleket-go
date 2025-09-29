package heleket

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
)

// signRequest generates a signature for the request body according to Heleket's algorithm:
// MD5(base64(requestBody) + apiKey)
func (c *Heleket) signRequest(apiKey string, reqBody []byte) string {
	data := base64.StdEncoding.EncodeToString(reqBody)
	hash := md5.Sum([]byte(data + apiKey))
	return hex.EncodeToString(hash[:])
}

// VerifySign verifies the webhook signature according to Heleket's algorithm.
//
// The signature verification must preserve the exact JSON formatting (including key order)
// that Heleket used when generating the signature. This implementation uses raw string
// manipulation to remove the 'sign' field without re-marshaling, which would change key order.
//
// Algorithm:
// 1. Extract the 'sign' field value from the webhook JSON
// 2. Remove the '"sign":"value"' entry from the raw JSON string (preserving all other formatting)
// 3. Generate signature: MD5(base64(json_without_sign) + apiKey)
// 4. Compare with the extracted signature using constant-time comparison
//
// This matches Heleket's PHP implementation:
// $sign = $data['sign'];
// unset($data['sign']);
// $hash = md5(base64_encode(json_encode($data, JSON_UNESCAPED_UNICODE)) . $apiPaymentKey);
func (c *Heleket) VerifySign(apiKey string, reqBody []byte) error {
	// First, parse to extract the signature value
	var jsonBody map[string]any
	if err := json.Unmarshal(reqBody, &jsonBody); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	// Extract the signature from the webhook
	reqSign, ok := jsonBody["sign"].(string)
	if !ok || reqSign == "" {
		return errors.New("missing or invalid 'sign' field in webhook")
	}

	// Remove the "sign" field from the raw JSON string using regex
	// This preserves the exact formatting and key order of the original payload
	//
	// Pattern explanation:
	// "sign"\s*:\s*"(?:[^"\\]|\\.)*"\s*,?\s*
	// - "sign"                    : literal "sign" key
	// - \s*:\s*                   : colon with optional whitespace
	// - "(?:[^"\\]|\\.)*"         : JSON string value with proper escape handling
	//   - (?:...)                 : non-capturing group
	//   - [^"\\]                  : match any char except " or \
	//   - |                       : OR
	//   - \\.                     : match backslash followed by any char (handles escapes)
	//   - *                       : repeat zero or more times
	// - \s*,?\s*                  : optional trailing comma with whitespace
	signPattern := regexp.MustCompile(`"sign"\s*:\s*"(?:[^"\\]|\\.)*"\s*,?\s*`)
	bodyWithoutSign := signPattern.ReplaceAll(reqBody, []byte{})

	// Clean up any trailing commas that might be left after removing "sign"
	// This handles cases like: {"a":"b","sign":"x",} -> {"a":"b",}
	bodyWithoutSign = regexp.MustCompile(`,\s*}`).ReplaceAll(bodyWithoutSign, []byte("}"))
	bodyWithoutSign = regexp.MustCompile(`,\s*]`).ReplaceAll(bodyWithoutSign, []byte("]"))

	// Generate the expected signature using the cleaned raw JSON
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
