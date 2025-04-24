package server

import (
	"crypto/rand"
	"encoding/base64"
)

func generateState() string {
	b := make([]byte, 128)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	return state
}
