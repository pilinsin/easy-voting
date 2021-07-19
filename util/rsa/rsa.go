package rsa

import (
	"crypto/rand"
	crsa "crypto/rsa"
	"encoding/json"
	"golang.org/x/crypto/sha3"
	"math/big"

	"EasyVoting/util"
)

func bytes2bigInt(b []byte) *big.Int {
	return new(big.Int).SetBytes(b)
}

type MarshalableCRTValue struct {
	Exp   []byte
	Coeff []byte
	R     []byte
}

func ToMarshalableCRTValue(crt crsa.CRTValue) MarshalableCRTValue {
	exp := crt.Exp.Bytes()
	coeff := crt.Coeff.Bytes()
	r := crt.R.Bytes()

	return MarshalableCRTValue{exp, coeff, r}
}
func (mcrt MarshalableCRTValue) FromMarshalableCRTValue() crsa.CRTValue {
	return crsa.CRTValue{
		bytes2bigInt(mcrt.Exp),
		bytes2bigInt(mcrt.Coeff),
		bytes2bigInt(mcrt.R),
	}
}
func CRTValueEquals(crt1, crt2 crsa.CRTValue) bool {
	if ok := crt1.Exp.Cmp(crt2.Exp) == 0; !ok {
		return false
	}
	if ok := crt1.Coeff.Cmp(crt2.Coeff) == 0; !ok {
		return false
	}
	return crt1.R.Cmp(crt2.R) == 0
}

type MarshalablePrecomputed struct {
	Dp    []byte
	Dq    []byte
	Qinv  []byte
	MCRTs []MarshalableCRTValue
}

func ToMarshalablePrecomputed(pc crsa.PrecomputedValues) MarshalablePrecomputed {
	dp := pc.Dp.Bytes()
	dq := pc.Dq.Bytes()
	qinv := pc.Qinv.Bytes()
	mcrts := make([]MarshalableCRTValue, len(pc.CRTValues))
	for idx, crt := range pc.CRTValues {
		mcrts[idx] = ToMarshalableCRTValue(crt)
	}

	return MarshalablePrecomputed{dp, dq, qinv, mcrts}
}
func (mpc MarshalablePrecomputed) FromMarshalablePreComputed() crsa.PrecomputedValues {
	crts := make([]crsa.CRTValue, len(mpc.MCRTs))
	for idx, mcrt := range mpc.MCRTs {
		crts[idx] = mcrt.FromMarshalableCRTValue()
	}
	return crsa.PrecomputedValues{
		bytes2bigInt(mpc.Dp),
		bytes2bigInt(mpc.Dq),
		bytes2bigInt(mpc.Qinv),
		crts,
	}
}
func PrecomputedEquals(pc1, pc2 crsa.PrecomputedValues) bool {
	if ok := pc1.Dp.Cmp(pc2.Dp) == 0; !ok {
		return false
	}
	if ok := pc1.Dq.Cmp(pc2.Dq) == 0; !ok {
		return false
	}
	if ok := pc1.Qinv.Cmp(pc2.Qinv) == 0; !ok {
		return false
	}

	if len(pc1.CRTValues) != len(pc2.CRTValues) {
		return false
	}
	for idx := range pc1.CRTValues {
		if ok := CRTValueEquals(pc1.CRTValues[idx], pc2.CRTValues[idx]); !ok {
			return false
		}
	}

	return true
}

type PriKey struct {
	priKey crsa.PrivateKey
}
type MarshalablePriKey struct {
	MarshalablePubKey
	D      []byte
	Primes [][]byte
	MPC    MarshalablePrecomputed
}

