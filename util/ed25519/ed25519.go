package ed25519

import (
	"encoding/json"
	ced "golang.org/x/crypto/ed25519"

	"EasyVoting/util"
)

type SignKey struct {
	signKey ced.PrivateKey
}

type VerfKey struct {
	verfKey ced.PublicKey
}

type KeyPair struct {
	signKey SignKey
	verfKey VerfKey
}

func GenKeyPair() *KeyPair {
	pub, pri, err := ced.GenerateKey(nil)
	util.CheckError(err)

	kp := &KeyPair{SignKey{pri}, VerfKey{pub}}
	return kp
}
func (kp *KeyPair) Sign() SignKey {
	return kp.signKey
}
func (kp *KeyPair) Verf() VerfKey {
	return kp.verfKey
}

func (key SignKey) Sign(msg []byte) []byte {
	return ced.Sign(key.signKey, msg)
}
func (key VerfKey) Verify(msg, sig []byte) bool {
	return ced.Verify(key.verfKey, msg, sig)
}

func (key SignKey) Equals(key2 SignKey) bool {
	return util.ConstTimeBytesEqual(key.signKey, key2.signKey)
}
func (key VerfKey) Equals(key2 VerfKey) bool {
	return util.ConstTimeBytesEqual(key.verfKey, key2.verfKey)
}

func (key SignKey) Marshal() []byte {
	b, err := json.Marshal(key.signKey)
	util.CheckError(err)
	return b
}
func UnmarshalSign(b []byte) SignKey {
	var key ced.PrivateKey
	err := json.Unmarshal(b, &key)
	util.CheckError(err)

	return SignKey{key}
}

func (key VerfKey) Marshal() []byte {
	b, err := json.Marshal(key.verfKey)
	util.CheckError(err)
	return b
}
func UnmarshalVerf(b []byte) VerfKey {
	var key ced.PublicKey
	err := json.Unmarshal(b, &key)
	util.CheckError(err)

	return VerfKey{key}
}
