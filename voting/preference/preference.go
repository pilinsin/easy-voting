package preferencevoting

import (
	crsa "crypto/rsa"
	"errors"
	"sort"

	"EasyVoting/util"
	"EasyVoting/voting"
)

type PreferenceVoting struct {
	voting.Voting
}

func New(cfg *voting.InitConfig, ipnsAddrs []string) *PreferenceVoting {
	if ok := voting.VerifyUserID(cfg.KeyFile, ipnsAddrs); !ok {
		util.CheckError(errors.New("Invalid KeyFile"))
		return nil
	}

	pv := &PreferenceVoting{}
	pv.Init(cfg)
	return pv
}

func (pv *PreferenceVoting) IsValidData(vi voting.VoteInt) bool {
	if !pv.NumCandsMatch(len(vi)) {
		return false
	}

	var ps []int
	for _, vote := range vi {
		ps = append(ps, vote)
	}
	sort.Ints(ps)

	for i, v := range ps {
		if i != v {
			return false
		}
	}

	return true
}

func (pv *PreferenceVoting) Type() string {
	return "preferencevoting"
}

func (pv *PreferenceVoting) Vote(data voting.VoteInt, pubKey crsa.PublicKey) string {
	if pv.WithinTime() && pv.IsValidData(data) {
		vd := pv.GenVotingData(data)
		mvd := vd.Marshal()
		return pv.BaseVote(mvd, pubKey)
	}

	util.CheckError(errors.New("Invalid Data"))
	return ""

}

func (pv *PreferenceVoting) Get(ipnsName string, priKey crsa.PrivateKey) voting.VotingData {
	mvd := pv.BaseGet(ipnsName, priKey)
	return voting.UnmarshalVotingData(mvd)
}

func (pv *PreferenceVoting) Count(nameList map[string]struct{}, proKey crsa.PrivateKey) {

}
