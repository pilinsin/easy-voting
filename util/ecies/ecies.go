package ecies

import (
	"encoding/json"

	"github.com/itrabbit/ecc"

	"EasyVoting/util"
)

type PubKey struct {
	pubKey ecc.PublicKey
}

type PriKey struct {
	priKey ecc.PrivateKey
}

type KeyPair struct {
	pubKey PubKey
	priKey PriKey
}

func GenKeyPair() *KeyPair {
	pri, err := ecc.GenerateKey()
	util.CheckError(err)

	kp := &KeyPair{PubKey{pri.PublicKey}, PriKey{*pri}}
	return kp
}
func (kp *KeyPair) Public() PubKey {
	return kp.pubKey
}
func (kp *KeyPair) Private() PriKey {
	return kp.priKey
}

func (key PubKey) Encrypt(message []byte) []byte {
	enc, err := ecc.Encrypt(&key.pubKey, message)
	util.CheckError(err)
	return enc
}
func (key PriKey) Decrypt(enc []byte) []byte {
	msg, err := ecc.Decrypt(&key.priKey, enc)
	util.CheckError(err)
	return msg
}

func (key PubKey) Equals(key2 PubKey) bool {
	return util.ConstTimeBytesEqual(key.Marshal(), key2.Marshal())
}
func (key PriKey) Equals(key2 PriKey) bool {
	return util.ConstTimeBytesEqual(key.Marshal(), key2.Marshal())
}

func (key PubKey) Marshal() []byte {
	b, err := json.Marshal(key.pubKey.String())
	util.CheckError(err)
	return b
}
func UnmarshalPublic(b []byte) PubKey {
	var pubstr string
	err := json.Unmarshal(b, &pubstr)
	util.CheckError(err)

	pub, err := ecc.PublicKeyFromString(pubstr)
	util.CheckError(err)
	return PubKey{*pub}
}

func (key PriKey) Marshal() []byte {
	b, err := json.Marshal(key.priKey.String())
	util.CheckError(err)
	return b
}
func UnmarshalPrivate(b []byte) PriKey {
	var pristr string
	err := json.Unmarshal(b, &pristr)
	util.CheckError(err)

	pri, err := ecc.KeyFromString(pristr)
	util.CheckError(err)
	return PriKey{*pri}
}
