package votingutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	pb "github.com/pilinsin/easy-voting/voting/util/pb"
	proto "google.golang.org/protobuf/proto"

	rutil "github.com/pilinsin/easy-voting/registration/util"
	evutil "github.com/pilinsin/easy-voting/util"
	i2p "github.com/pilinsin/go-libp2p-i2p"
	pv "github.com/pilinsin/p2p-verse"
	crdt "github.com/pilinsin/p2p-verse/crdt"
	ipfs "github.com/pilinsin/p2p-verse/ipfs"
	"github.com/pilinsin/util"
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

func encodeVoteParams(vp *VoteParams) *pb.Params {
	return &pb.Params{
		Min:   int32(vp.Min),
		Max:   int32(vp.Max),
		Total: int32(vp.Total),
	}
}
func decodeVoteParams(vp *pb.Params) *VoteParams {
	return &VoteParams{
		Min:   int(vp.GetMin()),
		Max:   int(vp.GetMax()),
		Total: int(vp.GetTotal()),
	}
}

type Candidate struct {
	Name    string
	Group   string
	Url     string
	Image   []byte
	ImgName string
}

func encodeCandidates(cands []*Candidate) []*pb.Candidate {
	pbCands := make([]*pb.Candidate, len(cands))
	for idx, cand := range cands {
		pbCands[idx] = &pb.Candidate{
			Name:    cand.Name,
			Group:   cand.Group,
			Url:     cand.Url,
			Image:   cand.Image,
			ImgName: cand.ImgName,
		}
	}
	return pbCands
}
func decodeCandidates(cands []*pb.Candidate) []*Candidate {
	pbCands := make([]*Candidate, len(cands))
	for idx, cand := range cands {
		pbCands[idx] = &Candidate{
			Name:    cand.GetName(),
			Group:   cand.GetGroup(),
			Url:     cand.GetUrl(),
			Image:   cand.GetImage(),
			ImgName: cand.GetImgName(),
		}
	}
	return pbCands
}

func encodeTimeInfo(ti *util.TimeInfo) *pb.TimeInfo {
	return &pb.TimeInfo{
		Begin: ti.Begin,
		End:   ti.End,
		Loc:   ti.Loc,
	}
}
func decodeTimeInfo(ti *pb.TimeInfo) *util.TimeInfo {
	return &util.TimeInfo{
		Begin: ti.Begin,
		End:   ti.End,
		Loc:   ti.Loc,
	}
}

type VotingStores struct {
	Is  ipfs.Ipfs
	Hkm crdt.IStore
	Ivm crdt.IUpdatableSignatureStore
}

func (vs *VotingStores) Close() {
	if vs.Is != nil {
		vs.Is.Close()
	}
	if vs.Hkm != nil {
		vs.Hkm.Close()
	}
	if vs.Ivm != nil {
		vs.Ivm.Close()
	}
}

type Config struct {
	Title      string
	Time       *util.TimeInfo
	Salt1      string
	Salt2      []byte
	Candidates []*Candidate
	ManPriCid  string
	PubKey     evutil.IPubKey
	Params     *VoteParams
	Type       VotingType
	HkmAddr    string
	IvmAddr    string
	Labels     []string
}

