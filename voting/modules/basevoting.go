package votingmodule

import (
	"fmt"
	"time"

	"github.com/pilinsin/ipfs-util"
	"github.com/pilinsin/util"
	"github.com/pilinsin/util/crypto"
	rutil "github.com/pilinsin/easy-voting/registration/util"
	vutil "github.com/pilinsin/easy-voting/voting/util"
)

type voting struct {
	is         *ipfs.IPFS
	identity   *rutil.UserIdentity
	vid        string
	salt2      string
	tInfo      *util.TimeInfo
	cands      []vutil.Candidate
	manPubKey  crypto.IPubKey
	verfTopic string
	voteTopic string
	logTopic string
	voteSub         iface.PubSubSubscription
	logSub         iface.PubSubSubscription
	hashVoteMap *vutil.HashVoteMap
	hbmCid string
	uhmCid string
	pimCid     string
	verfMapCid string
	verfMapName string
	resBoxName string
}

func (v *voting) init(vCfgCid string, identity *rutil.UserIdentity, is *ipfs.IPFS) {
	vCfg, err := vutil.ConfigFromCid(vCfgCid, is)
	if err != nil {
		return
	}
	verfCid, err := ipfs.Name.GetCid(vCfg.VerfMapName(), is)
	if err != nil {
		return
	}

	v.is = is
	v.identity = identity
	v.vid = vCfg.VotingID()
	v.salt2 = vCfg.Salt2()
	v.tInfo = vCfg.TimeInfo()
	v.cands = vCfg.Candidates()
	v.manPubKey = vCfg.ManPubKey()
	v.verfTopic = "verfKey_pubsub/" + vCfgCid,
	v.voteTopic = "voting_pubsub/" + vCfgCid,
	v.logTopic = "log_pubsub/" + vCfgCid,
	v.voteSub = is.PubSub().Subscribe("voting_pubsub/" + vCfgCid),
	v.logSub = is.PubSub().Subscribe("log_pubsub/" + vCfgCid),
	v.hashVoteMap = vutil.NewHashVoteMap(100000, vCfg.TimeInfo(), vCfg.VotingID())
	v.hbmCid = vCfg.HbmCid()
	v.uhmCid = vCfg.UhmCid()
	v.pimCid = vCfg.PimCid()
	v.psTopic = "voting_pubsub/" + vCfgCid
	v.verfMapCid = verfCid
	v.verfMapName = vCfg.VerfMapName()
	v.resBoxName = vCfg.ResBoxName()
}
func (v *voting) Load() error{
	hvkm, err := vutil.HashVerfMapFromName(v.verfMapName, v.is)
	if err != nil{return err}

	cids := v.is.PubSub().NextAll(v.logSub)
	for idx, _ := range cids{
		cid := util.Bytes64ToAnyStr(cids[len(cids) - idx - 1])
		hvtm, err := HashVoteMapFromCid(cid, v.is)
		if err != nil{continue}
		if ok := hvtm.VerifyMap(hvkm, v.is); ok{
			v.hashVoteMap = hvtm
			return nil
		}
	}
	return util.NewError("load failed: no valid log")
}
func (v *voting) Close() {
	v.voteSub.Close()
	v.voteSub = nil
	v.logSub.Close()
	v.logSub = nil
	v.is = nil
	v.identity = nil
}

func (v *voting) isCandsMatch(vi vutil.VoteInt) bool {
	if len(v.cands) != len(vi) {
		return false
	}

	for _, ng := range v.candNameGroups() {
		if _, ok := vi[ng]; !ok {
			return false
		}
	}
	return true
}

func (v *voting) candNameGroups() []string {
	ngs := make([]string, len(v.cands))
	for idx, candidate := range v.cands {
		ngs[idx] = candidate.Name + ", " + candidate.Group + " _" + v.vid
	}
	return ngs
}

func (v *voting) Update(){
	hvkm, err := vutil.HashVerfMapFromName(v.verfMapName, v.is)
	if err != nil{return}

	vInfos := v.is.PubSub().NextAll(v.voteSub)
	for _, mvi := range vInfos{
		vInfo := &vutil.VoteInfo{}
		if err := vInfo.Unmarshal(mvi); err != nil{continue}
		v.hashVoteMap.Append(vInfo, hvkm, v.is)
	}
}
func (v *voting) Log(){
	hvkm, err := vutil.HashVerfMapFromName(v.verfMapName, v.is)
	if err != nil{return}
	if ok := v.hashVoteMap.VerifyMap(hvkm, v.is); !ok{return}

	cid := ipfs.File.Add(v.hashVoteMap.Marshal(), v.is)
	v.is.PubSub().Publish(util.AnyStrToBytes64(cid), v.logTopic)
}


