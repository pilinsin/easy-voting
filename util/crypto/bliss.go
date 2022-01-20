package crypto

import (
	"github.com/LoCCS/bliss"
	"github.com/LoCCS/bliss/sampler"
	"github.com/pilinsin/easy-voting/util"
)

type blissKeyPair struct{
	signKey *blissSignKey
	verfKey *blissVerfKey
}
func newBlissKeyPair() ISignKeyPair{
	seedSize := int(sampler.SHA_512_DIGEST_LENGTH)
	seed := util.BytesToUint8s(util.GenRandomBytes(seedSize))
	e, _ := sampler.NewEntropy(seed)

	priKey, _ := bliss.GeneratePrivateKey(4, e)
	signKey := &blissSignKey{priKey, seed}
	pubKey := priKey.PublicKey()
	verfKey := &blissVerfKey{pubKey}

	return &blissKeyPair{signKey, verfKey}
}
func (kp *blissKeyPair) Sign() ISignKey{
	return kp.signKey
}
func (kp *blissKeyPair) Verify() IVerfKey{
	return kp.verfKey
}

type blissSignKey struct{
	signKey *bliss.PrivateKey
	seed []uint8
}
func (sk *blissSignKey) Close(){
	sk.signKey.Destroy()
}
func (sk *blissSignKey) Sign(data []byte) []byte{
	e, _ := sampler.NewEntropy(sk.seed)
	if signature, err := sk.signKey.Sign(data, e); err != nil{
		return nil
	}else{
		return signature.Encode()
	}
}
func (sk *blissSignKey) Verify() IVerfKey{
	verfKey := sk.signKey.PublicKey()
	return &blissVerfKey{verfKey}
}
func (sk *blissSignKey) Equals(sk2 ISignKey) bool{
	m := sk.Marshal()
	m2 := sk2.Marshal()
	return util.ConstTimeBytesEqual(m, m2)
}
func (sk *blissSignKey) Marshal() []byte{
	marshalSignKey := &struct{
		Sign []byte
		Seed []uint8
	}{sk.signKey.Encode(), sk.seed}
	m, _ := util.Marshal(marshalSignKey)
	return m
}
func (sk *blissSignKey) Unmarshal(m []byte) error{
	marshalSignKey := &struct{
		Sign []byte
		Seed []uint8
	}{}
	if err := util.Unmarshal(m, marshalSignKey); err != nil{
		return err
	}
	signKey, err := bliss.DecodePrivateKey(marshalSignKey.Sign)
	if err != nil{
		return err
	}

	sk.signKey = signKey
	sk.seed = marshalSignKey.Seed
	return nil
}

type blissVerfKey struct{
	verfKey *bliss.PublicKey
}
func (vk *blissVerfKey) Verify(data, sign []byte) bool{
	if signature, err := bliss.DecodeSignature(sign); err != nil{
		return false
	}else{
		ok, _ := vk.verfKey.Verify(data, signature)
		return ok
	}
}
func (vk *blissVerfKey) Equals(vk2 IVerfKey) bool{
	m := vk.Marshal()
	m2 := vk2.Marshal()
	return util.ConstTimeBytesEqual(m, m2)
}
func (vk *blissVerfKey) Marshal() []byte{
	return vk.verfKey.Encode()
}
func (vk *blissVerfKey) Unmarshal(m []byte) error{
	if verfKey, err := bliss.DecodePublicKey(m); err != nil{
		return nil
	}else{
		vk.verfKey = verfKey
		return nil
	}
}
