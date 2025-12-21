package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// GenerateKey creates a hash from any object for use as directory name
func GenerateKey(obj interface{}) (string, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}