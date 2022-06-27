package util

import (
	"errors"

	crdt "github.com/pilinsin/p2p-verse/crdt"
	ipub "github.com/pilinsin/util/public"
	ecies "github.com/pilinsin/util/public/ecies"
	isign "github.com/pilinsin/util/sign"
	ed25519 "github.com/pilinsin/util/sign/ed25519"
)

type ISignKey = isign.ISignKey
type IVerfKey = isign.IVerfKey

var NewSignKeyPair = ed25519.NewKeyPair
var UnmarshalSign = ed25519.UnmarshalSignKey
var UnmarshalVerf = ed25519.UnmarshalVerfKey

type IPriKey = ipub.IPriKey
type IPubKey = ipub.IPubKey

var NewPubKeyPair = ecies.NewKeyPair
var UnmarshalPri = ecies.UnmarshalPriKey
var UnmarshalPub = ecies.UnmarshalPubKey

func genKp() (crdt.IPrivKey, crdt.IPubKey, error) {
	kp := NewSignKeyPair()
	return kp.Sign(), kp.Verify(), nil
}
func marshalPub(pub crdt.IPubKey) ([]byte, error) {
	if pub == nil {
		return nil, errors.New("pub is nil")
	}
	return pub.Raw()
}
func unmarshalPub(m []byte) (crdt.IPubKey, error) {
	return UnmarshalVerf(m)
}

func Init() {
	crdt.InitCryptoFuncs(genKp, marshalPub, unmarshalPub)
}
