package utils

import (
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func CreatePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err == nil {
		return hex.EncodeToString(hash), nil
	}

	return "", err
}

func IsHashMatchPassword(hexHash, password string) bool {
	hash, err := hex.DecodeString(hexHash) 
	if err != nil {
		return false
	}

	return bcrypt.CompareHashAndPassword(hash, []byte(password)) == nil
}