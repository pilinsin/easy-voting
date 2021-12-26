package votingutil

import (
	"fyne.io/fyne/v2"

	"EasyVoting/ipfs"
	rutil "EasyVoting/registration/util"
	"EasyVoting/util"
	"EasyVoting/util/ecies"
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
	Image fyne.Resource
}

type config struct {
	title          string
	votingID       string
	tInfo          *util.TimeInfo
	salt1          string
	salt2          string
	candidates     []Candidate
	manPubKey      *ecies.PubKey
	vParam         VoteParams
	vType          VotingType
	chmCid         string
	nimCid         string
	ivmCid         string
	resMapName     string
	userDataLabels []string
}

func newConfig(title, begin, end, loc string, cands []Candidate, manPubKey *ecies.PubKey, vParam VoteParams, vType VotingType, rCfgCid string, resMapName string, is *ipfs.IPFS) (*config, error) {
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
	cfg.resMapName = resMapName
	cfg.userDataLabels = rCfg.UserDataLabels()

	hnm := &rutil.HashNameMap{}
	err = hnm.FromName(rCfg.HnmIpnsName(), is)
	if err != nil {
		return nil, util.AddError(err, "hnm unmarshal error")
	}
	chm := rutil.NewConstHashMap([]rutil.UhHash{}, 100000, is)
	nim := NewNameIdMap(100000, cfg.votingID)
	ivm := NewIdVotingMap(100000, tInfo)
	for kv := range hnm.NextKeyValue(is) {
		uhHash := kv.Key()
		chm.Append(uhHash, is)

		var uid string
		var uvHash UidVidHash
		for {
			uid = util.GenUniqueID(30, 6)
			uvHash = NewUidVidHash(uid, cfg.votingID)
			if _, _, ok := ivm.ContainHash(uvHash, is); !ok {
				break
			}
		}
		rIpnsName := kv.Value().Name()
		nim.Append(rIpnsName, uid, is)
		ivm.Append(uvHash, rIpnsName, is)
	}
	cfg.chmCid = ipfs.ToCidWithAdd(chm.Marshal(), is)
	cfg.nimCid = ipfs.ToCidWithAdd(nim.Marshal(), is)
	cfg.ivmCid = ipfs.ToCidWithAdd(ivm.Marshal(), is)

	return cfg, nil
}
func (cfg config) Title() string            { return cfg.title }
func (cfg config) VotingID() string         { return cfg.votingID }
func (cfg config) TimeInfo() *util.TimeInfo { return cfg.tInfo }
func (cfg config) Salt1() string            { return cfg.salt1 }
func (cfg config) Salt2() string            { return cfg.salt2 }
func (cfg config) Candidates() []Candidate  { return cfg.candidates }
func (cfg config) ManPubKey() *ecies.PubKey { return cfg.manPubKey }
func (cfg config) VParam() VoteParams       { return cfg.vParam }
func (cfg config) VType() VotingType        { return cfg.vType }
func (cfg config) UchmCid() string          { return cfg.chmCid }
func (cfg config) UnimCid() string          { return cfg.nimCid }
func (cfg config) UivmCid() string          { return cfg.ivmCid }
func (cfg config) ResMapName() string       { return cfg.resMapName }
func (cfg config) UserDataLabels() []string { return cfg.userDataLabels }

func ConfigFromCid(vCfgCid string, is *ipfs.IPFS) (*config, error) {
	m, err := ipfs.FromCid(vCfgCid, is)
	if err != nil {
		return nil, util.NewError("invalid vCfgCid")
	}
	vCfg, err := UnmarshalConfig(m)
	if err != nil {
		return nil, util.NewError("invalid vCfgCid")
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
		ChmCid         string
		NimCid         string
		IvmCid         string
		ResMapName     string
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
		ChmCid:         cfg.chmCid,
		NimCid:         cfg.nimCid,
		IvmCid:         cfg.ivmCid,
		ResMapName:     cfg.resMapName,
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
		ChmCid         string
		NimCid         string
		IvmCid         string
		ResMapName     string
		UserDataLabels []string
	}{}
	err := util.Unmarshal(m, mCfg)
	if err != nil {
		return nil, err
	}

	pub := &ecies.PubKey{}
	err = pub.Unmarshal(mCfg.ManPubKey)
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
		chmCid:         mCfg.ChmCid,
		nimCid:         mCfg.NimCid,
		ivmCid:         mCfg.IvmCid,
		resMapName:     mCfg.ResMapName,
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