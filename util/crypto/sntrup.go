package crypto

import (
	"crypto/rand"

	chacha "golang.org/x/crypto/chacha20poly1305"
	sntrup "github.com/companyzero/sntrup4591761"

	"github.com/pilinsin/easy-voting/util"
)

const (
	pubKeySize = sntrup.PublicKeySize
	priKeySize = sntrup.PrivateKeySize
	cipherSize = sntrup.CiphertextSize
	sharedSize = sntrup.SharedKeySize
)

type sntrupKeyPair struct{
	pubKey *sntrupPubKey
	priKey *sntrupPriKey
}
func newSntrupKeyPair() IEncryptKeyPair{
	pub, pri, _ := sntrup.GenerateKey(rand.Reader)
	pubKey := &sntrupPubKey{pub}
	priKey := &sntrupPriKey{pri}
	return &sntrupKeyPair{pubKey, priKey}
}
func (kp *sntrupKeyPair) Private() IPriKey{
	return kp.priKey
}
func (kp *sntrupKeyPair) Public() IPubKey{
	return kp.pubKey
}

type sntrupPriKey struct{
	priKey *[priKeySize]byte
}
func (pri *sntrupPriKey) Decrypt(m []byte) ([]byte, error){
	if len(m) <= cipherSize{return nil, util.NewError("decrypt fail: len(m) <= cipherSize")}
	cipher := new([cipherSize]byte)
	copy(cipher[:], m[:cipherSize])

	share, flag := sntrup.Decapsulate(cipher, pri.priKey)
	if flag <= 0{return nil, util.NewError("decrypt fail: decapsulate error")}
	aead, err := chacha.New(share[:])
	if err != nil{return nil, err}
	if len(m) <= cipherSize + aead.NonceSize() + aead.Overhead(){
		return nil, util.NewError("decrypt fail: len(m) <= cipherSize+nonceSize+overHead")
	}
	encTotal := m[cipherSize:]
	nonce, enc := encTotal[:aead.NonceSize()], encTotal[aead.NonceSize():]
	if data, err := aead.Open(nil, nonce, enc, nil); err != nil{
		return nil, err
	}else{
		return data, nil
	}
}
func (pri *sntrupPriKey) Public() IPubKey{
	pub := new([pubKeySize]byte)
	copy(pub[:], pri.priKey[382:])
	return &sntrupPubKey{pub}
}
func (pri *sntrupPriKey) Equals(pri2 IPriKey) bool{
	return util.ConstTimeBytesEqual(pri.Marshal(), pri2.Marshal())
}
func (pri *sntrupPriKey) Marshal() []byte{
	m, _ := util.Marshal(pri.priKey)
	return m
}
func (pri *sntrupPriKey) Unmarshal(m []byte) error{
	priKey := new([priKeySize]byte)
	if err := util.Unmarshal(m, priKey); err != nil{
		return err
	}
	pri.priKey = priKey
	return nil
}

type sntrupPubKey struct{
	pubKey *[pubKeySize]byte
}
func (pub *sntrupPubKey) Encrypt(data []byte) ([]byte, error){
	cipher, share, err := sntrup.Encapsulate(rand.Reader, pub.pubKey)
	if err != nil{return nil, err}
	
	aead, err := chacha.New(share[:])
	if err != nil{return nil, err}
	nonce := make([]byte, aead.NonceSize(), aead.NonceSize()+len(data)+aead.Overhead())
	if _, err := rand.Read(nonce); err != nil{return nil, err}
	enc := aead.Seal(nonce, nonce, data, nil)
	return append(cipher[:], enc...), nil
}
func (pub *sntrupPubKey) Equals(pub2 IPubKey) bool{
	return util.ConstTimeBytesEqual(pub.Marshal(), pub2.Marshal())
}
func (pub *sntrupPubKey) Marshal() []byte{
	m, _ := util.Marshal(pub.pubKey)
	return m
}
func (pub *sntrupPubKey) Unmarshal(m []byte) error{
	pubKey := new([pubKeySize]byte)
	if err := util.Unmarshal(m, pubKey); err != nil{
		return err
	}
	pub.pubKey = pubKey
	return nil
}
