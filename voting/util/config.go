package votingutil

import (
	"github.com/pilinsin/ipfs-util"
	"github.com/pilinsin/util"
	"github.com/pilinsin/util/crypto"
	rutil "github.com/pilinsin/easy-voting/registration/util"
)

type VotingType int

const (
	Single VotingType = iota
	Block
	Approval
	Range
	Preference
	Cumulative
)

func VotingTypes() []string {
	return []string{
		"Approval",
		"Block",
		"Cumulative",
		"Preference",
		"Range",
		"Single",
	}
}
func (vt VotingType) VotingType2String() string {
	switch vt {
	case Single:
		return "Single"
	case Block:
		return "Block"
	case Approval:
		return "Approval"
	case Range:
		return "Range"
	case Preference:
		return "Preference"
	case Cumulative:
		return "Cumulative"
	default:
		return ""
	}
}
func VotingTypeFromStr(str string) (VotingType, error) {
	switch str {
	case "Single":
		return Single, nil
	case "Block":
		return Block, nil
	case "Approval":
		return Approval, nil
	case "Range":
		return Range, nil
	case "Preference":
		return Preference, nil
	case "Cumulative":
		return Cumulative, nil
	default:
		return -1, util.NewError("invalid VotingType")
	}
}

type VoteParams struct {
	Min   int
	Max   int
	Total int
}

type Candidate struct {
	Name  string
	Group string
	Url   string
	Image []byte
	ImageName string
}

type config struct {
	title          string
	votingID       string
	tInfo          *util.TimeInfo
	salt1          string
	salt2          string
	candidates     []Candidate
	manPubKey      crypto.IPubKey
	vParam         VoteParams
	vType          VotingType
	hbmCid string
	uhmCid string
	pimCid         string
	verfMapName string
	resBoxName     string
	userDataLabels []string
}
func NewConfigs(title, begin, end, loc, rCfgCid string, cands []Candidate, vParam VoteParams, vType VotingType, is *ipfs.IPFS) (*ManIdentity, *config, error) {
	encKeyPair := crypto.NewPubEncryptKeyPair()
	pub := encKeyPair.Public()
	pri := encKeyPair.Private()
	verfKf := ipfs.Name.NewKeyFile()
	resKf := ipfs.Name.NewKeyFile()
	resName, _ := resKf.Name()

	vCfg, err := newConfig(title, begin, end, loc, cands, pub, vParam, vType, rCfgCid, verfKf, resName, is)
	if err != nil {
		return nil, nil, err
	}
	mId := &ManIdentity{
		manPriKey:     pri,
		verfMapKeyFile: verfKf,
		resBoxKeyFile: resKf,
	}
	return mId, vCfg, nil
}
func newConfig(title, begin, end, loc string, cands []Candidate, manPubKey crypto.IPubKey, vParam VoteParams, vType VotingType, rCfgCid string, verfMapKeyFile *ipfs.KeyFile, resBoxName string, is *ipfs.IPFS) (*config, error) {
	rCfg, err := rutil.ConfigFromCid(rCfgCid, is)
	if err != nil {
		return nil, util.AddError(err, "invalid rCfgCid")
	}
	tInfo, ok := util.NewTimeInfo(begin, end, loc)
	if !ok {
		return nil, util.NewError("invalid time info")
	}

	cfg := &config{}
	cfg.title = title
	cfg.votingID = util.GenUniqueID(30, 30)
	cfg.tInfo = tInfo
	cfg.candidates = cands
	cfg.manPubKey = manPubKey
	cfg.vParam = vParam
	cfg.vType = vType
	cfg.salt1 = rCfg.Salt1()
	cfg.salt2 = rCfg.Salt2()
	cfg.hbmCid = ipfs.Name.GetCid(rCfg.HbmIpnsName(), is)
	cfg.resBoxName = resBoxName
	cfg.userDataLabels = rCfg.UserDataLabels()

	hbm, err = rutil.HashBoxMapFromName(rCfg.HbmIpnsName(), is)
	if err != nil {
		return nil, util.AddError(err, "hnm unmarshal error")
	}
	uhm := NewUvhHashMap(100000, is)
	pim := NewNameIdMap(100000, cfg.votingID)
	for kv := range hbm.NextKeyValue(is) {
		var uid string
		var uvhHash UvhHash
		for {
			uid = util.GenUniqueID(30, 6)
			uvHash := NewUidVidHash(uid, cfg.votingID)
			uvhHash = NewUvhHash(uvHash, cfg.votingID)
			if _, exist := uhm.ContainHash(uvhHash, is); !exist {
				break
			}
		}
		uhm.append(uvhHash, nil, is)
		pim.append(kv.Value().Public(), uid, is)
	}
	cfg.uhmCid = ipfs.File.Add(uhm.Marshal(), is)
	cfg.pimCid = ipfs.File.Add(pim.Marshal(), is)
	idVerfKeyMap := NewIdVerfKeyMap(100000)
	cfg.verfMapName = ipfs.Name.PublishWithKeyFile(idVerfKeyMap.Marshal(), verfMapKeyFile, is)
	return cfg, nil
}
func (cfg config) Title() string            { return cfg.title }
func (cfg config) VotingID() string         { return cfg.votingID }
func (cfg config) TimeInfo() *util.TimeInfo { return cfg.tInfo }
func (cfg config) Salt1() string            { return cfg.salt1 }
func (cfg config) Salt2() string            { return cfg.salt2 }
func (cfg config) Candidates() []Candidate  { return cfg.candidates }
func (cfg config) ManPubKey() crypto.IPubKey { return cfg.manPubKey }
func (cfg config) VParam() VoteParams       { return cfg.vParam }
func (cfg config) VType() VotingType        { return cfg.vType }
func (cfg config) HbmCid() string          { return cfg.hbmCid }
func (cfg config) UhmCid() string          { return cfg.uhmCid }
func (cfg config) PimCid() string          { return cfg.pimCid }
func (cfg config) VerfMapName() string { return cfg.verfMapName}
func (cfg config) ResBoxName() string       { return cfg.resBoxName }
func (cfg config) UserDataLabels() []string { return cfg.userDataLabels }

