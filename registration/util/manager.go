package registrationutil

import (
	"github.com/pilinsin/easy-voting/ipfs"
	"github.com/pilinsin/easy-voting/util"
	"github.com/pilinsin/easy-voting/util/crypto"
)

type ManIdentity struct {
	rPriKey    crypto.IPriKey
	hnmKeyFile *ipfs.KeyFile
}
func (mi ManIdentity) Private() crypto.IPriKey { return mi.rPriKey }
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

	priKey, err := crypto.UnmarshalPriKey(mManIdentity.RPriKey)
	if err != nil {
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
