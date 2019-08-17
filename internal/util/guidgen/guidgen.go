package guidgen

import (
	"math/rand"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/pkg/errors"
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

// CheckProposedGUID validates that a proposedGUID conforms to the standard pattern of "<prefix>_<alphanmeric string>" of length.
func CheckProposedGUID(proposedGUID, prefix string, length int) error {
	if proposedGUID != "" && utf8.RuneCountInString(proposedGUID) != length {
		return errors.Errorf("proposed guid must be %v characters", length)
	}
	if prefix == "" {
		return errors.New("prefix must be at least one character")
	}
	if !isAlphaNumeric(prefix) {
		return errors.Errorf("prefix '%v' is not alphanumeric. Do not include '_' in the prefix. It will be added automatically", prefix)
	}
	if length < 3 {
		return errors.Errorf("length of %v is invalid. Must be at least 3 for a single-character prefix + '_' + a random character for the guid", length)
	}
	// prefix must at least allow for one randomized character after it plus "_"
	if utf8.RuneCountInString(prefix) > length-2 {
		return errors.Errorf("proposed prefix must be less than %v characters", length-2)
	}
	if proposedGUID != "" && !strings.HasPrefix(proposedGUID, prefix+"_") {
		return errors.Errorf("proposed guid must start with '%v'", prefix+"_")
	}
	if proposedGUID != "" && !isAlphaNumeric(proposedGUID[utf8.RuneCountInString(prefix)+1:]) {
		return errors.Errorf("characters after prefix, '%v', must be English alphanumeric", proposedGUID[:len(prefix)+1])
	}
	return nil
}

var isAlphaNumericRegex = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString

func isAlphaNumeric(str string) bool {
	return isAlphaNumericRegex(str)
}
