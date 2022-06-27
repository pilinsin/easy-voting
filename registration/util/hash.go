package registrationutil

import (
	"encoding/base64"
	"strings"

	hash "github.com/pilinsin/util/hash"
)

func NewUserHash(salt string, userData ...string) string {
	s := strings.Join(userData, " ")
	s = strings.Join(strings.Fields(s), "#@%@#")
	m := hash.Hash([]byte(s), []byte(salt))
	m = hash.HashWithSize([]byte(salt), m, 64)
	return base64.URLEncoding.EncodeToString(m)
}
