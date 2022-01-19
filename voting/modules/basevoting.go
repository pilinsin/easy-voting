package votingmodule

import (
	"fmt"
	"time"
	"strings"

	"EasyVoting/ipfs"
	rutil "EasyVoting/registration/util"
	"EasyVoting/util"
	"EasyVoting/util/crypto"
	vutil "EasyVoting/voting/util"
)

type voting struct {
	is         *ipfs.IPFS
	identity   *rutil.UserIdentity
	vid        string
	tInfo      *util.TimeInfo
	salt2      string
	cands      []vutil.Candidate
	manPubKey  crypto.IPubKey
	chmCid     string
	nimCid     string
	ivmCid     string
	resMapName string
	psTopic string
	verfMapCid string
	verfMapName string
}

func (v *voting) init(vCfgCid string, identity *rutil.UserIdentity, is *ipfs.IPFS) {
	vCfg, err := vutil.ConfigFromCid(vCfgCid, is)
	if err != nil {
		return
	}
	verfCid, err := ipfs.CidFromName(vCfg.VerfMapName(), is)
	if err != nil {
		return
	}

	v.is = is
	v.identity = identity
	v.vid = vCfg.VotingID()
	v.tInfo = vCfg.TimeInfo()
	v.salt2 = vCfg.Salt2()
	v.cands = vCfg.Candidates()
	v.manPubKey = vCfg.ManPubKey()
	v.chmCid = vCfg.UchmCid()
	v.nimCid = vCfg.UnimCid()
	v.ivmCid = vCfg.UivmCid()
	v.resMapName = vCfg.ResMapName()
	v.psTopic = "voting_pubsub/" + vCfgCid
	v.verfMapCid = verfCid
	v.verfMapName = vCfg.VerfMapName()
}
func (v *voting) Close() {
	v.is = nil
	v.identity = nil
}

func (v *voting) candNameGroups() []string {
	ngs := make([]string, len(v.cands))
	for idx, candidate := range v.cands {
		ngs[idx] = candidate.Name + ", " + candidate.Group + " _" + v.vid
	}
	return ngs
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

func (v voting) uid() (string, bool) {
	chm := &rutil.ConstHashMap{}
	err := chm.FromCid(v.chmCid, v.is)
	if err != nil {
		return "", false
	}
	nim := &vutil.NameIdMap{}
	err = nim.FromCid(v.nimCid, v.is)
	if err != nil {
		return "", false
	}

	uhHash := rutil.NewUhHash(v.is, v.salt2, v.identity.UserHash())
	hash := chm.ContainHash(uhHash, v.is)
	name := nim.VerifyIdentity(v.identity, v.is)
	if !(hash && name) {
		return "", false
	}
	return nim.ContainIdentity(v.identity, v.is)
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
		ivm, err := vutil.IdVotingMapFromCid(v.ivmCid, v.is)
		if err != nil {
			return err
		}
		ivm.Vote(uvHash, data, v.identity, v.manPubKey, v.is)

		uInfo := vutil.NewUserInfo(uvHash, v.identity.Sign().Verify())
		v.is.PubSubPublish(uInfo.Marshal(), v.psTopic)
		ticker := time.NewTicker(30*time.Second)
		defer ticker.Stop()
		for {
			verfMap, err := vutil.IdVerfKeyMapFromName(v.verfMapName, v.is)
			if err != nil {
				return err
			}
			if verfMap.VerifyUserInfo(uInfo, v.is) {
				fmt.Println("uInfo verified")
				return nil
			}
			//fmt.Println("wait for registration")
			<-ticker.C
		}
		return nil
	}
}

func (v voting) baseGetMyVote() (*vutil.VoteInt, error) {
	if uid, ok := v.uid(); !ok {
		return nil, util.NewError("cannot get uid")
	} else {
		uvHash := vutil.NewUidVidHash(uid, v.vid)
		ivm, err := vutil.IdVotingMapFromCid(v.ivmCid, v.is)
		if err != nil {
			return nil, err
		}
		if vb, ok := ivm.ContainHash(uvHash, v.is); !ok {
			return nil, util.NewError("invalid uvHash")
		} else {
			if sv, err := vb.GetMyVote(v.identity); err != nil{
				return nil, err
			} else {
				if ok := sv.Verify(v.identity.Sign().Verify()); !ok {
					return nil, util.NewError("invalid verfKey")
				} else {
					return sv.Vote(v.tInfo)
				}
			}
		}
	}
}
func (v voting) baseGetVotes() (<-chan *vutil.VoteInt, int, error) {
	resMap, err := vutil.ResultMapFromName(v.resMapName, v.is)
	if err != nil {
		return nil, -1, err
	}

	return resMap.Next(v.is), resMap.NumVoters(), nil
}

func (v *voting) VerifyIdVerfKeyMap() bool {
	ivm, err := vutil.IdVotingMapFromCid(v.ivmCid, v.is)
	if err != nil {
		return false
	}

	pth, _ := v.is.NameResolve(v.verfMapName)
	mVerfMap, err := v.is.FileGet(pth)
	if err != nil {
		fmt.Println("verfMapName error")
		return false
	}
	verfMap, err := vutil.UnmarshalIdVerfKeyMap(mVerfMap)
	if err != nil {
		fmt.Println("verfMap unmarshal error")
		return false
	}

	if ok := verfMap.VerifyIds(ivm, v.is); !ok {
		fmt.Println("a verfKey corresponding to a uvHash is not registered")
		return false
	}
	if verfMap.VerifyCid(v.verfMapCid, v.is) {
		v.verfMapCid = strings.TrimPrefix(pth.String(), "/ipfs/")
		return true
	} else {
		fmt.Println("invalid verfMap cid")
		return false
	}
}

func (v voting) VerifyResultMap() (bool, error) {
	verfMap, err := vutil.IdVerfKeyMapFromName(v.verfMapName, v.is)
	if err != nil {
		return false, util.NewError("invalid ivmCid")
	}
	resMap, err := vutil.ResultMapFromName(v.resMapName, v.is)
	if err != nil {
		return false, util.NewError("resMap does not exist")
	}

	return resMap.VerifyVotes(verfMap, v.is), nil
}
