package server

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
)

func generateState() string {
	b := make([]byte, 128)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	return state
}

func validateURL(values url.Values, secretKey string) bool {
	id := values.Get("id")
	username := values.Get("username")

	params := fmt.Sprintf("id=%s&username=%s", id, username)

	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(params))
	expectedMAC := mac.Sum(nil)

	signature := values.Get("signature")

	messageMAC, err := base64.URLEncoding.DecodeString(signature)
	if err != nil {
		return false
	}

	return hmac.Equal(messageMAC, expectedMAC)
}