func (pri PriKey) ToMarshalable() MarshalablePriKey {
	n := pri.priKey.PublicKey.N.Bytes()
	e := pri.priKey.PublicKey.E

	d := pri.priKey.D.Bytes()
	ps := make([][]byte, len(pri.priKey.Primes))
	for idx, v := range pri.priKey.Primes {
		ps[idx] = v.Bytes()
	}
	mpc := ToMarshalablePrecomputed(pri.priKey.Precomputed)

	return MarshalablePriKey{MarshalablePubKey{n, e}, d, ps, mpc}
}
func (mpri MarshalablePriKey) FromMarshalable() PriKey {
	ps := make([]*big.Int, len(mpri.Primes))
	for idx, p := range mpri.Primes {
		ps[idx] = bytes2bigInt(p)
	}
	prk := crsa.PrivateKey{
		crsa.PublicKey{bytes2bigInt(mpri.N), mpri.E},
		bytes2bigInt(mpri.D),
		ps,
		mpri.MPC.FromMarshalablePreComputed(),
	}
	return PriKey{prk}
}

type PubKey struct {
	pubKey crsa.PublicKey
}
type MarshalablePubKey struct {
	N []byte
	E int
}

func (pub PubKey) ToMarshalable() MarshalablePubKey {
	n := pub.pubKey.N.Bytes()
	e := pub.pubKey.E
	return MarshalablePubKey{n, e}
}
func (mpub MarshalablePubKey) FromMarshalable() PubKey {
	puk := crsa.PublicKey{
		bytes2bigInt(mpub.N),
		mpub.E,
	}
	return PubKey{puk}
}
func PublicKeyEquals(puk, puk2 crsa.PublicKey) bool {
	return puk.N.Cmp(puk2.N) == 0 && puk.E == puk2.E
}

type KeyPair struct {
	priKey PriKey
	pubKey PubKey
}

func (kp *KeyPair) Private() PriKey {
	return kp.priKey
}
func (kp *KeyPair) Public() PubKey {
	return kp.pubKey
}

func GenKeyPair(bit int) *KeyPair {
	priKey, err := crsa.GenerateKey(rand.Reader, bit)
	util.CheckError(err)
	pri := PriKey{*priKey}
	pub := PubKey{priKey.PublicKey}
	return &KeyPair{pri, pub}
}

func (key PubKey) Encrypt(message []byte) []byte {
	label := []byte("Encrypted")
	enc, err := crsa.EncryptOAEP(sha3.New512(), rand.Reader, &key.pubKey, message, label)
	util.CheckError(err)

	return enc
}

func (key PriKey) Decrypt(enc []byte) []byte {
	label := []byte("Encrypted")
	dec, err := crsa.DecryptOAEP(sha3.New512(), rand.Reader, &key.priKey, enc, label)
	util.CheckError(err)

	return dec
}

func (key PubKey) Equals(key2 PubKey) bool {
	return PublicKeyEquals(key.pubKey, key2.pubKey)
}
func (key PriKey) Equals(key2 PriKey) bool {
	if ok := PublicKeyEquals(key.priKey.PublicKey, key2.priKey.PublicKey); !ok {
		return false
	}
	if ok := key.priKey.D.Cmp(key2.priKey.D) == 0; !ok {
		return false
	}
	for idx, p1 := range key.priKey.Primes {
		if ok := p1.Cmp(key2.priKey.Primes[idx]) == 0; !ok {
			return false
		}
	}
	return PrecomputedEquals(key.priKey.Precomputed, key2.priKey.Precomputed)
}

func (key PriKey) Marshal() []byte {
	puk := key.ToMarshalable()
	mprk, err := json.Marshal(puk)
	util.CheckError(err)
	return mprk
}
func UnmarshalPrivate(mprk []byte) PriKey {
	var prk MarshalablePriKey
	err := json.Unmarshal(mprk, &prk)
	util.CheckError(err)

	return prk.FromMarshalable()
}

func (key PubKey) Marshal() []byte {
	puk := key.ToMarshalable()
	mpuk, err := json.Marshal(puk)
	util.CheckError(err)
	return mpuk
}
func UnmarshalPublic(mpuk []byte) PubKey {
	var puk MarshalablePubKey
	err := json.Unmarshal(mpuk, &puk)
	util.CheckError(err)

	return puk.FromMarshalable()
}
