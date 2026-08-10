package gitcode

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// VerifyWebhookSignature verifies the signature of a webhook payload.
// The signature is expected in the format "sha256=<hex>".
//
// Example usage:
//
//	func handleWebhook(w http.ResponseWriter, r *http.Request) {
//	    payload, _ := io.ReadAll(r.Body)
//	    signature := r.Header.Get("X-Gitcode-Signature")
//	    if !gitcode.VerifyWebhookSignature(payload, "your-secret", signature) {
//	        http.Error(w, "Invalid signature", http.StatusUnauthorized)
//	        return
//	    }
//	    // Process webhook...
//	}
func VerifyWebhookSignature(payload []byte, secret, signature string) bool {
	if secret == "" || signature == "" {
		return false
	}

	// Remove algorithm prefix if present
	sig := strings.TrimPrefix(signature, "sha256=")
	if sig == signature && !strings.HasPrefix(signature, "sha256=") {
		// No sha256= prefix, try raw hex
		sig = signature
	}

	// Compute expected signature
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expected := hex.EncodeToString(mac.Sum(nil))

	// Constant-time comparison
	return hmac.Equal([]byte(strings.ToLower(sig)), []byte(strings.ToLower(expected)))
}

// ComputeWebhookSignature computes the HMAC-SHA256 signature for a payload.
// Returns the signature in the format "sha256=<hex>".
func ComputeWebhookSignature(payload []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}