func NewConfig(title, rCfgAddr string, tInfo *util.TimeInfo, cands []*Candidate, vParam *VoteParams, vType VotingType) (string, string, *VotingStores, error) {
	bAddr, rCfgCid, err := evutil.ParseConfigAddr(rCfgAddr)
	if err != nil {
		return "", "", nil, err
	}
	bootstraps := pv.AddrInfosFromString(bAddr)

	baseDir := evutil.BaseDir("voting", "setup")

	ipfsDir := filepath.Join(baseDir, "ipfs")
	os.RemoveAll(ipfsDir)
	is, err := evutil.NewIpfs(i2p.NewI2pHost, ipfsDir, false, bootstraps)
	if err != nil {
		return "", "", nil, err
	}

	rCfg := &rutil.Config{}
	if err := rCfg.FromCid(rCfgCid, is); err != nil {
		is.Close()
		return "", "", nil, err
	}

	storeDir := filepath.Join(baseDir, "store")
	os.RemoveAll(storeDir)
	v := crdt.NewVerse(i2p.NewI2pHost, storeDir, false, bootstraps...)

	uhm, err := v.NewStore(rCfg.UhmAddr, "hash")
	if err != nil {
		is.Close()
		return "", "", nil, err
	}

	skp := evutil.NewSignKeyPair()
	ch := make(chan string)
	go func() {
		defer close(ch)
		ch <- crdt.PubKeyToStr(skp.Verify())
	}()
	hkm, err := v.NewStore(pv.RandString(8), "signature", &crdt.StoreOpts{Pub: skp.Verify(), Priv: skp.Sign()})
	if err != nil {
		is.Close()
		return "", "", nil, err
	}
	hkm, err = v.NewAccessStore(hkm, ch)
	if err != nil {
		is.Close()
		return "", "", nil, err
	}
	hkmAddr := hkm.Address()

	accesses := make(chan string)
	go func() {
		defer close(accesses)
		defer uhm.Close()
		rs, err := uhm.Query()
		if err != nil {
			return
		}
		for res := range rs.Next() {
			userPubKey, err := evutil.UnmarshalPub(res.Value)
			if err != nil {
				continue
			}
			ukp := NewUserKeyPair()
			enc, err := userPubKey.Encrypt(ukp.Marshal())
			if err != nil {
				continue
			}
			if err := hkm.Put(res.Key, enc); err != nil {
				continue
			}
			accesses <- crdt.PubKeyToStr(ukp.Verify())
			fmt.Println(res.Key, "is appended")
		}
		rs.Close()
	}()

	end := tInfo.EndTime()
	opt := &crdt.StoreOpts{TimeLimit: end}
	tmp, err := v.NewStore(pv.RandString(8), "updatableSignature", opt)
	if err != nil {
		is.Close()
		hkm.Close()
		return "", "", nil, err
	}
	tmp, err = v.NewAccessStore(tmp, accesses)
	if err != nil {
		is.Close()
		hkm.Close()
		return "", "", nil, err
	}
	ivm := tmp.(crdt.IUpdatableSignatureStore)
	ivmAddr := ivm.Address()

	encKeyPair := evutil.NewPubKeyPair()
	cg, err := ipfs.NewCidGetter()
	if err != nil {
		is.Close()
		hkm.Close()
		ivm.Close()
		return "", "", nil, err
	}
	defer cg.Close()
	m, _ := encKeyPair.Private().Raw()
	cid, err := cg.Get(m)
	if err != nil {
		is.Close()
		hkm.Close()
		ivm.Close()
		return "", "", nil, err
	}

	vCfg := &Config{
		Title:      title,
		Time:       tInfo,
		Salt1:      rCfg.Salt1,
		Salt2:      rCfg.Salt2,
		Candidates: cands,
		ManPriCid:  cid,
		PubKey:     encKeyPair.Public(),
		Params:     vParam,
		Type:       vType,
		HkmAddr:    hkmAddr,
		IvmAddr:    ivmAddr,
		Labels:     rCfg.Labels,
	}
	vCfgCid, err := vCfg.toCid(is)
	if err != nil {
		return "", "", nil, err
	}

	manId := &ManIdentity{
		Priv: encKeyPair.Private(),
	}
	vStores := &VotingStores{
		Is:  is,
		Hkm: hkm,
		Ivm: ivm,
	}

	return "v/" + vCfgCid, manId.toString(), vStores, nil
}

func (cfg Config) Marshal() []byte {
	mpub, _ := cfg.PubKey.Raw()
	pbCfg := &pb.Config{
		Title:      cfg.Title,
		Time:       encodeTimeInfo(cfg.Time),
		Salt1:      cfg.Salt1,
		Salt2:      cfg.Salt2,
		Candidates: encodeCandidates(cfg.Candidates),
		ManCid:     cfg.ManPriCid,
		PubKey:     mpub,
		Params:     encodeVoteParams(cfg.Params),
		Type:       int32(cfg.Type),
		HkmAddr:    cfg.HkmAddr,
		IvmAddr:    cfg.IvmAddr,
		Labels:     cfg.Labels,
	}
	m, _ := proto.Marshal(pbCfg)
	return m
}
func (cfg *Config) Unmarshal(m []byte) error {
	pbCfg := &pb.Config{}
	if err := proto.Unmarshal(m, pbCfg); err != nil {
		return err
	}

	pubKey, err := evutil.UnmarshalPub(pbCfg.GetPubKey())
	if err != nil {
		return err
	}

	cfg.Title = pbCfg.Title
	cfg.Time = decodeTimeInfo(pbCfg.Time)
	cfg.Salt1 = pbCfg.Salt1
	cfg.Salt2 = pbCfg.Salt2
	cfg.Candidates = decodeCandidates(pbCfg.Candidates)
	cfg.ManPriCid = pbCfg.ManCid
	cfg.PubKey = pubKey
	cfg.Params = decodeVoteParams(pbCfg.Params)
	cfg.Type = VotingType(pbCfg.Type)
	cfg.HkmAddr = pbCfg.HkmAddr
	cfg.IvmAddr = pbCfg.IvmAddr
	cfg.Labels = pbCfg.Labels
	return nil
}

func (cfg *Config) toCid(is ipfs.Ipfs) (string, error) {
	return is.Add(cfg.Marshal(), time.Second*10)
}
func (cfg *Config) FromCid(vCfgCid string, is ipfs.Ipfs) error {
	m, err := is.Get(vCfgCid, time.Second*5)
	if err != nil {
		return errors.New("get from vCfgCid error")
	}
	if err := cfg.Unmarshal(m); err != nil {
		return errors.New("unmarshal vCfg error")
	}
	return nil
}

func (cfg *Config) ShuffleCandidates() {
	n := len(cfg.Candidates)
	for i := 0; i < n-1; i++ {
		j := i + util.RandInt(n-i) //[i, n)
		cfg.Candidates[i], cfg.Candidates[j] = cfg.Candidates[j], cfg.Candidates[i]
	}
}
func (cfg *Config) CandNameGroups() []string {
	ngs := make([]string, len(cfg.Candidates))
	for idx, candidate := range cfg.Candidates {
		ngs[idx] = candidate.Name + ", " + candidate.Group
	}
	return ngs
}
