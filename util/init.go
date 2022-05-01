package util

import (
	crdt "github.com/pilinsin/p2p-verse/crdt"
	"github.com/pilinsin/util/crypto"
)

func genKp() (crdt.IPrivKey, crdt.IPubKey, error) {
	kp := crypto.NewSignKeyPair()
	return kp.Sign(), kp.Verify(), nil
}
func marshalPub(pub crdt.IPubKey) ([]byte, error) {
	return crypto.MarshalVerfKey(pub.(crypto.IVerfKey))
}
func unmarshalPub(m []byte) (crdt.IPubKey, error) {
	return crypto.UnmarshalVerfKey(m)
}

func Init() {
	crdt.InitCryptoFuncs(genKp, marshalPub, unmarshalPub)
}
