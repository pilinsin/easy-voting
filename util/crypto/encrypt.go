package crypto

import(
	"EasyVoting/util"
)

type encryptMode int
const(
	Sntrup encryptMode = iota
	Ecies
)
var SelectedEncMode = Sntrup


func NewEncryptKeyPair() IEncryptKeyPair{
	switch SelectedEncMode {
	case Sntrup:
		return newSntrupKeyPair()
	case Ecies:
		return newEciesKeyPair()
	default:
		return nil
	}
}

func UnmarshalPriKey(m []byte) (IPriKey, error){
	switch SelectedEncMode {
	case Sntrup:
		pri := &sntrupPriKey{}
		err := pri.Unmarshal(m)
		return pri, err
	case Ecies:
		pri := &eciesPriKey{}
		err := pri.Unmarshal(m)
		return pri, err
	default:
		return nil, util.NewError("invalid EncryptMode is selected")
	}
}

func UnmarshalPubKey(m []byte) (IPubKey, error){
	switch SelectedEncMode {
	case Sntrup:
		pub := &sntrupPubKey{}
		err := pub.Unmarshal(m)
		return pub, err
	case Ecies:
		pub := &eciesPubKey{}
		err := pub.Unmarshal(m)
		return pub, err
	default:
		return nil, util.NewError("invalid EncryptMode is selected")
	}
}