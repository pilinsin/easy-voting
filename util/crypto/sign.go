package crypto

import(
	"github.com/pilinsin/easy-voting/util"
)

type signMode int
const(
	Bliss signMode = iota
	Ed25519
)
var SelectedSignMode = Bliss


func NewSignKeyPair() ISignKeyPair{
	switch SelectedSignMode {
	case Bliss:
		return newBlissKeyPair()
	case Ed25519:
		return newEd25519KeyPair()
	default:
		return nil
	}
}

func UnmarshalSignKey(m []byte) (ISignKey, error){
	switch SelectedSignMode {
	case Bliss:
		sk := &blissSignKey{}
		err := sk.Unmarshal(m)
		return sk, err
	case Ed25519:
		sk := &ed25519SignKey{}
		err := sk.Unmarshal(m)
		return sk, err
	default:
		return nil, util.NewError("invalid SignMode is selected")
	}
}

func UnmarshalVerfKey(m []byte) (IVerfKey, error){
	switch SelectedSignMode {
	case Bliss:
		vk := &blissVerfKey{}
		err := vk.Unmarshal(m)
		return vk, err
	case Ed25519:
		vk := &ed25519VerfKey{}
		err := vk.Unmarshal(m)
		return vk, err
	default:
		return nil, util.NewError("invalid SignMode is selected")
	}
}