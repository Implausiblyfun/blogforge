// Package encrypt shall be a super silly fake encoding strategy but allows for future slotting in of better patterns.-
package encrypt

import (
	"encoding/base64"
	"errors"
	"fmt"
)

var scaryKey string

// SetKey to the secret, this is a silly way to denote a singular encryption key but...
// This is a quick and dirty app anyways for now.
func SetKey(secret string) (err error) {
	if len(scaryKey) > 0 {
		err = errors.New("Overriding the existing key which was not empty. You may break everything.")
	}

	scaryKey = secret
	return
}

const sillyFormat = "%s::%s"

// Encode the given string with our secrets.
func Encode(toEncode string) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(sillyFormat, scaryKey, toEncode)))
}

// Decode the given string with our secrets.
func Decode(toDecode string) (string, error) {
	baseBytes, err := base64.StdEncoding.DecodeString(toDecode)
	baseStr := string(baseBytes)
	if err != nil {
		return string(baseStr), err
	}

	if len(baseStr) < len(scaryKey)+2 {
		return string(baseStr), errors.New("incorrectly formatted source")
	}
	return baseStr[len(scaryKey)+2:], nil
}
