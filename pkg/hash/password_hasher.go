package hash

import (
	"crypto/sha1"
	"encoding/hex"
)

type SHA1PasswordHasher struct {
	salt string
}

func NewSHA1PasswordHasher(salt string) *SHA1PasswordHasher {
	return &SHA1PasswordHasher{salt: salt}
}

func (h *SHA1PasswordHasher) Hash(password string) (string, error) {
	hash := sha1.New()
	if _, err := hash.Write([]byte(password)); err != nil {
		return "", err
	}

	hashedPswdBytes := hash.Sum([]byte(h.salt))
	hashedPasswordHex := hex.EncodeToString(hashedPswdBytes)
	return hashedPasswordHex, nil
}
