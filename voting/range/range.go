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

func New(cfg *voting.InitConfig, huidListAddrs []string, min int, max int) *RangeVoting {
	if ok := voting.VerifyUserID(cfg.Is, cfg.UserID, huidListAddrs); !ok {
		util.CheckError(errors.New("Invalid userID"))
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
		mvi := data.Marshal()
		return rv.BaseVote(mvi, pubKey)
	}

	util.CheckError(errors.New("Invalid Data"))
	return ""

}

func (rv *RangeVoting) Get(ipnsName string, priKey crsa.PrivateKey) voting.VoteInt {
	mvi := rv.BaseGet(ipnsName, priKey)
	return voting.UnmarshalVoteInt(mvi)
}

func (rv *RangeVoting) Count(nameList map[string]struct{}, proKey crsa.PrivateKey) {

}
