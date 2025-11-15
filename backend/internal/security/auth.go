package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"
)

var ErrInvalidToken = errors.New("invalid or expired token")

type TokenData struct {
	UserID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

func GenerateDeepLinkToken(userID string, secret []byte, ttl time.Duration) (string, error) {
	data := TokenData{
		UserID:    userID,
		ExpiresAt: time.Now().Add(ttl),
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	encoded := base64.URLEncoding.EncodeToString(payload)
	sig := signHMAC([]byte(encoded), secret)
	token := encoded + "." + sig

	return token, nil
}

func VerifyDeepLinkToken(token string, secret []byte) (*TokenData, error) {
	var payload, sig string
	for i := len(token) - 1; i >= 0; i-- {
		if token[i] == '.' {
			payload = token[:i]
			sig = token[i+1:]
			break
		}
	}

	if payload == "" || sig == "" {
		return nil, ErrInvalidToken
	}

	expectedSig := signHMAC([]byte(payload), secret)
	if !hmac.Equal([]byte(sig), []byte(expectedSig)) {
		return nil, ErrInvalidToken
	}

	decoded, err := base64.URLEncoding.DecodeString(payload)
	if err != nil {
		return nil, ErrInvalidToken
	}

	var data TokenData
	if err := json.Unmarshal(decoded, &data); err != nil {
		return nil, ErrInvalidToken
	}

	if time.Now().After(data.ExpiresAt) {
		return nil, ErrInvalidToken
	}

	return &data, nil
}

func signHMAC(data, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write(data)
	return base64.URLEncoding.EncodeToString(mac.Sum(nil))
}

func VerifyWebhookSignature(secret, signature string, body []byte) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expected))
}
