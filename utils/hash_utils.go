package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// GenerateHash generates a SHA-256 hash for any data that can be marshaled into JSON.
func GenerateHash(data interface{}) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(bytes)
	return hex.EncodeToString(hash[:]), nil
}
