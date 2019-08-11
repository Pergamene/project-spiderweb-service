package guidgen

import (
	"errors"
	"math/rand"
)

// ErrMaxGUIDRetryAttempts is the err if the maximum number of attempts on generating a unique guid as been reached.
var ErrMaxGUIDRetryAttempts = errors.New("max number of retries (5) have been attempted for generating a new guid")

// MaxGUIDRetryAttempts represents a reasonable number of attempts to generate a unique guid.
const MaxGUIDRetryAttempts = 5

// GenerateGUID generates a pseudorandom guid with a set prefix and length, unrelated to the timestamp.
// Rather than having any guarantees of being unique, this function should be called several times over if
// a guid is found to not be unique, but this should be near unneccesary if lenght is sufficient.
func GenerateGUID(prefix string, length int) string {
	randomCharacters := length - len(prefix) - 1
	if randomCharacters <= 0 {
		return prefix
	}
	return prefix + "_" + getRandomString(randomCharacters)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

// SOURCE: https://medium.com/@kpbird/golang-generate-fixed-size-random-string-dd6dbd5e63c0
func getRandomString(length int) string {
	str := make([]rune, length)
	for i := range str {
		str[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(str)
}
