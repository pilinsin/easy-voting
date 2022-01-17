package sign
/*
import (
	"github.com/LoCCS/bliss"
	"github.com/LoCCS/bliss/sampler"
	"EasyVoting/util"
)

type KeyPair struct{
	signKey *SignKey
	verfKey *VerfKey
}
func NewKeyPair() *KeyPair{
	seedSize := int(sampler.SHA_512_DIGEST_LENGTH)
	seed := util.BytesToUint8s(util.GenRandomBytes(seedSize))
	e, _ := sampler.NewEntropy(seed)

	priKey, _ := bliss.GeneratePrivateKey(4, e)
	signKey := &SignKey{priKey, seed}
	pubKey := priKey.PublicKey()
	verfKey := &VerfKey{pubKey}

	return &KeyPair{signKey, verfKey}
}
func (kp *KeyPair) Sign() *SignKey{
	return kp.signKey
}
func (kp *KeyPair) Verify() *VerfKey{
	return kp.verfKey
}

type SignKey struct{
	signKey *bliss.PrivateKey
	seed []uint8
}
func (sk *SignKey) Close(){
	sk.signKey.Destroy()
}
func (sk *SignKey) Sign(data []byte) []byte{
	e, _ := sampler.NewEntropy(sk.seed)
	if signature, err := sk.signKey.Sign(data, e); err != nil{
		return nil
	}else{
		return signature.Encode()
	}
}
func (sk *SignKey) Verify() *VerfKey{
	verfKey := sk.signKey.PublicKey()
	return &VerfKey{verfKey}
}
func (sk *SignKey) Equals(sk2 *SignKey) bool{
	m := sk.Marshal()
	m2 := sk2.Marshal()
	return util.ConstTimeBytesEqual(m, m2)
}
func (sk *SignKey) Marshal() []byte{
	marshalSignKey := &struct{
		Sign []byte
		Seed []uint8
	}{sk.signKey.Encode(), sk.seed}
	m, _ := util.Marshal(marshalSignKey)
	return m
}
func (sk *SignKey) Unmarshal(m []byte) error{
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

type VerfKey struct{
	verfKey *bliss.PublicKey
}
func (vk *VerfKey) Verify(data, sign []byte) bool{
	if signature, err := bliss.DecodeSignature(sign); err != nil{
		return false
	}else{
		ok, _ := vk.verfKey.Verify(data, signature)
		return ok
	}
}
func (vk *VerfKey) Equals(vk2 *VerfKey) bool{
	m := vk.Marshal()
	m2 := vk2.Marshal()
	return util.ConstTimeBytesEqual(m, m2)
}
func (vk *VerfKey) Marshal() []byte{
	return vk.verfKey.Encode()
}
func (vk *VerfKey) Unmarshal(m []byte) error{
	if verfKey, err := bliss.DecodePublicKey(m); err != nil{
		return nil
	}else{
		vk.verfKey = verfKey
		return nil
	}
}
*/