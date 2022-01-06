package votingutil

import (
	"EasyVoting/ipfs"
	"EasyVoting/util"
	"EasyVoting/util/ecies"
)

type manConfig struct {
	*config
	manPriKey     *ecies.PriKey
	resMapKeyFile *ipfs.KeyFile
}

func NewConfigs(title, begin, end, loc, rCfgCid string, cands []Candidate, vParam VoteParams, vType VotingType, is *ipfs.IPFS) (*manConfig, *config, error) {
	encKeyPair := ecies.NewKeyPair()
	pub := encKeyPair.Public()
	pri := encKeyPair.Private()
	kf := ipfs.NewKeyFile()
	name, _ := kf.Name()

	vCfg, err := newConfig(title, begin, end, loc, cands, pub, vParam, vType, rCfgCid, name, is)
	if err != nil {
		return nil, nil, err
	}
	mCfg := &manConfig{
		config:        vCfg,
		manPriKey:     pri,
		resMapKeyFile: kf,
	}
	return mCfg, vCfg, nil
}
func (mCfg manConfig) Private() *ecies.PriKey { return mCfg.manPriKey }
func (mCfg manConfig) KeyFile() *ipfs.KeyFile { return mCfg.resMapKeyFile }
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
	mMCfg := &struct {
		Cfg, Pri, Kf []byte
	}{mCfg.config.Marshal(), mCfg.manPriKey.Marshal(), mCfg.resMapKeyFile.Marshal()}
	m, _ := util.Marshal(mMCfg)
	return m
}
func UnmarshalManConfig(m []byte) (*manConfig, error) {
	mCfg := &struct{ Cfg, Pri, Kf []byte }{}
	err := util.Unmarshal(m, mCfg)
	if err != nil {
		return nil, err
	}

	cfg, err := UnmarshalConfig(mCfg.Cfg)
	if err != nil {
		return nil, err
	}
	priKey := &ecies.PriKey{}
	if err := priKey.Unmarshal(mCfg.Pri); err != nil {
		return nil, err
	}
	kf := &ipfs.KeyFile{}
	if err := kf.Unmarshal(mCfg.Kf); err != nil {
		return nil, err
	}

	return &manConfig{cfg, priKey, kf}, nil
}
