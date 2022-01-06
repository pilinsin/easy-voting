package votingmodule

import (
	"EasyVoting/ipfs"
	rutil "EasyVoting/registration/util"
	"EasyVoting/util"
	"EasyVoting/util/ecies"
	vutil "EasyVoting/voting/util"
)

type voting struct {
	is         *ipfs.IPFS
	identity   *rutil.UserIdentity
	vid        string
	tInfo      *util.TimeInfo
	salt2      string
	cands      []vutil.Candidate
	manPubKey  *ecies.PubKey
	chmCid     string
	nimCid     string
	ivmCid     string
	resMapName string
}

func (v *voting) init(vCfgCid string, identity *rutil.UserIdentity, is *ipfs.IPFS) {
	vCfg, err := vutil.ConfigFromCid(vCfgCid, is)
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
		if vb, _, ok := ivm.ContainHash(uvHash, v.is); !ok {
			return nil, util.NewError("invalid uvHash")
		} else {
			msv, _ := vb.GetMyVote(v.identity)
			if sv, err := vutil.UnmarshalSignedVote(msv); err != nil {
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
func (v voting) VerifyResultMap() (bool, error) {
	ivm, err := vutil.IdVotingMapFromCid(v.ivmCid, v.is)
	if err != nil {
		return false, util.NewError("invalid ivmCid")
	}
	resMap, err := vutil.ResultMapFromName(v.resMapName, v.is)
	if err != nil {
		return false, util.NewError("resMap does not exist")
	}

	return resMap.VerifyVotes(ivm, v.is), nil
}
