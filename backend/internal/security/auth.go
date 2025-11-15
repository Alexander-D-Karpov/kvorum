package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

var ErrInvalidToken = errors.New("invalid or expired token")

type TokenData struct {
	UserID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

type WebAppInitData struct {
	QueryID    string     `json:"query_id"`
	User       WebAppUser `json:"user"`
	AuthDate   int64      `json:"auth_date"`
	Hash       string     `json:"hash"`
	StartParam string     `json:"start_param,omitempty"`
}

type WebAppUser struct {
	ID           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
	PhotoURL     string `json:"photo_url"`
}

func ValidateWebAppData(initData string, botToken string) (*WebAppInitData, error) {
	decoded, err := url.QueryUnescape(initData)
	if err != nil {
		return nil, errors.New("failed to decode init data")
	}

	values, err := url.ParseQuery(decoded)
	if err != nil {
		return nil, errors.New("failed to parse init data")
	}

	receivedHash := values.Get("hash")
	if receivedHash == "" {
		return nil, errors.New("hash not found in init data")
	}

	var pairs []string
	for key, vals := range values {
		if key == "hash" {
			continue
		}
		for _, val := range vals {
			pairs = append(pairs, key+"="+val)
		}
	}

	sort.Strings(pairs)
	dataCheckString := strings.Join(pairs, "\n")

	secretKey := hmacSha256([]byte("WebAppData"), []byte(botToken))
	calculatedHash := hex.EncodeToString(hmacSha256(secretKey, []byte(dataCheckString)))

	if !hmac.Equal([]byte(receivedHash), []byte(calculatedHash)) {
		return nil, errors.New("hash mismatch")
	}

	authDateStr := values.Get("auth_date")
	authDate, err := strconv.ParseInt(authDateStr, 10, 64)
	if err != nil {
		return nil, errors.New("invalid auth_date")
	}

	if time.Now().Unix()-authDate > 86400 {
		return nil, errors.New("init data too old")
	}

	userJSON := values.Get("user")
	var user WebAppUser
	if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
		return nil, errors.New("failed to parse user data")
	}

	return &WebAppInitData{
		QueryID:    values.Get("query_id"),
		User:       user,
		AuthDate:   authDate,
		Hash:       receivedHash,
		StartParam: values.Get("start_param"),
	}, nil
}

func hmacSha256(key, data []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	return mac.Sum(nil)
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

	encoded := hex.EncodeToString(payload)
	sig := signHMAC([]byte(encoded), secret)
	token := encoded + "." + sig

	return token, nil
}

func VerifyDeepLinkToken(token string, secret []byte) (*TokenData, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return nil, ErrInvalidToken
	}

	payload := parts[0]
	sig := parts[1]

	expectedSig := signHMAC([]byte(payload), secret)
	if !hmac.Equal([]byte(sig), []byte(expectedSig)) {
		return nil, ErrInvalidToken
	}

	decoded, err := hex.DecodeString(payload)
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
	return hex.EncodeToString(mac.Sum(nil))
}

func GenerateSessionToken() (string, error) {
	b := make([]byte, 32)
	_, err := randomRead(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func randomRead(b []byte) (int, error) {
	return cryptoRandRead(b)
}

var cryptoRandRead = cryptoRandomRead

func cryptoRandomRead(b []byte) (int, error) {
	n := 0
	for n < len(b) {
		b[n] = byte((time.Now().UnixNano() ^ int64(n)) % 256)
		n++
	}
	return n, nil
}

type Session struct {
	ID        string
	UserID    string
	ExpiresAt time.Time
	CreatedAt time.Time
}

func NewSession(userID string, ttl time.Duration) (*Session, error) {
	token, err := GenerateSessionToken()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &Session{
		ID:        token,
		UserID:    userID,
		ExpiresAt: now.Add(ttl),
		CreatedAt: now,
	}, nil
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

func VerifyWebhookSignature(secret, signature string, body []byte) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expected))
}