func (v voting) uid() (string, bool) {
	mhbm, err := ipfs.File.Get(v.hbmCid, v.is)
	if err != nil {
		return "", false
	}
	hbm, err := rutil.UnmarshalHashBoxMap(mhbm, v.is)
	if err != nil {
		return "", false
	}
	pim, err = vutil.PubIdMapFromCid(v.pimCid, v.is)
	if err != nil {
		return "", false
	}

	uhHash := rutil.NewUhHash(v.salt2, v.identity.UserHash())
	_, ok := hbm.ContainHash(uhHash, v.is)
	ok2 := pim.VerifyIdentity(v.identity, v.is)
	if !ok || !ok2 {
		return "", false
	}
	return pim.ContainIdentity(v.identity, v.is)
}

func (v voting) VerifyIdentity() bool {
	_, ok := v.uid()
	return ok
}

func (v *voting) baseVote(data vutil.VoteInt) error {
	if uid, ok := v.uid(); !ok {
		return util.NewError("cannot get uid")
	} else {
		uvHash := vutil.NewUidVidHash(uid, v.vid)
		uvhHash := vutil.NewUvhHash(uvHash, v.vid)
		uhm, err := vutil.UvhHashMapFromCid(v.uhmCid, v.is)
		if ok := (err == nil) && uhm.ContainHash(uvhHash, v.is); !ok{
			return util.NewError("invalid vote")
		}

		uInfo := vutil.NewUserInfo(uvHash, v.identity.Verify())
		v.is.PubSub().Publish(uInfo.Marshal(), v.verfTopic)
		ticker := time.NewTicker(30*time.Second)
		defer ticker.Stop()
		for {
			hvkm, err := vutil.HashVerfMapFromName(v.verfMapName, v.is)
			if err != nil{return err}
			if hvkm.VerifyUserInfo(uInfo, v.is) {
				fmt.Println("uInfo verified")
				break
			}
			//fmt.Println("wait for registration")
			<-ticker.C
		}

		vb := vutil.NewVotingBox()
		vb.Vote(data, v.identity, v.manPubKey)
		vInfo := vutil.NewVoteInfo(uvHash, vb)
		v.is.PubSub().Publish(vInfo.Marshal(), v.voteTopic)
		return nil
	}
}

func (v voting) baseGetMyVotes() (<-chan *vutil.VoteInt, int, error) {
	mRes, err := ipfs.Name.Get(v.resBoxName, v.is)
	if err != nil{return nil, -1, err}
	resBox, err := vutil.UnmarshalHashVoteMap(mRes)
	if err != nil{return nil, -1, err}
	priKey := resBox.ManPriKey()

	hvkm, err := vutil.HashVerfMapFromName(v.verfMapName, v.is)
	if err != nil{return nil, -1, err}
	if ok := v.hashVoteMap.VerifyMap(hvkm, v.is); !ok{
		return nil, -1, util.NewError("invalid my result")
	}

	ch := make(chan *vutil.VoteInt)
	go func(){
		defer close(ch)
		for vb := range v.hashVoteMap.Next(v.is){
			vi, err := vb.GetVote(v.tInfo, priKey)
			if err == nil{
				ch <- &vi
			}else{
				ch <- nil
			}
		}
	}
	return ch, v.hashVoteMap.Len(), nil
}
func (v voting) baseGetVotes() (<-chan *vutil.VoteInt, int, error) {
	mRes, err := ipfs.Name.Get(v.resBoxName, v.is)
	if err != nil{return nil, -1, err}
	resBox, err := vutil.UnmarshalHashVoteMap(mRes)
	if err != nil{return nil, -1, err}
	res := resBox.HashVoteMap()
	priKey := resBox.ManPriKey()

	hvkm, err := vutil.HashVerfMapFromName(v.verfMapName, v.is)
	if err != nil{return nil, -1, err}
	if ok := res.VerifyMap(hvkm, v.is); !ok{
		return nil, -1, util.NewError("invalid result")
	}

	ch := make(chan *vutil.VoteInt)
	go func(){
		defer close(ch)
		for vb := range res.Next(v.is){
			vi, err := vb.GetVote(v.tInfo, priKey)
			if err == nil{
				ch <- &vi
			}else{
				ch <- nil
			}
		}
	}
	return ch, res.Len(), nil
}

func (v *voting) VerifyHashVerfMap() bool {
	pth, _ := v.is.NameResolve(v.verfMapName)
	mVerfMap, err := v.is.FileGet(pth)
	if err != nil {
		fmt.Println("verfMapName error")
		return false
	}
	hvkm, err := vutil.UnmarshalHashVerfMap(mVerfMap)
	if err != nil {
		fmt.Println("verfMap unmarshal error")
		return false
	}

	if ok := hvkm.VerifyCid(v.verfMapCid, v.is); ok{
		v.verfMapCid = pth.Cid().String()
		return true
	} else {
		fmt.Println("invalid verfMap cid")
		return false
	}
}
