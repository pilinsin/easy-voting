package registrationutil

import (
	"github.com/pilinsin/util"
	"github.com/pilinsin/util/crypto"
)


func NewUserHash(salt string, userData ...string) string {
	m, _ := util.Marshal(userData)
	m = crypto.Hash([]byte(m), []byte(salt))
	m = crypto.HashWithSize([]byte(salt), m, 64)
	return util.AnyBytes64ToStr(m)
}
