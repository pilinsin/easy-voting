package registrationutil

import (
	"EasyVoting/ipfs"
	"EasyVoting/util"
	"EasyVoting/util/crypto"
)

type config struct {
	title          string
	rPubKey        crypto.IPubKey
	salt1          string
	salt2          string
	chmCid         string
	hnmIpnsName    string
	userDataLabels []string
}

func NewConfigs(title string, userDataset <-chan []string, userDataLabels []string, is *ipfs.IPFS) (*ManIdentity, *config) {
	encKeyPair := crypto.NewEncryptKeyPair()
	kf := ipfs.NewKeyFile()
	rCfg := newConfig(title, userDataset, userDataLabels, encKeyPair.Public(), is, kf)
	mi := &ManIdentity{
		rPriKey:    encKeyPair.Private(),
		hnmKeyFile: kf,
	}
	return mi, rCfg
}
func newConfig(title string, userDataset <-chan []string, userDataLabels []string, rPubKey crypto.IPubKey, is *ipfs.IPFS, kf *ipfs.KeyFile) *config {
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
func (cfg config) RPubKey() crypto.IPubKey   { return cfg.rPubKey }
func (cfg config) Salt1() string            { return cfg.salt1 }
func (cfg config) Salt2() string            { return cfg.salt2 }
func (cfg config) ChMapCid() string         { return cfg.chmCid }
func (cfg config) HnmIpnsName() string      { return cfg.hnmIpnsName }
func (cfg config) UserDataLabels() []string { return cfg.userDataLabels }

func (cfg config) IsCompatible(mi *ManIdentity) bool{
	pub := cfg.rPubKey.Equals(mi.rPriKey.Public())
	name, err := mi.hnmKeyFile.Name()
	nm := cfg.hnmIpnsName == name
	return pub && nm && (err == nil)
}

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

	pubKey, err := crypto.UnmarshalPubKey(mCfg.RPubKey)
	if err != nil {
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
