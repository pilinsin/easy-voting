package registrationutil

import (
	"EasyVoting/ipfs"
	"EasyVoting/util"
	"EasyVoting/util/crypto/encrypt"
)

type ManIdentity struct {
	rPriKey    *encrypt.PriKey
	hnmKeyFile *ipfs.KeyFile
}
func (mi ManIdentity) Private() *encrypt.PriKey { return mi.rPriKey }
func (mi ManIdentity) KeyFile() *ipfs.KeyFile { return mi.hnmKeyFile }

func (mi ManIdentity) Marshal() []byte {
	mManIdentity := &struct {
		RPriKey    []byte
		HnmKeyFile []byte
	}{mi.rPriKey.Marshal(), mi.hnmKeyFile.Marshal()}
	m, _ := util.Marshal(mManIdentity)
	return m
}
func (mi *ManIdentity) Unmarshal(m []byte) error {
	mManIdentity := &struct {
		RPriKey    []byte
		HnmKeyFile []byte
	}{}
	if err := util.Unmarshal(m, mManIdentity); err != nil {
		return err
	}

	priKey := &encrypt.PriKey{}
	if err := priKey.Unmarshal(mManIdentity.RPriKey); err != nil {
		return err
	}
	kf := &ipfs.KeyFile{}
	if err := kf.Unmarshal(mManIdentity.HnmKeyFile); err != nil {
		return err
	}
		
	mi.rPriKey = priKey
	mi.hnmKeyFile = kf
	return nil
}