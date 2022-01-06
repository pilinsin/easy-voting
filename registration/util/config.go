package registrationutil

import (
	"EasyVoting/ipfs"
	"EasyVoting/util"
	"EasyVoting/util/ecies"
)

type config struct {
	title          string
	rPubKey        *ecies.PubKey
	salt1          string
	salt2          string
	chmCid         string
	hnmIpnsName    string
	userDataLabels []string
}

func newConfig(title string, userDataset <-chan []string, userDataLabels []string, rPubKey *ecies.PubKey, is *ipfs.IPFS, kf *ipfs.KeyFile) *config {
	cfg := &config{
		title:          title,
		rPubKey:        rPubKey,
		salt1:          title + util.GenUniqueID(30, 30),
		salt2:          title + "_:_" + util.GenUniqueID(30, 30),
		userDataLabels: userDataLabels,
	}

	chm := NewConstHashMap([]UhHash{}, 100000, is)
	for userData := range userDataset {
		userHash := NewUserHash(is, cfg.salt1, userData...)
		chm.Append(NewUhHash(is, cfg.salt2, userHash), is)
	}
	cfg.chmCid = ipfs.ToCidWithAdd(chm.Marshal(), is)

	hnm := NewHashNameMap(100000)
	cfg.hnmIpnsName = ipfs.ToNameWithKeyFile(hnm.Marshal(), kf, is)

	return cfg
}
func (cfg config) Title() string            { return cfg.title }
func (cfg config) RPubKey() *ecies.PubKey   { return cfg.rPubKey }
func (cfg config) Salt1() string            { return cfg.salt1 }
func (cfg config) Salt2() string            { return cfg.salt2 }
func (cfg config) ChMapCid() string         { return cfg.chmCid }
func (cfg config) HnmIpnsName() string      { return cfg.hnmIpnsName }
func (cfg config) UserDataLabels() []string { return cfg.userDataLabels }

func ConfigFromCid(rCfgCid string, is *ipfs.IPFS) (*config, error) {
	m, err := ipfs.FromCid(rCfgCid, is)
	if err != nil {
		return nil, util.NewError("from rCfgCid error")
	}
	rCfg, err := UnmarshalConfig(m)
	if err != nil {
		return nil, util.NewError("unmarshal rCfgCid error")
	}
	return rCfg, nil
}
func (cfg config) Marshal() []byte {
	mCfg := &struct {
		Title          string
		RPubKey        []byte
		Salt1          string
		Salt2          string
		ChmCid         string
		HnmIpnsName    string
		UserDataLabels []string
	}{cfg.title, cfg.rPubKey.Marshal(), cfg.salt1, cfg.salt2, cfg.chmCid, cfg.hnmIpnsName, cfg.userDataLabels}
	m, _ := util.Marshal(mCfg)
	return m
}
func UnmarshalConfig(m []byte) (*config, error) {
	mCfg := &struct {
		Title          string
		RPubKey        []byte
		Salt1          string
		Salt2          string
		ChmCid         string
		HnmIpnsName    string
		UserDataLabels []string
	}{}
	err := util.Unmarshal(m, mCfg)
	if err != nil {
		return nil, err
	}

	pubKey := &ecies.PubKey{}
	if err := pubKey.Unmarshal(mCfg.RPubKey); err != nil {
		return nil, err
	}
	cfg := &config{
		title:          mCfg.Title,
		rPubKey:        pubKey,
		salt1:          mCfg.Salt1,
		salt2:          mCfg.Salt2,
		chmCid:         mCfg.ChmCid,
		hnmIpnsName:    mCfg.HnmIpnsName,
		userDataLabels: mCfg.UserDataLabels,
	}
	return cfg, nil
}
