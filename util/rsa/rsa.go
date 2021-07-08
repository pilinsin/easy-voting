package rsa

import (
	"crypto"
	"crypto/rand"
	crsa "crypto/rsa"
	"encoding/json"
	"golang.org/x/crypto/sha3"

	"EasyVoting/util"
)

type KeyPair struct {
	priKey crsa.PrivateKey
	pubKey crsa.PublicKey
}

func (kp *KeyPair) Private() crsa.PrivateKey {
	return kp.priKey
}
func (kp *KeyPair) Public() crsa.PublicKey {
	return kp.pubKey
}

func GenKeyPair(bit int) *KeyPair {
	priKey, err := crsa.GenerateKey(rand.Reader, bit)
	util.CheckError(err)
	return &KeyPair{*priKey, priKey.PublicKey}
}

func Encrypt(message []byte, key crsa.PublicKey) []byte {
	label := []byte("Encrypted")
	enc, err := crsa.EncryptOAEP(sha3.New512(), rand.Reader, &key, message, label)
	util.CheckError(err)

	return enc
}

func Decrypt(enc []byte, key crsa.PrivateKey) []byte {
	label := []byte("Encrypted")
	dec, err := crsa.DecryptOAEP(sha3.New512(), rand.Reader, &key, enc, label)
	util.CheckError(err)

	return dec
}

func Sign(message []byte, signKey crsa.PrivateKey) []byte {
	hash := sha3.Sum512(message)
	sign, err := crsa.SignPSS(rand.Reader, &signKey, crypto.SHA3_512, hash[:], nil)
	util.CheckError(err)

	return sign
}
func Verify(message []byte, sign []byte, verifyKey crsa.PublicKey) bool {
	hash := sha3.Sum512(message)
	err := crsa.VerifyPSS(&verifyKey, crypto.SHA3_512, hash[:], sign, nil)
	return err == nil
}

func MarshalPrivate(priKey crsa.PrivateKey) []byte {
	mprk, err := json.Marshal(priKey)
	util.CheckError(err)
	return mprk
}
func UnmarshalPrivate(mprk []byte) crsa.PrivateKey {
	var prk crsa.PrivateKey
	err := json.Unmarshal(mprk, &prk)
	util.CheckError(err)

	return prk
}

func MarshalPublic(pubKey crsa.PublicKey) []byte {
	mpuk, err := json.Marshal(pubKey)
	util.CheckError(err)
	return mpuk
}
func UnmarshalPublic(mpuk []byte) crsa.PublicKey {
	var puk crsa.PublicKey
	err := json.Unmarshal(mpuk, &puk)
	util.CheckError(err)

	return puk
}
