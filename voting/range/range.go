package rangevoting

import (
	crsa "crypto/rsa"
	"errors"

	"EasyVoting/util"
	"EasyVoting/voting"
)

type RangeVoting struct {
	voting.Voting
	min int
	max int
}

func New(cfg *voting.InitConfig, ipnsAddrs []string, min int, max int) *RangeVoting {
	if ok := voting.VerifyUserID(cfg.KeyFile, ipnsAddrs); !ok {
		util.CheckError(errors.New("Invalid KeyFile"))
		return nil
	}

	rv := &RangeVoting{min: min, max: max}
	rv.Init(cfg)
	return rv
}

func (rv *RangeVoting) IsValidData(vi voting.VoteInt) bool {
	if !rv.NumCandsMatch(len(vi)) {
		return false
	}

	for _, vote := range vi {
		if vote < rv.min || vote > rv.max {
			return false
		}
	}

	return true
}

func (rv *RangeVoting) Type() string {
	return "rangevoting"
}

func (rv *RangeVoting) Vote(data voting.VoteInt, pubKey crsa.PublicKey) string {
	if rv.WithinTime() && rv.IsValidData(data) {
		vd := rv.GenVotingData(data)
		mvd := vd.Marshal()
		return rv.BaseVote(mvd, pubKey)
	}

	util.CheckError(errors.New("Invalid Data"))
	return ""

}

func (rv *RangeVoting) Get(ipnsName string, priKey crsa.PrivateKey) voting.VotingData {
	mvd := rv.BaseGet(ipnsName, priKey)
	return voting.UnmarshalVotingData(mvd)
}

func (rv *RangeVoting) Count(nameList map[string]struct{}, proKey crsa.PrivateKey) {

}
