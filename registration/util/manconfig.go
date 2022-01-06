package registrationutil

import (
	"EasyVoting/ipfs"
	"EasyVoting/util"
	"EasyVoting/util/ecies"
)

type manConfig struct {
	*config
	rPriKey    *ecies.PriKey
	hnmKeyFile *ipfs.KeyFile
}

func NewConfigs(title string, userDataset <-chan []string, userDataLabels []string, is *ipfs.IPFS) (*manConfig, *config) {
	encKeyPair := ecies.NewKeyPair()
	kf := ipfs.NewKeyFile()
	rCfg := newConfig(title, userDataset, userDataLabels, encKeyPair.Public(), is, kf)
	mCfg := &manConfig{
		config:     rCfg,
		rPriKey:    encKeyPair.Private(),
		hnmKeyFile: kf,
	}
	return mCfg, rCfg
}
func (mCfg manConfig) Private() *ecies.PriKey { return mCfg.rPriKey }
func (mCfg manConfig) KeyFile() *ipfs.KeyFile { return mCfg.hnmKeyFile }
func (mCfg manConfig) Config() *config        { return mCfg.config }

func ManConfigFromCid(mCfgCid string, is *ipfs.IPFS) (*manConfig, error) {
	m, err := ipfs.FromCid(mCfgCid, is)
	if err != nil {
		return nil, util.NewError("from mCfgCid error")
	}
	mCfg, err := UnmarshalManConfig(m)
	if err != nil {
		return nil, util.NewError("unmarshal mCfgCid error")
	}
	return mCfg, nil
}
func (mCfg manConfig) Marshal() []byte {
	mManCfg := &struct {
		MRCfg      []byte
		RPriKey    []byte
		HnmKeyFile []byte
	}{mCfg.config.Marshal(), mCfg.rPriKey.Marshal(), mCfg.hnmKeyFile.Marshal()}
	m, _ := util.Marshal(mManCfg)
	return m
}
func UnmarshalManConfig(m []byte) (*manConfig, error) {
	mManCfg := &struct {
		MRCfg      []byte
		RPriKey    []byte
		HnmKeyFile []byte
	}{}
	if err := util.Unmarshal(m, mManCfg); err != nil {
		return nil, err
	}

	rCfg, err := UnmarshalConfig(mManCfg.MRCfg)
	if err != nil {
		return nil, err
	}
	priKey := &ecies.PriKey{}
	if err := priKey.Unmarshal(mManCfg.RPriKey); err != nil {
		return nil, err
	}
	kf := &ipfs.KeyFile{}
	if err := kf.Unmarshal(mManCfg.HnmKeyFile); err != nil {
		return nil, err
	}
	mCfg := &manConfig{
		config:     rCfg,
		rPriKey:    priKey,
		hnmKeyFile: kf,
	}
	return mCfg, nil
}
