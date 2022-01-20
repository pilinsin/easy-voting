package crypto

import (
	eciesgo "github.com/ecies/go"

	"github.com/pilinsin/easy-voting/util"
)

type eciesPubKey struct {
	pubKey *eciesgo.PublicKey
}

type eciesPriKey struct {
	priKey *eciesgo.PrivateKey
}

type eciesKeyPair struct {
	pubKey *eciesPubKey
	priKey *eciesPriKey
}

func newEciesKeyPair() IEncryptKeyPair {
	pri, _ := eciesgo.GenerateKey()
	kp := &eciesKeyPair{&eciesPubKey{pri.PublicKey}, &eciesPriKey{pri}}
	return kp
}
func (kp *eciesKeyPair) Public() IPubKey {
	return kp.pubKey
}
func (kp *eciesKeyPair) Private() IPriKey {
	return kp.priKey
}

func (key *eciesPriKey) Public() IPubKey {
	return &eciesPubKey{key.priKey.PublicKey}
}

func (key *eciesPubKey) Encrypt(message []byte) ([]byte, error) {
	return eciesgo.Encrypt(key.pubKey, message)
}
func (key *eciesPriKey) Decrypt(enc []byte) ([]byte, error) {
	return eciesgo.Decrypt(key.priKey, enc)
}

func (key *eciesPubKey) Equals(key2 IPubKey) bool {
	return util.ConstTimeBytesEqual(key.Marshal(), key2.Marshal())
}
func (key *eciesPriKey) Equals(key2 IPriKey) bool {
	return util.ConstTimeBytesEqual(key.Marshal(), key2.Marshal())
}

func (key *eciesPubKey) Marshal() []byte {
	return key.pubKey.Bytes(true)
}
func (key *eciesPubKey) Unmarshal(b []byte) error {
	pub, err := eciesgo.NewPublicKeyFromBytes(b)
	if err == nil {
		key.pubKey = pub
		return nil
	} else {
		return err
	}
}

func (key *eciesPriKey) Marshal() []byte {
	return key.priKey.Bytes()
}
func (key *eciesPriKey) Unmarshal(b []byte) error {
	pri := eciesgo.NewPrivateKeyFromBytes(b)
	key.priKey = pri
	return nil
}
