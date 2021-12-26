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

func NewKeyPair() *KeyPair {
	pri, _ := eciesgo.GenerateKey()
	kp := &KeyPair{&PubKey{pri.PublicKey}, &PriKey{pri}}
	return kp
}
func (kp *KeyPair) Public() *PubKey {
	return kp.pubKey
}
func (kp *KeyPair) Private() *PriKey {
	return kp.priKey
}

func (key *PriKey) Public() *PubKey {
	return &PubKey{key.priKey.PublicKey}
}

func (key *PubKey) Encrypt(message []byte) ([]byte, error) {
	return eciesgo.Encrypt(key.pubKey, message)
}
func (key *PriKey) Decrypt(enc []byte) ([]byte, error) {
	return eciesgo.Decrypt(key.priKey, enc)
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
func (key *PubKey) Unmarshal(b []byte) error {
	pub, err := eciesgo.NewPublicKeyFromBytes(b)
	if err == nil {
		key.pubKey = pub
		return nil
	} else {
		return err
	}
}

func (key *PriKey) Marshal() []byte {
	return key.priKey.Bytes()
}
func (key *PriKey) Unmarshal(b []byte) error {
	pri := eciesgo.NewPrivateKeyFromBytes(b)
	key.priKey = pri
	return nil
}