func (cfg config) IsCompatible(mi *ManIdentity) bool{
	txt := "test pubKey message"
	enc, err  := cfg.manPubKey.Encrypt([]byte(txt))
	dec, err2 := mi.manPriKey.Decrypt(enc)
	pub := (txt == string(dec)) && (err == nil) && (err2 == nil)
	verfName, vErr := mi.verfMapKeyFile.Name()
	vnm := cfg.verfMapName == verfName
	resName, rErr := mi.resMapKeyFile.Name()
	rnm := cfg.resBoxName == resName
	return pub && vnm && (vErr == nil) && rnm && (rErr == nil)
}

func ConfigFromCid(vCfgCid string, is *ipfs.IPFS) (*config, error) {
	m, err := ipfs.File.Get(vCfgCid, is)
	if err != nil {
		return nil, util.NewError("from vCfgCid error")
	}
	vCfg, err := UnmarshalConfig(m)
	if err != nil {
		return nil, util.NewError("unmarshal vCfgCid error")
	}
	return vCfg, nil
}
func (cfg config) Marshal() []byte {
	mCfg := &struct {
		Title          string
		VotingID       string
		TimeInfo       *util.TimeInfo
		Salt1          string
		Salt2          string
		Candidates     []Candidate
		ManPubKey      []byte
		VParam         VoteParams
		VType          VotingType
		HbmCid string
		UhmCid string
		PimCid         string
		VerfMapName string
		ResBoxName     string
		UserDataLabels []string
	}{
		Title:          cfg.title,
		VotingID:       cfg.votingID,
		TimeInfo:       cfg.tInfo,
		Salt1:          cfg.salt1,
		Salt2:          cfg.salt2,
		Candidates:     cfg.candidates,
		ManPubKey:      cfg.manPubKey.Marshal(),
		VParam:         cfg.vParam,
		VType:          cfg.vType,
		HbmCid: cfg.hbmCid,
		UhmCid: cfg.uhmCid,
		PimCid:         cfg.pimCid,
		VerfMapName: cfg.verfMapName,
		ResBoxName:     cfg.resBoxName,
		UserDataLabels: cfg.userDataLabels,
	}
	m, _ := util.Marshal(mCfg)
	return m
}
func UnmarshalConfig(m []byte) (*config, error) {
	mCfg := &struct {
		Title          string
		VotingID       string
		TimeInfo       *util.TimeInfo
		Salt1          string
		Salt2          string
		Candidates     []Candidate
		ManPubKey      []byte
		VParam         VoteParams
		VType          VotingType
		HbmCid string
		UhmCid string
		PimCid         string
		VerfMapName string
		ResBoxName     string
		UserDataLabels []string
	}{}
	err := util.Unmarshal(m, mCfg)
	if err != nil {
		return nil, err
	}

	pub, err := crypto.UnmarshalPubKey(mCfg.ManPubKey)
	if err != nil {
		return nil, err
	}
	cfg := &config{
		title:          mCfg.Title,
		votingID:       mCfg.VotingID,
		tInfo:          mCfg.TimeInfo,
		salt1:          mCfg.Salt1,
		salt2:          mCfg.Salt2,
		candidates:     mCfg.Candidates,
		manPubKey:      pub,
		vParam:         mCfg.VParam,
		vType:          mCfg.VType,
		hbmCid: mCfg.HbmCid,
		uhmCid: mCfg.UhmCid,
		pimCid:         mCfg.PimCid,
		verfMapName: mCfg.VerfMapName,
		resBoxName:     mCfg.ResBoxName,
		userDataLabels: mCfg.UserDataLabels,
	}
	return cfg, nil
}
func (cfg *config) ShuffleCandidates() {
	n := len(cfg.candidates)
	for i := 0; i < n-1; i++ {
		j := i + util.RandInt(n-i) //[i, n)
		cfg.candidates[i], cfg.candidates[j] = cfg.candidates[j], cfg.candidates[i]
	}
}
func (cfg *config) CandNameGroups() []string {
	ngs := make([]string, len(cfg.candidates))
	for idx, candidate := range cfg.candidates {
		ngs[idx] = candidate.Name + ", " + candidate.Group + " _" + cfg.votingID
	}
	return ngs
}
