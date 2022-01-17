package encrypt
/*
import (
	"crypto/rand"

	chacha "golang.org/x/crypto/chacha20poly1305"
	sntrup "github.com/companyzero/sntrup4591761"

	"EasyVoting/util"
)

const (
	pubKeySize = sntrup.PublicKeySize
	priKeySize = sntrup.PrivateKeySize
	cipherSize = sntrup.CiphertextSize
	sharedSize = sntrup.SharedKeySize
)

type KeyPair struct{
	pubKey *PubKey
	priKey *PriKey
}
func NewKeyPair() *KeyPair{
	pub, pri, _ := sntrup.GenerateKey(rand.Reader)
	pubKey := &PubKey{pub}
	priKey := &PriKey{pri}
	return &KeyPair{pubKey, priKey}
}
func (kp *KeyPair) Private() *PriKey{
	return kp.priKey
}
func (kp *KeyPair) Public() *PubKey{
	return kp.pubKey
}

type PriKey struct{
	priKey *[priKeySize]byte
}
func (pri *PriKey) Decrypt(m []byte) ([]byte, error){
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
func (pri *PriKey) Public() *PubKey{
	pub := new([pubKeySize]byte)
	copy(pub[:], pri.priKey[382:])
	return &PubKey{pub}
}
func (pri *PriKey) Equals(pri2 *PriKey) bool{
	return util.ConstTimeBytesEqual(pri.priKey[:], pri2.priKey[:])
}
func (pri *PriKey) Marshal() []byte{
	m, _ := util.Marshal(pri.priKey)
	return m
}
func (pri *PriKey) Unmarshal(m []byte) error{
	priKey := new([priKeySize]byte)
	if err := util.Unmarshal(m, priKey); err != nil{
		return err
	}
	pri.priKey = priKey
	return nil
}

type PubKey struct{
	pubKey *[pubKeySize]byte
}
func (pub *PubKey) Encrypt(data []byte) ([]byte, error){
	cipher, share, err := sntrup.Encapsulate(rand.Reader, pub.pubKey)
	if err != nil{return nil, err}
	
	aead, err := chacha.New(share[:])
	if err != nil{return nil, err}
	nonce := make([]byte, aead.NonceSize(), aead.NonceSize()+len(data)+aead.Overhead())
	if _, err := rand.Read(nonce); err != nil{return nil, err}
	enc := aead.Seal(nonce, nonce, data, nil)
	return append(cipher[:], enc...), nil
}
func (pub *PubKey) Equals(pub2 *PubKey) bool{
	return util.ConstTimeBytesEqual(pub.pubKey[:], pub2.pubKey[:])
}
func (pub *PubKey) Marshal() []byte{
	m, _ := util.Marshal(pub.pubKey)
	return m
}
func (pub *PubKey) Unmarshal(m []byte) error{
	pubKey := new([pubKeySize]byte)
	if err := util.Unmarshal(m, pubKey); err != nil{
		return err
	}
	pub.pubKey = pubKey
	return nil
}
*/