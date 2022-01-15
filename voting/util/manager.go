package votingutil

import (
	"EasyVoting/ipfs"
	"EasyVoting/util"
	"EasyVoting/util/ecies"
)

type ManIdentity struct {
	manPriKey     *ecies.PriKey
	resMapKeyFile *ipfs.KeyFile
}

func (mi ManIdentity) Private() *ecies.PriKey { return mi.manPriKey }
func (mi ManIdentity) KeyFile() *ipfs.KeyFile { return mi.resMapKeyFile }

func (mi ManIdentity) Marshal() []byte {
	mManId := &struct {
		Pri, Kf []byte
	}{mi.manPriKey.Marshal(), mi.resMapKeyFile.Marshal()}
	m, _ := util.Marshal(mManId)
	return m
}
func (mi *ManIdentity) Unmarshal(m []byte) error {
	mManId := &struct{ Pri, Kf []byte }{}
	if err := util.Unmarshal(m, mManId); err != nil {
		return err
	}

	priKey := &ecies.PriKey{}
	if err := priKey.Unmarshal(mManId.Pri); err != nil {
		return err
	}
	kf := &ipfs.KeyFile{}
	if err := kf.Unmarshal(mManId.Kf); err != nil {
		return err
	}

	mi.manPriKey = priKey
	mi.resMapKeyFile = kf
	return nil
}
