package ecies

import (
	eciesgo "github.com/ecies/go"

	"EasyVoting/util"
)

type PubKey struct {
	pubKey *eciesgo.PublicKey
}

type PriKey struct {
	priKey *eciesgo.PrivateKey
}

type KeyPair struct {
	pubKey *PubKey
	priKey *PriKey
}

func GenKeyPair() *KeyPair {
	pri, err := eciesgo.GenerateKey()
	util.CheckError(err)

	kp := &KeyPair{&PubKey{pri.PublicKey}, &PriKey{pri}}
	return kp
}
func (kp *KeyPair) Public() *PubKey {
	return kp.pubKey
}
func (kp *KeyPair) Private() *PriKey {
	return kp.priKey
}

func (key *PubKey) Encrypt(message []byte) []byte {
	enc, err := eciesgo.Encrypt(key.pubKey, message)
	util.CheckError(err)
	return enc
}
func (key *PriKey) Decrypt(enc []byte) []byte {
	msg, err := eciesgo.Decrypt(key.priKey, enc)
	util.CheckError(err)
	return msg
}

func (key *PubKey) Equals(key2 *PubKey) bool {
	return util.ConstTimeBytesEqual(key.Marshal(), key2.Marshal())
}
func (key *PriKey) Equals(key2 *PriKey) bool {
	return util.ConstTimeBytesEqual(key.Marshal(), key2.Marshal())
}

func (key *PubKey) Marshal() []byte {
	return key.pubKey.Bytes(true)
}
func UnmarshalPublic(b []byte) *PubKey {
	pub, err := eciesgo.NewPublicKeyFromBytes(b)
	util.CheckError(err)
	return &PubKey{pub}
}

func (key *PriKey) Marshal() []byte {
	return key.priKey.Bytes()
}
func UnmarshalPrivate(b []byte) *PriKey {
	pri := eciesgo.NewPrivateKeyFromBytes(b)
	return &PriKey{pri}
}
