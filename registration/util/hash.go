package registrationutil

import (
	"github.com/pilinsin/ipfs-util"
	"github.com/pilinsin/util"
	"github.com/pilinsin/util/crypto"
)

type StoreHash string
func NewStoreHash(kw string) StoreHash{
	salt := "identity_store_hash "
	if kw == ""{kw = "self"}
	kw = util.AnyBytes64ToStr(crypto.Hash([]byte(kw), []byte(salt)))
	return StoreHash(kw)
}

type UserHash string
func NewUserHash(salt string, userData ...string) UserHash {
	m, _ := util.Marshal(userData)
	m = crypto.Hash([]byte(m), []byte(salt))
	m = crypto.HashWithSize([]byte(salt), m, 64)
	return UserHash(util.AnyBytes64ToStr(m))
}

type UhHash string
func NewUhHash(salt string, userHash UserHash) UhHash {
	m := util.AnyStrToBytes64(string(userHash) + salt)
	m = crypto.Hash([]byte(m), []byte(salt))
	m = crypto.HashWithSize([]byte(salt), m, 64)
	return UhHash(util.AnyBytes64ToStr(m))
}
