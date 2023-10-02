package utils

import (
	`math/rand`
	`time`
)

const smallLetterBytes = "abcdefghijklmnopqrstuvwxyz"
const capitalLetterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const digitBytes = "1234567890"
const specialCharBytes = "`~!@#$%^&*()_+-={}[]|;:'\",<.>/?"

// Marker: =======================================
// Marker: Basic string functions
// Marker: =======================================

// GetRandomAlphaString will get a n character long random alphabetic string
func GetRandomAlphaString(n int) string {
	letterBytes := smallLetterBytes + capitalLetterBytes
	b := make([]byte, n)
	for i := range b {
		rand.Seed(time.Now().UTC().UnixNano())
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
