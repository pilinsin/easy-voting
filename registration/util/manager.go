package registrationutil

import (
	"github.com/pilinsin/ipfs-util"
	"github.com/pilinsin/util"
	"github.com/pilinsin/util/crypto"
)

type ManIdentity struct {
	rPriKey    crypto.IPriKey
	hbmKeyFile *ipfs.KeyFile
}
func (mi ManIdentity) Private() crypto.IPriKey { return mi.rPriKey }
func (mi ManIdentity) KeyFile() *ipfs.KeyFile { return mi.hbmKeyFile }

func (mi ManIdentity) Marshal() []byte {
	mManIdentity := &struct {
		RPriKey    []byte
		HbmKeyFile []byte
	}{mi.rPriKey.Marshal(), mi.hbmKeyFile.Marshal()}
	m, _ := util.Marshal(mManIdentity)
	return m
}
func (mi *ManIdentity) Unmarshal(m []byte) error {
	mManIdentity := &struct {
		RPriKey    []byte
		HbmKeyFile []byte
	}{}
	if err := util.Unmarshal(m, mManIdentity); err != nil {
		return err
	}

	priKey, err := crypto.UnmarshalPriKey(mManIdentity.RPriKey)
	if err != nil {
		return err
	}
	kf := &ipfs.KeyFile{}
	if err := kf.Unmarshal(mManIdentity.HbmKeyFile); err != nil {
		return err
	}
		
	mi.rPriKey = priKey
	mi.hbmKeyFile = kf
	return nil
}
