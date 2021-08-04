package ed25519

import (
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
	signKey *SignKey
	verfKey *VerfKey
}

func GenKeyPair() *KeyPair {
	pub, pri, err := ced.GenerateKey(nil)
	util.CheckError(err)

	kp := &KeyPair{&SignKey{pri}, &VerfKey{pub}}
	return kp
}
func (kp *KeyPair) Sign() *SignKey {
	return kp.signKey
}
func (kp *KeyPair) Verf() *VerfKey {
	return kp.verfKey
}

func (key *SignKey) Sign(msg []byte) []byte {
	return ced.Sign(key.signKey, msg)
}
func (key *VerfKey) Verify(msg, sig []byte) bool {
	return ced.Verify(key.verfKey, msg, sig)
}

func (key *SignKey) Equals(key2 *SignKey) bool {
	return util.ConstTimeBytesEqual(key.Marshal(), key2.Marshal())
}
func (key *VerfKey) Equals(key2 *VerfKey) bool {
	return util.ConstTimeBytesEqual(key.Marshal(), key2.Marshal())
}

func (key *SignKey) Marshal() []byte {
	return key.signKey
}
func UnmarshalSign(b []byte) *SignKey {
	return &SignKey{b}
}

func (key *VerfKey) Marshal() []byte {
	return key.verfKey
}
func UnmarshalVerf(b []byte) *VerfKey {
	return &VerfKey{b}
}
