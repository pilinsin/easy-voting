package votingutil

import (
	"github.com/pilinsin/easy-voting/ipfs"
	"github.com/pilinsin/easy-voting/util"
	"github.com/pilinsin/easy-voting/util/crypto"
)

type ManIdentity struct {
	manPriKey     crypto.IPriKey
	verfMapKeyFile *ipfs.KeyFile
	resMapKeyFile *ipfs.KeyFile
}

func (mi ManIdentity) Private() crypto.IPriKey { return mi.manPriKey }
func (mi ManIdentity) VerfMapKeyFile() *ipfs.KeyFile { return mi.verfMapKeyFile }
func (mi ManIdentity) ResMapKeyFile() *ipfs.KeyFile { return mi.resMapKeyFile }

func (mi ManIdentity) Marshal() []byte {
	mManId := &struct {
		Pri, VerfKf, ResKf []byte
	}{mi.manPriKey.Marshal(), mi.verfMapKeyFile.Marshal(), mi.resMapKeyFile.Marshal()}
	m, _ := util.Marshal(mManId)
	return m
}
func (mi *ManIdentity) Unmarshal(m []byte) error {
	mManId := &struct{ Pri, VerfKf, ResKf []byte }{}
	if err := util.Unmarshal(m, mManId); err != nil {
		return err
	}

	priKey, err := crypto.UnmarshalPriKey(mManId.Pri)
	if err != nil {
		return err
	}
	verfKf := &ipfs.KeyFile{}
	if err := verfKf.Unmarshal(mManId.VerfKf); err != nil {
		return err
	}
	resKf := &ipfs.KeyFile{}
	if err := resKf.Unmarshal(mManId.ResKf); err != nil {
		return err
	}

	mi.manPriKey = priKey
	mi.verfMapKeyFile = verfKf
	mi.resMapKeyFile = resKf
	return nil
}
