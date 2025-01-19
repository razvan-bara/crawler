package queue

import (
	"crypto/sha256"
	"encoding/hex"
)

var deduplicationMap = make(map[string]bool)

func isUniqueMessage(textBytes []byte) bool {
	hash := sha256.Sum256(textBytes)
	hashStr := hex.EncodeToString(hash[:])
	_, ok := deduplicationMap[hashStr]
	if ok {
		return false
	}

	deduplicationMap[hashStr] = true
	return true
}
