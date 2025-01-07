package util

import (
	"math/rand"
	"time"
)

func RandomCode() string {
	const charset = "0123456789abcdefghijklmnopqrlstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, 6)
	rand.Seed(time.Now().UnixNano())
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
