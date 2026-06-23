package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func makeSignature(secret, body string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(body))
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

func TestValidateSignature_Valid(t *testing.T) {
	secret := "test-secret-123"
	body := `{"action":"opened","number":1}`
	sig := makeSignature(secret, body)

	h := &Handler{WebhookSecret: secret}
	if !h.validateSignature([]byte(body), sig) {
		t.Error("expected valid signature to pass, but it failed")
	}
}

func TestValidateSignature_Invalid(t *testing.T) {
	h := &Handler{WebhookSecret: "test-secret-123"}
	if h.validateSignature([]byte(`{"action":"opened"}`), "sha256=invalidsignature") {
		t.Error("expected invalid signature to fail, but it passed")
	}
}

func TestValidateSignature_Missing(t *testing.T) {
	h := &Handler{WebhookSecret: "test-secret-123"}
	if h.validateSignature([]byte(`{"action":"opened"}`), "") {
		t.Error("expected empty signature to fail, but it passed")
	}
}

func TestValidateSignature_WrongSecret(t *testing.T) {
	body := `{"action":"opened","number":1}`
	sig := makeSignature("correct-secret", body)

	h := &Handler{WebhookSecret: "wrong-secret"}
	if h.validateSignature([]byte(body), sig) {
		t.Error("expected wrong secret to fail, but it passed")
	}
}

func TestValidateSignature_MalformedPrefix(t *testing.T) {
	h := &Handler{WebhookSecret: "test-secret-123"}
	if h.validateSignature([]byte(`{"action":"opened"}`), "md5=somehash") {
		t.Error("expected wrong prefix to fail, but it passed")
	}
}
