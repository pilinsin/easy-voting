package registrationutil

import (
	"github.com/pilinsin/ipfs-util"
	"github.com/pilinsin/util"
	"github.com/pilinsin/util/crypto"
)

type config struct {
	title          string
	rPubKey        crypto.IPubKey
	salt1          string
	salt2          string
	uhmCid         string
	hbmIpnsName    string
	userDataLabels []string
}

func NewConfigs(title string, userDataset <-chan []string, userDataLabels []string, is *ipfs.IPFS) (*ManIdentity, *config) {
	encKeyPair := crypto.NewPubEncryptKeyPair()
	kf := ipfs.Name.NewKeyFile()
	rCfg := newConfig(title, userDataset, userDataLabels, encKeyPair.Public(), is, kf)
	mi := &ManIdentity{
		rPriKey:    encKeyPair.Private(),
		hbmKeyFile: kf,
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

	uhHashes := make(chan UhHash)
	go func(){
		defer cloce(uhHashes)
		for userData := range userDataset {
			userHash := NewUserHash(cfg.salt1, userData...)
			uhHashes <- NewUhHash(cfg.salt2, userHash)
		}
	}()
	uhm := NewUhHashMap(uhHashes, 100000, is)
	cfg.uhmCid = ipfs.File.Add(uhm.Marshal(), is)

	hbm := NewHashBoxMap(100000)
	cfg.hbmIpnsName = ipfs.Name.PublishWithKeyFile(hbm.Marshal(), kf, is)

	return cfg
}
func (cfg config) Title() string            { return cfg.title }
func (cfg config) RPubKey() crypto.IPubKey   { return cfg.rPubKey }
func (cfg config) Salt1() string            { return cfg.salt1 }
func (cfg config) Salt2() string            { return cfg.salt2 }
func (cfg config) UhmCid() string         { return cfg.uhmCid }
func (cfg config) HbmIpnsName() string      { return cfg.hbmIpnsName }
func (cfg config) UserDataLabels() []string { return cfg.userDataLabels }

func (cfg config) IsCompatible(mi *ManIdentity) bool{
	txt := []byte("test pubKey message")
	enc, err  := cfg.rPubKey.Encrypt(txt)
	dec, err2 := mi.rPriKey.Decrypt(enc)
	pub := util.ConstTimeBytesEqual(txt, dec) && (err == nil) && (err2 == nil)
	name, err3 := mi.hbmKeyFile.Name()
	nm := (cfg.hbmIpnsName == name) && (err3 == nil)
	return pub && nm
}

func ConfigFromCid(rCfgCid string, is *ipfs.IPFS) (*config, error) {
	m, err := ipfs.File.Get(rCfgCid, is)
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
		UhmCid         string
		HbmIpnsName    string
		UserDataLabels []string
	}{cfg.title, cfg.rPubKey.Marshal(), cfg.salt1, cfg.salt2, cfg.uhmCid, cfg.hbmIpnsName, cfg.userDataLabels}
	m, _ := util.Marshal(mCfg)
	return m
}
func UnmarshalConfig(m []byte) (*config, error) {
	mCfg := &struct {
		Title          string
		RPubKey        []byte
		Salt1          string
		Salt2          string
		UhmCid         string
		HbmIpnsName    string
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
		uhmCid:         mCfg.UhmCid,
		hbmIpnsName:    mCfg.HbmIpnsName,
		userDataLabels: mCfg.UserDataLabels,
	}
	return cfg, nil
}
