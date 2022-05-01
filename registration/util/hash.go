package registrationutil

import (
	"encoding/base64"
	"github.com/pilinsin/util/crypto"
	"strings"
)

func NewUserHash(salt string, userData ...string) string {
	s := strings.Join(userData, " ")
	s = strings.Join(strings.Fields(s), "#@%@#")
	m := crypto.Hash([]byte(s), []byte(salt))
	m = crypto.HashWithSize([]byte(salt), m, 64)
	return base64.URLEncoding.EncodeToString(m)
}
