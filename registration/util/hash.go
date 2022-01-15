package registrationutil

import (
	"EasyVoting/ipfs"
	"EasyVoting/util"
)

type StoreHash string
func NewStoreHash(kw string) StoreHash{
	salt := "identity_store_hash "
	if kw == ""{kw = "self"}
	kw = util.AnyBytes64ToStr(util.Hash([]byte(kw), []byte(salt)))
	return StoreHash(kw)
}

type UserHash string
func NewUserHash(is *ipfs.IPFS, salt string, userData ...string) UserHash {
	m, _ := util.Marshal(userData)
	cidStr := ipfs.ToCid(m, is)
	return UserHash(ipfs.ToCid(util.Hash([]byte(cidStr), []byte(salt)), is))
}

type UhHash string
func NewUhHash(is *ipfs.IPFS, salt string, userHash UserHash) UhHash {
	m := util.AnyStrToBytes64(string(userHash) + salt)
	cidStr := ipfs.ToCid(m, is)
	return UhHash(ipfs.ToCid(util.Hash([]byte(cidStr), []byte(salt)), is))
}
