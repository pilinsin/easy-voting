package rsa

import (
	"crypto/rand"
	crsa "crypto/rsa"
	"encoding/base64"
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

func Encrypt(message string, key crsa.PublicKey) string {
	label := []byte("Encrypted")

	enc, err := crsa.EncryptOAEP(sha3.New512(), rand.Reader, &key, []byte(message), label)
	util.CheckError(err)

	return base64.StdEncoding.EncodeToString(enc)
}

func Decrypt(enc string, key crsa.PrivateKey) string {
	encBytes, err := base64.StdEncoding.DecodeString(enc)
	util.CheckError(err)
	label := []byte("Encrypted")
	decBytes, err := crsa.DecryptOAEP(sha3.New512(), rand.Reader, &key, encBytes, label)
	util.CheckError(err)

	message := string(decBytes)
	return message
}
