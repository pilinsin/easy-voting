package ed25519

import (
	sed "crypto/ed25519"
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

func NewKeyPair() *KeyPair {
	pub, pri, _ := ced.GenerateKey(nil)
	kp := &KeyPair{&SignKey{pri}, &VerfKey{pub}}
	return kp
}
func (kp *KeyPair) Sign() *SignKey {
	return kp.signKey
}
func (kp *KeyPair) Verify() *VerfKey {
	return kp.verfKey
}

func (key *SignKey) Verify() *VerfKey {
	var priKey sed.PrivateKey = key.signKey
	pubKey := priKey.Public().(ced.PublicKey)
	return &VerfKey{pubKey}
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
func (key *SignKey) Unmarshal(b []byte) error {
	key.signKey = b
	return nil
}
func (key *VerfKey) Marshal() []byte {
	return key.verfKey
}
func (key *VerfKey) Unmarshal(b []byte) error {
	key.verfKey = b
	return nil
}
